package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bamboo-firewall/be"
	"github.com/bamboo-firewall/be/cmd/server/pkg/common/errlist"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
	"github.com/bamboo-firewall/be/cmd/server/pkg/repository"
	"github.com/bamboo-firewall/be/cmd/server/pkg/selector"
	"github.com/bamboo-firewall/be/domain/model"
)

func NewHEP(policyMongo *repository.PolicyDB) *hep {
	return &hep{
		storage: policyMongo,
	}
}

type hep struct {
	storage be.Storage
}

func (ds *hep) Create(ctx context.Context, input *model.CreateHostEndpointInput) (*entity.HostEndpoint, *ierror.Error) {
	// ToDo: use transaction and lock row
	hepExisted, coreErr := ds.storage.GetHostEndpointByName(ctx, input.Metadata.Name)
	if coreErr != nil && !errors.Is(coreErr, errlist.ErrNotFoundHostEndpoint) {
		return nil, httpbase.ErrDatabase(ctx, "get host endpoint failed").SetSubError(coreErr)
	}

	var ports []entity.HostEndpointSpecPort
	for _, port := range input.Spec.Ports {
		ports = append(ports, entity.HostEndpointSpecPort{
			Name:     port.Name,
			Port:     port.Port,
			Protocol: port.Protocol,
		})
	}

	hepEntity := &entity.HostEndpoint{
		ID:      primitive.NewObjectID(),
		UUID:    uuid.New().String(),
		Version: 1,
		Metadata: entity.HostEndpointMetadata{
			Name:   input.Metadata.Name,
			Labels: input.Metadata.Labels,
		},
		Spec: entity.HostEndpointSpec{
			InterfaceName: input.Spec.InterfaceName,
			IPs:           input.Spec.IPs,
			Ports:         ports,
		},
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if !errors.Is(coreErr, errlist.ErrNotFoundHostEndpoint) {
		hepEntity.ID = hepExisted.ID
		hepEntity.UUID = hepExisted.UUID
		hepEntity.Version = hepExisted.Version + 1
		hepEntity.CreatedAt = hepExisted.CreatedAt
	}

	if coreErr = ds.storage.UpsertHostEndpoint(ctx, hepEntity); coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "Create host endpoint failed").SetSubError(coreErr)
	}
	return hepEntity, nil
}

func (ds *hep) Delete(ctx context.Context, name string) *ierror.Error {
	if coreErr := ds.storage.DeleteHostEndpointByName(ctx, name); coreErr != nil {
		return httpbase.ErrDatabase(ctx, "Delete host endpoint failed").SetSubError(coreErr)
	}
	return nil
}

func (ds *hep) FetchPolicies(ctx context.Context, input *model.FetchPoliciesInput) (*model.HostEndPointPolicy, *ierror.Error) {
	hepEntity, coreErr := ds.storage.GetHostEndpointByName(ctx, input.Name)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrNotFoundHostEndpoint) {
			return nil, httpbase.ErrNotFound(ctx, "NotFound").SetSubError(coreErr)
		}
		return nil, httpbase.ErrDatabase(ctx, "Get host endpoint failed").SetSubError(coreErr)
	}

	policies, err := ds.storage.ListGNP(ctx)
	if err != nil {
		return nil, httpbase.ErrDatabase(ctx, "List policies failed").SetSubError(coreErr)
	}

	sets, err := ds.storage.ListGNS(ctx)
	if err != nil {
		return nil, httpbase.ErrDatabase(ctx, "List sets failed").SetSubError(coreErr)
	}

	var (
		parsedPolicies []*model.ParsedPolicy
		parsedSets     []*model.ParsedSet
	)
	parsedSetsMap := make(map[string]struct{})
	gnpVersions := make(map[string]uint)
	gnsVersions := make(map[string]uint)
	for _, policy := range policies {
		sel, errParse := selector.Parse(policy.Spec.Selector)
		if errParse != nil {
			slog.Warn("malformed selector", "policy_uuid", policy.UUID, "selector", policy.Spec.Selector, "err", errParse)
			continue
		}
		if !sel.Evaluate(hepEntity.Metadata.Labels) {
			continue
		}
		gnpVersions[policy.UUID] = policy.Version

		inboundRules := make([]*model.ParsedRule, 0)
		outboundRules := make([]*model.ParsedRule, 0)
		for _, rule := range policy.Spec.Ingress {
			parsedRule, parsedSetsPerRule := parseRule(policy, &rule, sets, parsedSetsMap, gnsVersions)
			inboundRules = append(inboundRules, parsedRule)
			parsedSets = append(parsedSets, parsedSetsPerRule...)
		}
		for _, rule := range policy.Spec.Egress {
			parsedRule, parsedSetsPerRule := parseRule(policy, &rule, sets, parsedSetsMap, gnsVersions)
			outboundRules = append(outboundRules, parsedRule)
			parsedSets = append(parsedSets, parsedSetsPerRule...)
		}
		parsedPolicies = append(parsedPolicies, &model.ParsedPolicy{
			UUID:          policy.UUID,
			Version:       policy.Version,
			Name:          policy.Metadata.Name,
			InboundRules:  inboundRules,
			OutboundRules: outboundRules,
		})
	}

	return &model.HostEndPointPolicy{
		MetaData: model.HostEndPointPolicyMetadata{
			HEPVersion:  hepEntity.Version,
			GNPVersions: gnpVersions,
			GNSVersions: gnsVersions,
		},
		HEP:            hepEntity,
		ParsedPolicies: parsedPolicies,
		ParsedSets:     parsedSets,
	}, nil
}

func parseRule(policy *entity.GlobalNetworkPolicy, rule *entity.GNPSpecRule, sets []*entity.GlobalNetworkSet,
	parsedSetsMap map[string]struct{}, gnsVersions map[string]uint) (*model.ParsedRule, []*model.ParsedSet) {
	var (
		protocol           string
		isProtocolNegative bool
		srcGNSSetNames     []string
		srcNets            []string
		isSrcNetNegative   bool
		srcPorts           []string
		isSrcPortNegative  bool
		dstGNSSetNames     []string
		dstNets            []string
		isDstNetNegative   bool
		dstPorts           []string
		isDstPortNegative  bool
		parsedSets         []*model.ParsedSet
	)
	if rule.Protocol != "" {
		protocol = rule.Protocol
		isProtocolNegative = false
	} else if rule.NotProtocol != "" {
		protocol = rule.NotProtocol
		isProtocolNegative = true
	}

	// get sets match if selector is available
	if len(rule.Source.Selector) > 0 {
		for {
			selSource, errParseSource := selector.Parse(rule.Source.Selector)
			if errParseSource != nil {
				slog.Warn("malformed selector in source", "policy_uuid", policy.UUID, "selector", rule.Source.Selector, "err", errParseSource)
				break
			}
			for _, set := range sets {
				if !selSource.Evaluate(set.Metadata.Labels) {
					continue
				}
				srcGNSSetNames = append(srcGNSSetNames, set.Metadata.Name)
				if _, ok := parsedSetsMap[set.UUID]; !ok {
					parsedSetsMap[set.UUID] = struct{}{}
					gnsVersions[set.UUID] = set.Version
					parsedSets = append(parsedSets, entityToParsedSet(set))
				}
			}
			break
		}
	}

	if len(rule.Source.Nets) > 0 {
		srcNets = rule.Source.Nets
		isSrcNetNegative = false
	} else {
		srcNets = rule.Source.NotNets
		isSrcNetNegative = true
	}
	if len(rule.Source.Ports) > 0 {
		srcPorts = convertPorts(rule.Source.Ports)
		isSrcPortNegative = false
	} else {
		srcPorts = convertPorts(rule.Source.NotPorts)
		isSrcPortNegative = true
	}

	// get sets match if selector is available
	if len(rule.Destination.Selector) > 0 {
		for {
			selDst, errParseDst := selector.Parse(rule.Destination.Selector)
			if errParseDst != nil {
				slog.Warn("malformed selector in destination", "policy_uuid", policy.UUID, "selector", rule.Source.Selector, "err", errParseDst)
				break
			}
			for _, set := range sets {
				if !selDst.Evaluate(set.Metadata.Labels) {
					continue
				}
				dstGNSSetNames = append(dstGNSSetNames, set.Metadata.Name)
				if _, ok := parsedSetsMap[set.UUID]; !ok {
					parsedSetsMap[set.UUID] = struct{}{}
					gnsVersions[set.UUID] = set.Version
					parsedSets = append(parsedSets, entityToParsedSet(set))
				}
			}
			break
		}
	}

	if len(rule.Destination.Nets) > 0 {
		dstNets = rule.Destination.Nets
		isDstNetNegative = false
	} else {
		dstNets = rule.Destination.NotNets
		isDstNetNegative = true
	}
	if len(rule.Destination.Ports) > 0 {
		dstPorts = convertPorts(rule.Destination.Ports)
		isDstPortNegative = false
	} else {
		dstPorts = convertPorts(rule.Destination.NotPorts)
		isDstPortNegative = true
	}
	return &model.ParsedRule{
		Action:             rule.Action,
		IPVersion:          rule.IPVersion,
		Protocol:           protocol,
		IsProtocolNegative: isProtocolNegative,
		SrcGNSNetNames:     srcGNSSetNames,
		SrcNets:            srcNets,
		IsSrcNetNegative:   isSrcNetNegative,
		SrcPorts:           srcPorts,
		IsSrcPortNegative:  isSrcPortNegative,
		DstGNSNetNames:     dstGNSSetNames,
		DstNets:            dstNets,
		IsDstNetNegative:   isDstNetNegative,
		DstPorts:           dstPorts,
		IsDstPortNegative:  isDstPortNegative,
	}, parsedSets
}

func entityToParsedSet(set *entity.GlobalNetworkSet) *model.ParsedSet {
	return &model.ParsedSet{
		Name:      set.Metadata.Name,
		IPVersion: set.Metadata.IPVersion,
		Nets:      set.Spec.Nets,
	}
}

func convertPorts(ports []interface{}) []string {
	var portStrings []string
	for _, port := range ports {
		portStrings = append(portStrings, fmt.Sprint(port))
	}
	return portStrings
}
