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
		specIngress = append(specIngress, dto.GNPSpecRule{
			Metadata: rule.Metadata,
			Action:   rule.Action,
			Protocol: rule.Protocol,
			Source: dto.GNPSpecRuleEntity{
				Nets:  rule.Source.Nets,
				Ports: rule.Source.Ports,
			},
			Destination: dto.GNPSpecRuleEntity{
				Nets:  rule.Destination.Nets,
				Ports: rule.Destination.Ports,
			},
		})
	}

	var specEgress []dto.GNPSpecRule
	for _, rule := range gnp.Spec.Egress {
		specEgress = append(specEgress, dto.GNPSpecRule{
			Metadata: rule.Metadata,
			Action:   rule.Action,
			Protocol: rule.Protocol,
			Source: dto.GNPSpecRuleEntity{
				Nets:  rule.Source.Nets,
				Ports: rule.Source.Ports,
			},
			Destination: dto.GNPSpecRuleEntity{
				Nets:  rule.Destination.Nets,
				Ports: rule.Destination.Ports,
			},
		})
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

func ToCreateGlobalNetworkPolicyInput(in *dto.CreateGlobalNetworkPolicyInput) *model.CreateGlobalNetworkPolicyInput {
	var specIngress []model.GNPSpecRuleInput
	for _, rule := range in.Spec.Ingress {
		specIngress = append(specIngress, model.GNPSpecRuleInput{
			Metadata: rule.Metadata,
			Action:   rule.Action,
			Protocol: rule.Protocol,
			Source: model.GNPSpecRuleEntityInput{
				Nets:  rule.Source.Nets,
				Ports: rule.Source.Ports,
			},
			Destination: model.GNPSpecRuleEntityInput{
				Nets:  rule.Destination.Nets,
				Ports: rule.Destination.Ports,
			},
		})
	}

	var specEgress []model.GNPSpecRuleInput
	for _, rule := range in.Spec.Egress {
		specEgress = append(specEgress, model.GNPSpecRuleInput{
			Metadata: rule.Metadata,
			Action:   rule.Action,
			Protocol: rule.Protocol,
			Source: model.GNPSpecRuleEntityInput{
				Nets:  rule.Source.Nets,
				Ports: rule.Source.Ports,
			},
			Destination: model.GNPSpecRuleEntityInput{
				Nets:  rule.Destination.Nets,
				Ports: rule.Destination.Ports,
			},
		})
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
