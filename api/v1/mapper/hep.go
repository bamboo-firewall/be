package mapper

import (
	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/net"
	"github.com/bamboo-firewall/be/domain/model"
)

func ToListHostEndpointDTOs(heps []*entity.HostEndpoint) []*dto.HostEndpoint {
	hepDTOs := make([]*dto.HostEndpoint, 0, len(heps))
	for _, hep := range heps {
		hepDTOs = append(hepDTOs, ToHostEndpointDTO(hep))
	}
	return hepDTOs
}

func ToHostEndpointDTO(hep *entity.HostEndpoint) *dto.HostEndpoint {
	if hep == nil {
		return nil
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
			TenantID:      hep.Spec.TenantID,
			IP:            net.IntToIP(hep.Spec.IP).String(),
			IPs:           hep.Spec.IPs,
		},
		Description: hep.Description,
		CreatedAt:   hep.CreatedAt.Local(),
		UpdatedAt:   hep.UpdatedAt.Local(),
	}
}

func ToCreateHostEndpointInput(in *dto.CreateHostEndpointInput) *model.CreateHostEndpointInput {
	return &model.CreateHostEndpointInput{
		Metadata: model.HostEndpointMetadataInput{
			Name:   in.Metadata.Name,
			Labels: in.Metadata.Labels,
		},
		Spec: model.HostEndpointSpecInput{
			InterfaceName: in.Spec.InterfaceName,
			IP:            in.Spec.IP,
			TenantID:      in.Spec.TenantID,
			IPs:           in.Spec.IPs,
		},
		Description: in.Description,
	}
}

func ToGetHostEndpointInput(in *dto.GetHostEndpointInput) *model.GetHostEndpointInput {
	var ipInt uint32
	netIP := net.ParseIP(in.IP)
	if netIP != nil {
		ipInt = net.IPToInt(*netIP)
	}
	return &model.GetHostEndpointInput{
		TenantID: in.TenantID,
		IP:       ipInt,
	}
}

func ToListHostEndpointsInput(in *dto.ListHostEndpointsInput) *model.ListHostEndpointsInput {
	var ipInt *uint32
	if in.IP != nil {
		netIP := net.ParseIP(*in.IP)
		if netIP != nil {
			ip := net.IPToInt(*netIP)
			ipInt = &ip
		}
	}
	return &model.ListHostEndpointsInput{
		TenantID: in.TenantID,
		IP:       ipInt,
	}
}

func ToFetchHostEndpointPolicyInput(in *dto.FetchHostEndpointPoliciesInput) *model.ListHostEndpointsInput {
	var ipInt *uint32
	if in.IP != nil {
		netIP := net.ParseIP(*in.IP)
		if netIP != nil {
			ip := net.IPToInt(*netIP)
			ipInt = &ip
		}
	}
	return &model.ListHostEndpointsInput{
		TenantID: in.TenantID,
		IP:       ipInt,
	}
}

func ToFetchHEPPoliciesOutput(hepPolicies []*model.HostEndpointPolicy) []*dto.HostEndpointPolicy {
	result := make([]*dto.HostEndpointPolicy, 0, len(hepPolicies))
	for _, hepPolicy := range hepPolicies {
		result = append(result, ToFetchHEPPolicyOutput(hepPolicy))
	}
	return result
}

func ToFetchHEPPolicyOutput(hostEndpointPolicy *model.HostEndpointPolicy) *dto.HostEndpointPolicy {
	parsedGNPDTOs := make([]*dto.ParsedGNP, len(hostEndpointPolicy.ParsedGNPs))
	for i, policy := range hostEndpointPolicy.ParsedGNPs {
		parsedGNPDTOs[i] = toParsedGNPDTO(policy)
	}
	parsedHEPDTOs := make([]*dto.ParsedHEP, len(hostEndpointPolicy.ParsedHEPs))
	for i, endpoint := range hostEndpointPolicy.ParsedHEPs {
		parsedHEPDTOs[i] = toParsedHEPDTO(endpoint)
	}
	parsedGNSDTOs := make([]*dto.ParsedGNS, len(hostEndpointPolicy.ParsedGNSs))
	for i, set := range hostEndpointPolicy.ParsedGNSs {
		parsedGNSDTOs[i] = toParsedGNSDTO(set)
	}
	return &dto.HostEndpointPolicy{
		MetaData: dto.HostEndpointPolicyMetadata{
			HEPVersions: hostEndpointPolicy.MetaData.HEPVersions,
			GNPVersions: hostEndpointPolicy.MetaData.GNPVersions,
			GNSVersions: hostEndpointPolicy.MetaData.GNSVersions,
		},
		HEP:        ToHostEndpointDTO(hostEndpointPolicy.HEP),
		ParsedGNPs: parsedGNPDTOs,
		ParsedHEPs: parsedHEPDTOs,
		ParsedGNSs: parsedGNSDTOs,
	}
}

func toParsedGNPDTO(parsedGNP *model.ParsedGNP) *dto.ParsedGNP {
	var inboundRules []*dto.ParsedRule
	for _, rule := range parsedGNP.InboundRules {
		inboundRules = append(inboundRules, toParsedRuleDTO(rule))
	}
	var outboundRules []*dto.ParsedRule
	for _, rule := range parsedGNP.OutboundRules {
		outboundRules = append(outboundRules, toParsedRuleDTO(rule))
	}
	return &dto.ParsedGNP{
		UUID:          parsedGNP.UUID,
		Version:       parsedGNP.Version,
		Name:          parsedGNP.Name,
		InboundRules:  inboundRules,
		OutboundRules: outboundRules,
	}
}

func toParsedRuleDTO(parsedRule *model.ParsedRule) *dto.ParsedRule {
	return &dto.ParsedRule{
		Action:             parsedRule.Action,
		IPVersion:          parsedRule.IPVersion,
		Protocol:           parsedRule.Protocol,
		IsProtocolNegative: parsedRule.IsProtocolNegative,
		SrcNets:            parsedRule.SrcNets,
		IsSrcNetNegative:   parsedRule.IsSrcNetNegative,
		SrcGNSUUIDs:        parsedRule.SrcGNSUUIDs,
		SrcHEPUUIDs:        parsedRule.SrcHEPUUIDs,
		SrcPorts:           parsedRule.SrcPorts,
		IsSrcPortNegative:  parsedRule.IsSrcPortNegative,
		DstNets:            parsedRule.DstNets,
		IsDstNetNegative:   parsedRule.IsDstPortNegative,
		DstGNSUUIDs:        parsedRule.DstGNSUUIDs,
		DstHEPUUIDs:        parsedRule.DstHEPUUIDs,
		DstPorts:           parsedRule.DstPorts,
		IsDstPortNegative:  parsedRule.IsDstPortNegative,
	}
}

func toParsedHEPDTO(parsedHEP *model.ParsedHEP) *dto.ParsedHEP {
	return &dto.ParsedHEP{
		UUID:  parsedHEP.UUID,
		Name:  parsedHEP.Name,
		IPsV4: parsedHEP.IPsV4,
		IPsV6: parsedHEP.IPsV6,
	}
}

func toParsedGNSDTO(parsedGNS *model.ParsedGNS) *dto.ParsedGNS {
	return &dto.ParsedGNS{
		UUID:   parsedGNS.UUID,
		Name:   parsedGNS.Name,
		NetsV4: parsedGNS.NetsV4,
		NetsV6: parsedGNS.NetsV6,
	}
}
