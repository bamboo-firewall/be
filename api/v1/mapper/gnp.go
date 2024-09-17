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
			Types:    gnp.Spec.Types,
			Ingress:  specIngress,
			Egress:   specEgress,
		},
		Description: gnp.Description,
		CreatedAt:   gnp.CreatedAt,
		UpdatedAt:   gnp.UpdatedAt,
	}
}

func toRuleDTO(rule entity.GNPSpecRule) dto.GNPSpecRule {
	return dto.GNPSpecRule{
		Metadata:    rule.Metadata,
		Action:      rule.Action,
		Protocol:    rule.Protocol,
		NotProtocol: rule.NotProtocol,
		Source: dto.GNPSpecRuleEntity{
			Nets:     rule.Source.Nets,
			NotNets:  rule.Source.NotNets,
			Ports:    rule.Source.Ports,
			NotPorts: rule.Source.NotPorts,
		},
		Destination: dto.GNPSpecRuleEntity{
			Nets:     rule.Destination.Nets,
			NotNets:  rule.Destination.NotNets,
			Ports:    rule.Destination.Ports,
			NotPorts: rule.Destination.NotPorts,
		},
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
			Types:    in.Spec.Types,
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
		Source: model.GNPSpecRuleEntityInput{
			Nets:     rule.Source.Nets,
			NotNets:  rule.Source.NotNets,
			Ports:    rule.Source.Ports,
			NotPorts: rule.Source.NotPorts,
		},
		Destination: model.GNPSpecRuleEntityInput{
			Nets:     rule.Destination.Nets,
			NotNets:  rule.Destination.NotNets,
			Ports:    rule.Destination.Ports,
			NotPorts: rule.Destination.NotPorts,
		},
	}
}
