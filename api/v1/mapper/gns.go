package mapper

import (
	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/domain/model"
)

func ToListGlobalNetworkSetDTOs(gnss []*entity.GlobalNetworkSet) []*dto.GlobalNetworkSet {
	gnsDTOs := make([]*dto.GlobalNetworkSet, 0, len(gnss))
	for _, gns := range gnss {
		gnsDTOs = append(gnsDTOs, ToGlobalNetworkSetDTO(gns))
	}
	return gnsDTOs
}

func ToGlobalNetworkSetDTO(gns *entity.GlobalNetworkSet) *dto.GlobalNetworkSet {
	if gns == nil {
		return nil
	}
	return &dto.GlobalNetworkSet{
		ID:      gns.ID.Hex(),
		UUID:    gns.UUID,
		Version: gns.Version,
		Metadata: dto.GNSMetadata{
			Name:   gns.Metadata.Name,
			Labels: gns.Metadata.Labels,
		},
		Spec: dto.GNSSpec{
			Nets: gns.Spec.Nets,
		},
		Description: gns.Description,
		CreatedAt:   gns.CreatedAt.Local(),
		UpdatedAt:   gns.UpdatedAt.Local(),
	}
}

func ToCreateGlobalNetworkSetInput(in *dto.CreateGlobalNetworkSetInput) *model.CreateGlobalNetworkSetInput {
	return &model.CreateGlobalNetworkSetInput{
		Metadata: model.GNSMetadataInput{
			Name:   in.Metadata.Name,
			Labels: in.Metadata.Labels,
		},
		Spec: model.GNSSpecInput{
			Nets: in.Spec.Nets,
		},
		Description: in.Description,
	}
}
