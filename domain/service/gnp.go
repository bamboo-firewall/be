package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bamboo-firewall/be"
	"github.com/bamboo-firewall/be/cmd/server/pkg/common/errlist"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
	"github.com/bamboo-firewall/be/cmd/server/pkg/net"
	"github.com/bamboo-firewall/be/cmd/server/pkg/repository"
	"github.com/bamboo-firewall/be/domain/model"
)

func NewGNP(policyMongo *repository.PolicyDB) *gnp {
	return &gnp{
		storage: policyMongo,
	}
}

type gnp struct {
	storage be.Storage
}

func (ds *gnp) Create(ctx context.Context, input *model.CreateGlobalNetworkPolicyInput) (*entity.GlobalNetworkPolicy, *ierror.Error) {
	var specIngress []entity.GNPSpecRule
	for _, rule := range input.Spec.Ingress {
		specIngress = append(specIngress, modelToRule(rule))
	}

	var order uint32
	if input.Spec.Order != nil {
		order = *input.Spec.Order
	} else {
		order = entity.PolicyOrderLowest
	}
	var specEgress []entity.GNPSpecRule
	for _, rule := range input.Spec.Egress {
		specEgress = append(specEgress, modelToRule(rule))
	}

	gnpEntity := &entity.GlobalNetworkPolicy{
		ID:   primitive.NewObjectID(),
		UUID: entity.NewMinifyUUID(),
		Metadata: entity.GNPMetadata{
			Name:   input.Metadata.Name,
			Labels: input.Metadata.Labels,
		},
		Spec: entity.GNPSpec{
			Order:    order,
			Selector: input.Spec.Selector,
			Ingress:  specIngress,
			Egress:   specEgress,
		},
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	gnpEntity, coreErr := ds.storage.UpsertGroupPolicy(ctx, gnpEntity)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrDuplicateGlobalNetworkPolicy) {
			return nil, httpbase.ErrBadRequest(ctx, "duplicate global network policy").SetSubError(coreErr)
		}
		return nil, httpbase.ErrDatabase(ctx, "create global network policy failed").SetSubError(coreErr)
	}
	return gnpEntity, nil
}

func (ds *gnp) Get(ctx context.Context, name string) (*entity.GlobalNetworkPolicy, *ierror.Error) {
	gnpEntity, coreErr := ds.storage.GetGNPByName(ctx, name)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrNotFoundGlobalNetworkPolicy) {
			return nil, httpbase.ErrNotFound(ctx, "not found").SetSubError(coreErr)
		}
		return nil, httpbase.ErrDatabase(ctx, "get global network policy failed").SetSubError(coreErr)
	}
	return gnpEntity, nil
}

func (ds *gnp) Delete(ctx context.Context, name string) *ierror.Error {
	if coreErr := ds.storage.DeleteGNPByName(ctx, name); coreErr != nil {
		return httpbase.ErrDatabase(ctx, "delete global network policy failed").SetSubError(coreErr)
	}
	return nil
}

func (ds *gnp) List(ctx context.Context, input *model.ListGNPsInput) ([]*entity.GlobalNetworkPolicy, *ierror.Error) {
	gnpsEntity, coreErr := ds.storage.ListGNPs(ctx, input)
	if coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "list global network policies failed").SetSubError(coreErr)
	}
	return gnpsEntity, nil
}

func modelToRule(rule model.GNPSpecRuleInput) entity.GNPSpecRule {
	return entity.GNPSpecRule{
		Metadata:    rule.Metadata,
		Action:      rule.Action,
		Protocol:    rule.Protocol,
		NotProtocol: rule.NotProtocol,
		IPVersion:   rule.IPVersion,
		Source:      modelToRuleEntity(rule.Source),
		Destination: modelToRuleEntity(rule.Destination),
	}
}

func modelToRuleEntity(ruleEntity *model.GNPSpecRuleEntityInput) *entity.GNPSpecRuleEntity {
	if ruleEntity == nil {
		return nil
	}
	return &entity.GNPSpecRuleEntity{
		Selector: ruleEntity.Selector,
		Nets:     parseNets(ruleEntity.Nets),
		NotNets:  parseNets(ruleEntity.NotNets),
		Ports:    ruleEntity.Ports,
		NotPorts: ruleEntity.NotPorts,
	}
}

func parseNets(nets []string) []string {
	var netResults []string
	for _, n := range nets {
		_, ipnet, err := net.ParseCIDROrIP(n)
		if err == nil {
			netResults = append(netResults, ipnet.String())
		}
	}
	return netResults
}
