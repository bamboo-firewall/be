package mapper

import (
	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/domain/model"
)

func ToGlobalNetworkPolicyDTO(gnp *entity.GlobalNetworkPolicy) *dto.GlobalNetworkPolicy {
	if gnp == nil {
		return nil
	}

	var specIngress []dto.GNPSpecRule
	for _, rule := range gnp.Spec.Ingress {
		specIngress = append(specIngress, toRuleDTO(rule))
	}

	var specEgress []dto.GNPSpecRule
	for _, rule := range gnp.Spec.Egress {
		specEgress = append(specEgress, toRuleDTO(rule))
	}
	return &dto.GlobalNetworkPolicy{
		ID:      gnp.ID.Hex(),
		UUID:    gnp.UUID,
		Version: gnp.Version,
		Metadata: dto.GNPMetadata{
			Name:   gnp.Metadata.Name,
			Labels: gnp.Metadata.Labels,
		},
		Spec: dto.GNPSpec{
			Selector: gnp.Spec.Selector,
			Ingress:  specIngress,
			Egress:   specEgress,
		},
		Description: gnp.Description,
		CreatedAt:   gnp.CreatedAt.Local(),
		UpdatedAt:   gnp.UpdatedAt.Local(),
	}
}

func toRuleDTO(rule entity.GNPSpecRule) dto.GNPSpecRule {
	return dto.GNPSpecRule{
		Metadata:    rule.Metadata,
		Action:      rule.Action,
		Protocol:    rule.Protocol,
		NotProtocol: rule.NotProtocol,
		IPVersion:   int(rule.IPVersion),
		Source:      toRuleEntityDTO(rule.Source),
		Destination: toRuleEntityDTO(rule.Destination),
	}
}

func toRuleEntityDTO(ruleEntity *entity.GNPSpecRuleEntity) *dto.GNPSpecRuleEntity {
	if ruleEntity == nil {
		return nil
	}
	return &dto.GNPSpecRuleEntity{
		Selector: ruleEntity.Selector,
		Nets:     ruleEntity.Nets,
		NotNets:  ruleEntity.NotNets,
		Ports:    ruleEntity.Ports,
		NotPorts: ruleEntity.NotPorts,
	}
}

func ToCreateGlobalNetworkPolicyInput(in *dto.CreateGlobalNetworkPolicyInput) *model.CreateGlobalNetworkPolicyInput {
	var specIngress []model.GNPSpecRuleInput
	for _, rule := range in.Spec.Ingress {
		specIngress = append(specIngress, toRuleInput(rule))
	}

	var specEgress []model.GNPSpecRuleInput
	for _, rule := range in.Spec.Egress {
		specEgress = append(specEgress, toRuleInput(rule))
	}

	return &model.CreateGlobalNetworkPolicyInput{
		Metadata: model.GNPMetadataInput{
			Name:   in.Metadata.Name,
			Labels: in.Metadata.Labels,
		},
		Spec: model.GNPSpecInput{
			Selector: in.Spec.Selector,
			Ingress:  specIngress,
			Egress:   specEgress,
		},
		Description: in.Description,
	}
}

func toRuleInput(rule dto.GNPSpecRuleInput) model.GNPSpecRuleInput {
	return model.GNPSpecRuleInput{
		Metadata:    rule.Metadata,
		Action:      rule.Action,
		Protocol:    rule.Protocol,
		NotProtocol: rule.NotProtocol,
		IPVersion:   rule.IPVersion,
		Source:      toRuleEntityInput(rule.Source),
		Destination: toRuleEntityInput(rule.Destination),
	}
}

func toRuleEntityInput(ruleEntity *dto.GNPSpecRuleEntityInput) *model.GNPSpecRuleEntityInput {
	if ruleEntity == nil {
		return nil
	}
	return &model.GNPSpecRuleEntityInput{
		Selector: ruleEntity.Selector,
		Nets:     ruleEntity.Nets,
		NotNets:  ruleEntity.NotNets,
		Ports:    ruleEntity.Ports,
		NotPorts: ruleEntity.NotPorts,
	}
}
