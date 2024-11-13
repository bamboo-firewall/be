package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bamboo-firewall/be"
	"github.com/bamboo-firewall/be/cmd/server/pkg/common/errlist"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
	"github.com/bamboo-firewall/be/cmd/server/pkg/net"
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
	ipsV4, ipsV6 := exactIPs(input.Spec.IPs)
	if len(ipsV4) == 0 {
		return nil, httpbase.ErrBadRequest(ctx, "required at least one ip version 4")
	}
	if input.Spec.TenantID == 0 {
		input.Spec.TenantID = entity.DefaultTenantID
	}
	var ipString string
	if input.Spec.IP == "" {
		ipString = ipsV4[0]
	} else {
		ipString = input.Spec.IP
	}

	ip := net.ParseIP(ipString)

	hepEntity := &entity.HostEndpoint{
		ID:   primitive.NewObjectID(),
		UUID: entity.NewMinifyUUID(),
		Metadata: entity.HostEndpointMetadata{
			Name:   input.Metadata.Name,
			Labels: input.Metadata.Labels,
		},
		Spec: entity.HostEndpointSpec{
			InterfaceName: input.Spec.InterfaceName,
			IP:            net.IPToInt(*ip),
			TenantID:      input.Spec.TenantID,
			IPs:           input.Spec.IPs,
			IPsV4:         ipsV4,
			IPsV6:         ipsV6,
		},
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	hepEntity, coreErr := ds.storage.UpsertHostEndpoint(ctx, hepEntity)
	if coreErr != nil {
		slog.Error("warning", "err", coreErr)
		if errors.Is(coreErr, errlist.ErrDuplicateHostEndpoint) {
			return nil, httpbase.ErrBadRequest(ctx, "duplicate host endpoint").SetSubError(coreErr)
		}
		return nil, httpbase.ErrDatabase(ctx, "create host endpoint failed").SetSubError(coreErr)
	}
	return hepEntity, nil
}

func exactIPs(ips []string) (ipsV4, ipsV6 []string) {
	for _, ipString := range ips {
		ip := net.ParseIP(ipString)
		if ip == nil {
			slog.Warn("malformed ip", "ip", ipString)
			continue
		}
		if ip.Version() == int(entity.IPVersion4) {
			ipsV4 = append(ipsV4, ip.String())
		} else if ip.Version() == int(entity.IPVersion6) {
			ipsV6 = append(ipsV6, ip.String())
		}
	}
	return
}

func (ds *hep) Get(ctx context.Context, input *model.GetHostEndpointInput) (*entity.HostEndpoint, *ierror.Error) {
	hepEntity, coreErr := ds.storage.GetHostEndpoint(ctx, input)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrNotFoundHostEndpoint) {
			return nil, httpbase.ErrNotFound(ctx, "not found").SetSubError(coreErr)
		}
		return nil, httpbase.ErrDatabase(ctx, "get host endpoint failed").SetSubError(coreErr)
	}
	return hepEntity, nil
}

func (ds *hep) List(ctx context.Context, input *model.ListHostEndpointsInput) ([]*entity.HostEndpoint, *ierror.Error) {
	hepsEntity, coreErr := ds.storage.ListHostEndpoints(ctx, input)
	if coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "list host endpoints failed").SetSubError(coreErr)
	}
	return hepsEntity, nil
}

func (ds *hep) Delete(ctx context.Context, input *model.DeleteHostEndpointInput) *ierror.Error {
	if input.TenantID == 0 {
		input.TenantID = entity.DefaultTenantID
	}
	var ipString string
	if input.IP == "" {
		ipsV4, _ := exactIPs(input.IPs)
		if len(ipsV4) == 0 {
			return httpbase.ErrBadRequest(ctx, "required at least one ip version 4")
		}
		ipString = ipsV4[0]
	} else {
		ipString = input.IP
	}

	ip := net.ParseIP(ipString)
	if coreErr := ds.storage.DeleteHostEndpoint(ctx, input.TenantID, net.IPToInt(*ip)); coreErr != nil {
		return httpbase.ErrDatabase(ctx, "delete host endpoint failed").SetSubError(coreErr)
	}
	return nil
}

func (ds *hep) FetchPolicies(ctx context.Context, input *model.ListHostEndpointsInput) ([]*model.HostEndpointPolicy, *ierror.Error) {
	heps, coreErr := ds.storage.ListHostEndpoints(ctx, nil)
	if coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "list host endpoint failed").SetSubError(coreErr)
	}

	gnps, coreErr := ds.storage.ListGNPs(ctx, &model.ListGNPsInput{IsOrder: true})
	if coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "list global network policy failed").SetSubError(coreErr)
	}

	gnss, coreErr := ds.storage.ListGNSs(ctx)
	if coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "list global network set failed").SetSubError(coreErr)
	}

	var (
		hepPolicies []*model.HostEndpointPolicy
	)

	rp := &ruleParser{
		parsedHEPsMap: make(map[string]struct{}),
		hepVersions:   make(map[string]uint),
		parsedGNSsMap: make(map[string]struct{}),
		gnsVersions:   make(map[string]uint),
	}

	for _, hepEntity := range heps {
		if input != nil {
			if input.TenantID != nil && input.IP != nil {
				if hepEntity.Spec.TenantID != *input.TenantID || hepEntity.Spec.IP != *input.IP {
					continue
				}
			}
		}

		var (
			parsedGNPs  []*model.ParsedGNP
			gnpVersions = make(map[string]uint)
		)
		for _, policy := range gnps {
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
				inboundRules = append(inboundRules, rp.parseRule(policy, &rule, heps, gnss))
			}
			for _, rule := range policy.Spec.Egress {
				outboundRules = append(outboundRules, rp.parseRule(policy, &rule, heps, gnss))
			}
			parsedGNPs = append(parsedGNPs, &model.ParsedGNP{
				UUID:          policy.UUID,
				Version:       policy.Version,
				Name:          policy.Metadata.Name,
				InboundRules:  inboundRules,
				OutboundRules: outboundRules,
			})
		}
		hepPolicies = append(hepPolicies, &model.HostEndpointPolicy{
			MetaData: model.HostEndpointPolicyMetadata{
				GNPVersions: gnpVersions,
				HEPVersions: rp.hepVersions,
				GNSVersions: rp.gnsVersions,
			},
			HEP:        hepEntity,
			ParsedGNPs: parsedGNPs,
			ParsedHEPs: rp.parsedHEPs,
			ParsedGNSs: rp.parsedGNSs,
		})
	}

	return hepPolicies, nil
}

type ruleParser struct {
	parsedHEPs    []*model.ParsedHEP
	parsedHEPsMap map[string]struct{}
	hepVersions   map[string]uint
	parsedGNSs    []*model.ParsedGNS
	parsedGNSsMap map[string]struct{}
	gnsVersions   map[string]uint
}

func (r *ruleParser) parseRule(policy *entity.GlobalNetworkPolicy, rule *entity.GNPSpecRule, heps []*entity.HostEndpoint, gnss []*entity.GlobalNetworkSet) *model.ParsedRule {
	var (
		protocol           interface{}
		isProtocolNegative bool
		srcGNSUUIDs        []string
		srcHEPUUIDs        []string
		srcNets            []string
		isSrcNetNegative   bool
		srcPorts           []string
		isSrcPortNegative  bool
		dstGNSUUIDs        []string
		dstHEPUUIDs        []string
		dstNets            []string
		isDstNetNegative   bool
		dstPorts           []string
		isDstPortNegative  bool
	)
	if rule.Protocol != nil {
		protocol = rule.Protocol
		isProtocolNegative = false
	} else if rule.NotProtocol != nil {
		protocol = rule.NotProtocol
		isProtocolNegative = true
	}

	// get global network set match if selector is available
	if rule.Source != nil {
		if len(rule.Source.Selector) > 0 {
			hepUUIDs, gnsUUIDs, err := r.handleSelector(rule.Source.Selector, rule.IPVersion, heps, gnss)
			if err != nil {
				slog.Warn("malformed selector in source", "policy_uuid", policy.UUID, "selector", rule.Source.Selector, "err", err)
			}
			if len(hepUUIDs) > 0 {
				srcHEPUUIDs = append(srcHEPUUIDs, hepUUIDs...)
			}
			if len(gnsUUIDs) > 0 {
				srcGNSUUIDs = append(srcGNSUUIDs, gnsUUIDs...)
			}
		}

		if len(rule.Source.Nets) > 0 {
			srcNets = rule.Source.Nets
			isSrcNetNegative = false
		} else if len(rule.Source.Nets) > 0 {
			srcNets = rule.Source.NotNets
			isSrcNetNegative = true
		}
		if len(rule.Source.Ports) > 0 {
			srcPorts = convertPorts(rule.Source.Ports)
			isSrcPortNegative = false
		} else if len(rule.Source.NotPorts) > 0 {
			srcPorts = convertPorts(rule.Source.NotPorts)
			isSrcPortNegative = true
		}
	}
	// get global network set match if selector is available
	if rule.Destination != nil {
		if len(rule.Destination.Selector) > 0 {
			hepUUIDs, gnsUUIDs, err := r.handleSelector(rule.Destination.Selector, rule.IPVersion, heps, gnss)
			if err != nil {
				slog.Warn("malformed selector in destination", "policy_uuid", policy.UUID, "selector", rule.Source.Selector, "err", err)
			}
			if len(hepUUIDs) > 0 {
				dstHEPUUIDs = append(dstHEPUUIDs, hepUUIDs...)
			}
			if len(gnsUUIDs) > 0 {
				dstGNSUUIDs = append(dstGNSUUIDs, gnsUUIDs...)
			}
		}

		if len(rule.Destination.Nets) > 0 {
			dstNets = rule.Destination.Nets
			isDstNetNegative = false
		} else if len(rule.Destination.NotNets) > 0 {
			dstNets = rule.Destination.NotNets
			isDstNetNegative = true
		}
		if len(rule.Destination.Ports) > 0 {
			dstPorts = convertPorts(rule.Destination.Ports)
			isDstPortNegative = false
		} else if len(rule.Destination.NotPorts) > 0 {
			dstPorts = convertPorts(rule.Destination.NotPorts)
			isDstPortNegative = true
		}
	}
	return &model.ParsedRule{
		Action:             rule.Action,
		IPVersion:          int(rule.IPVersion),
		Protocol:           protocol,
		IsProtocolNegative: isProtocolNegative,
		SrcGNSUUIDs:        srcGNSUUIDs,
		SrcHEPUUIDs:        srcHEPUUIDs,
		SrcNets:            srcNets,
		IsSrcNetNegative:   isSrcNetNegative,
		SrcPorts:           srcPorts,
		IsSrcPortNegative:  isSrcPortNegative,
		DstGNSUUIDs:        dstGNSUUIDs,
		DstHEPUUIDs:        dstHEPUUIDs,
		DstNets:            dstNets,
		IsDstNetNegative:   isDstNetNegative,
		DstPorts:           dstPorts,
		IsDstPortNegative:  isDstPortNegative,
	}
}

func (r *ruleParser) handleSelector(selectorString string, ruleIPVersion entity.IPVersion, heps []*entity.HostEndpoint, gnss []*entity.GlobalNetworkSet) ([]string, []string, error) {
	var (
		hepUUIDs []string
		gnsUUIDs []string
	)
	sel, errParse := selector.Parse(selectorString)
	if errParse != nil {
		return nil, nil, errParse
	}
	for _, ep := range heps {
		if !sel.Evaluate(ep.Metadata.Labels) {
			continue
		}
		if !((ruleIPVersion == entity.IPVersion4 && len(ep.Spec.IPsV4) > 0) || (ruleIPVersion == entity.IPVersion6 && len(ep.Spec.IPsV6) > 0)) {
			continue
		}
		hepUUIDs = append(hepUUIDs, ep.UUID)
		if _, ok := r.parsedHEPsMap[ep.UUID]; !ok {
			r.parsedHEPsMap[ep.UUID] = struct{}{}
			r.hepVersions[ep.UUID] = ep.Version
			r.parsedHEPs = append(r.parsedHEPs, entityToParsedHEP(ep))
		}
	}

	for _, set := range gnss {
		if !sel.Evaluate(set.Metadata.Labels) {
			continue
		}
		if !((ruleIPVersion == entity.IPVersion4 && len(set.Spec.NetsV4) > 0) || (ruleIPVersion == entity.IPVersion6 && len(set.Spec.NetsV6) > 0)) {
			continue
		}
		gnsUUIDs = append(gnsUUIDs, set.UUID)
		if _, ok := r.parsedGNSsMap[set.UUID]; !ok {
			r.parsedGNSsMap[set.UUID] = struct{}{}
			r.gnsVersions[set.UUID] = set.Version
			r.parsedGNSs = append(r.parsedGNSs, entityToParsedGNS(set))
		}
	}

	// if selector not match any hep and gns. Using match empty to prevent
	if len(gnsUUIDs) == 0 && len(hepUUIDs) == 0 {
		var set entity.GlobalNetworkSet
		if ruleIPVersion == entity.IPVersion4 {
			set = entity.GNSV4Empty
		} else if ruleIPVersion == entity.IPVersion6 {
			set = entity.GNSV6Empty
		}
		gnsUUIDs = append(gnsUUIDs, set.UUID)
		if _, ok := r.parsedGNSsMap[set.UUID]; !ok {
			r.parsedGNSsMap[set.UUID] = struct{}{}
			r.gnsVersions[set.UUID] = set.Version
			r.parsedGNSs = append(r.parsedGNSs, entityToParsedGNS(&set))
		}
	}

	return hepUUIDs, gnsUUIDs, nil
}

func entityToParsedHEP(hep *entity.HostEndpoint) *model.ParsedHEP {
	return &model.ParsedHEP{
		UUID:     hep.UUID,
		Name:     hep.Metadata.Name,
		TenantID: hep.Spec.TenantID,
		IP:       net.IntToIP(hep.Spec.IP).String(),
		IPsV4:    hep.Spec.IPsV4,
		IPsV6:    hep.Spec.IPsV6,
	}
}

func entityToParsedGNS(set *entity.GlobalNetworkSet) *model.ParsedGNS {
	return &model.ParsedGNS{
		UUID:   set.UUID,
		Name:   set.Metadata.Name,
		NetsV4: set.Spec.NetsV4,
		NetsV6: set.Spec.NetsV6,
	}
}

func convertPorts(ports []interface{}) []string {
	var portStrings []string
	for _, port := range ports {
		portStrings = append(portStrings, fmt.Sprint(port))
	}
	return portStrings
}
