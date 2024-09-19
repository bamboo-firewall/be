package mapper

import (
	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/domain/model"
)

func ToHostEndpointDTO(hep *entity.HostEndpoint) *dto.HostEndpoint {
	if hep == nil {
		return nil
	}
	var ports []dto.HostEndpointSpecPort
	for _, port := range hep.Spec.Ports {
		ports = append(ports, dto.HostEndpointSpecPort{
			Name:     port.Name,
			Port:     port.Port,
			Protocol: port.Protocol,
		})
	}
	return &dto.HostEndpoint{
		ID:      hep.ID.Hex(),
		UUID:    hep.UUID,
		Version: hep.Version,
		Metadata: dto.HostEndpointMetadata{
			Name:   hep.Metadata.Name,
			Labels: hep.Metadata.Labels,
		},
		Spec: dto.HostEndpointSpec{
			InterfaceName: hep.Spec.InterfaceName,
			IPs:           hep.Spec.IPs,
			Ports:         ports,
		},
		Description: hep.Description,
		CreatedAt:   hep.CreatedAt,
		UpdatedAt:   hep.UpdatedAt,
	}
}

func ToCreateHostEndpointInput(in *dto.CreateHostEndpointInput) *model.CreateHostEndpointInput {
	var ports []model.HostEndpointSpecPortInput
	for _, port := range in.Spec.Ports {
		ports = append(ports, model.HostEndpointSpecPortInput{
			Name:     port.Name,
			Port:     port.Port,
			Protocol: port.Protocol,
		})
	}

	return &model.CreateHostEndpointInput{
		Metadata: model.HostEndpointMetadataInput{
			Name:   in.Metadata.Name,
			Labels: in.Metadata.Labels,
		},
		Spec: model.HostEndpointSpecInput{
			InterfaceName: in.Spec.InterfaceName,
			IPs:           in.Spec.IPs,
			Ports:         ports,
		},
		Description: in.Description,
	}
}

func ToFetchPoliciesInput(in *dto.FetchPoliciesInput) *model.FetchPoliciesInput {
	return &model.FetchPoliciesInput{
		Name: in.Name,
	}
}

func ToFetchPoliciesOutput(hostEndpointPolicy *model.HostEndPointPolicy) *dto.FetchPoliciesOutput {
	policiesDTOs := make([]*dto.ParsedPolicy, len(hostEndpointPolicy.ParsedPolicies))
	for i, policy := range hostEndpointPolicy.ParsedPolicies {
		policiesDTOs[i] = ToParsedPolicyDTO(policy)
	}
	setDTOs := make([]*dto.ParsedSet, len(hostEndpointPolicy.ParsedSets))
	for i, set := range hostEndpointPolicy.ParsedSets {
		setDTOs[i] = ToParsedSetDTO(set)
	}
	return &dto.FetchPoliciesOutput{
		MetaData: dto.HostEndPointPolicyMetadata{
			HEPVersion:  hostEndpointPolicy.MetaData.HEPVersion,
			GNPVersions: hostEndpointPolicy.MetaData.GNPVersions,
			GNSVersions: hostEndpointPolicy.MetaData.GNSVersions,
		},
		HEP:            ToHostEndpointDTO(hostEndpointPolicy.HEP),
		ParsedPolicies: policiesDTOs,
		ParsedSets:     setDTOs,
	}
}

func ToParsedPolicyDTO(parsedPolicy *model.ParsedPolicy) *dto.ParsedPolicy {
	var inboundRules []*dto.ParsedRule
	for _, rule := range parsedPolicy.InboundRules {
		inboundRules = append(inboundRules, ToParsedRuleDTO(rule))
	}
	var outboundRules []*dto.ParsedRule
	for _, rule := range parsedPolicy.OutboundRules {
		outboundRules = append(outboundRules, ToParsedRuleDTO(rule))
	}
	return &dto.ParsedPolicy{
		UUID:          parsedPolicy.UUID,
		Version:       parsedPolicy.Version,
		Name:          parsedPolicy.Name,
		InboundRules:  inboundRules,
		OutboundRules: outboundRules,
	}
}

func ToParsedRuleDTO(parsedRule *model.ParsedRule) *dto.ParsedRule {
	return &dto.ParsedRule{
		Action:             parsedRule.Action,
		IPVersion:          parsedRule.IPVersion,
		Protocol:           parsedRule.Protocol,
		IsProtocolNegative: parsedRule.IsProtocolNegative,
		SrcNets:            parsedRule.SrcNets,
		IsSrcNetNegative:   parsedRule.IsSrcNetNegative,
		SrcGNSNetNames:     parsedRule.SrcGNSNetNames,
		SrcPorts:           parsedRule.SrcPorts,
		IsSrcPortNegative:  parsedRule.IsSrcPortNegative,
		DstNets:            parsedRule.DstNets,
		IsDstNetNegative:   parsedRule.IsDstPortNegative,
		DstGNSNetNames:     parsedRule.DstGNSNetNames,
		DstPorts:           parsedRule.DstPorts,
		IsDstPortNegative:  parsedRule.IsDstPortNegative,
	}
}

func ToParsedSetDTO(parsedSet *model.ParsedSet) *dto.ParsedSet {
	return &dto.ParsedSet{
		Name:      parsedSet.Name,
		IPVersion: parsedSet.IPVersion,
		Nets:      parsedSet.Nets,
	}
}
