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
		Name:    in.Name,
		Version: in.Version,
	}
}

func ToFetchPoliciesOutput(hep *entity.HostEndpoint, policies []*entity.GlobalNetworkPolicy, sets []*entity.GlobalNetworkSet) *dto.FetchPoliciesOutput {
	policiesDTOs := make([]*dto.GlobalNetworkPolicy, len(policies))
	for i, policy := range policies {
		policiesDTOs[i] = ToGlobalNetworkPolicyDTO(policy)
	}
	setDTOs := make([]*dto.GlobalNetworkSet, len(sets))
	for i, set := range sets {
		setDTOs[i] = ToGlobalNetworkSetDTO(set)
	}
	return &dto.FetchPoliciesOutput{
		IsNew:        true,
		HostEndpoint: ToHostEndpointDTO(hep),
		GNPs:         policiesDTOs,
		GNSs:         setDTOs,
	}
}
