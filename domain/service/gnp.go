package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bamboo-firewall/be"
	"github.com/bamboo-firewall/be/cmd/server/pkg/common/errlist"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
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
	// ToDo: use transaction and lock row
	gnpExisted, coreErr := ds.storage.GetGNPByName(ctx, input.Metadata.Name)
	if coreErr != nil && !errors.Is(coreErr, errlist.ErrNotFoundGlobalNetworkPolicy) {
		return nil, httpbase.ErrDatabase(ctx, "get global network policy failed").SetSubError(coreErr)
	}

	var specIngress []entity.GNPSpecRule
	for _, rule := range input.Spec.Ingress {
		specIngress = append(specIngress, toRuleEntity(rule))
	}

	var specEgress []entity.GNPSpecRule
	for _, rule := range input.Spec.Egress {
		specEgress = append(specEgress, toRuleEntity(rule))
	}

	gnpEntity := &entity.GlobalNetworkPolicy{
		ID:      primitive.NewObjectID(),
		UUID:    uuid.New().String(),
		Version: 1,
		Metadata: entity.GNPMetadata{
			Name:   input.Metadata.Name,
			Labels: input.Metadata.Labels,
		},
		Spec: entity.GNPSpec{
			Selector: input.Spec.Selector,
			Types:    input.Spec.Types,
			Ingress:  specIngress,
			Egress:   specEgress,
		},
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if !errors.Is(coreErr, errlist.ErrNotFoundGlobalNetworkPolicy) {
		gnpEntity.ID = gnpExisted.ID
		gnpEntity.UUID = gnpExisted.UUID
		gnpEntity.Version = gnpExisted.Version + 1
		gnpEntity.CreatedAt = gnpExisted.CreatedAt
	}

	if coreErr = ds.storage.UpsertGroupPolicy(ctx, gnpEntity); coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "Create global network failed").SetSubError(coreErr)
	}
	return gnpEntity, nil
}

func (ds *gnp) Delete(ctx context.Context, name string) *ierror.Error {
	if coreErr := ds.storage.DeleteGNPByName(ctx, name); coreErr != nil {
		return httpbase.ErrDatabase(ctx, "Delete global network policy failed").SetSubError(coreErr)
	}
	return nil
}

func toRuleEntity(rule model.GNPSpecRuleInput) entity.GNPSpecRule {
	return entity.GNPSpecRule{
		Metadata:    rule.Metadata,
		Action:      rule.Action,
		Protocol:    rule.Protocol,
		NotProtocol: rule.NotProtocol,
		IPVersion:   rule.IPVersion,
		Source: entity.GNPSpecRuleEntity{
			Selector: rule.Source.Selector,
			Nets:     rule.Source.Nets,
			NotNets:  rule.Source.NotNets,
			Ports:    rule.Source.Ports,
			NotPorts: rule.Source.NotPorts,
		},
		Destination: entity.GNPSpecRuleEntity{
			Selector: rule.Destination.Selector,
			Nets:     rule.Destination.Nets,
			NotNets:  rule.Destination.NotNets,
			Ports:    rule.Destination.Ports,
			NotPorts: rule.Destination.NotPorts,
		},
	}
}
