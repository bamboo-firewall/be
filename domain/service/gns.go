package service

import (
	"context"
	"errors"
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
	"github.com/bamboo-firewall/be/domain/model"
)

func NewGNS(policyMongo *repository.PolicyDB) *gns {
	return &gns{
		storage: policyMongo,
	}
}

type gns struct {
	storage be.Storage
}

func (ds *gns) Create(ctx context.Context, input *model.CreateGlobalNetworkSetInput) (*entity.GlobalNetworkSet, *ierror.Error) {
	netsV4, netsV6 := exactNets(input.Spec.Nets)
	gnsEntity := &entity.GlobalNetworkSet{
		ID:   primitive.NewObjectID(),
		UUID: entity.NewMinifyUUID(),
		Metadata: entity.GNSMetadata{
			Name:   input.Metadata.Name,
			Labels: input.Metadata.Labels,
		},
		Spec: entity.GNSSpec{
			Nets:   input.Spec.Nets,
			NetsV4: netsV4,
			NetsV6: netsV6,
		},
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	gnsEntity, coreErr := ds.storage.UpsertGNS(ctx, gnsEntity)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrDuplicateGlobalNetworkSet) {
			return nil, httpbase.ErrBadRequest(ctx, "duplicate global network set").SetSubError(coreErr)
		}
		return nil, httpbase.ErrDatabase(ctx, "create global network set failed").SetSubError(coreErr)
	}
	return gnsEntity, nil
}

func exactNets(nets []string) (netsV4 []string, netsV6 []string) {
	for _, netString := range nets {
		ip, ipnet, err := net.ParseCIDROrIP(netString)
		if err != nil {
			slog.Warn("malformed net", "net", netString)
			continue
		}
		var netV4V6 string
		if ip.String() == ipnet.IP.String() {
			netV4V6 = ipnet.String()
		} else {
			netV4V6 = ip.Network().String()
		}
		if ip.Version() == entity.IPVersion4 {
			netsV4 = append(netsV4, netV4V6)
		} else if ip.Version() == entity.IPVersion6 {
			netsV6 = append(netsV6, netV4V6)
		}
	}
	return
}

func (ds *gns) Get(ctx context.Context, name string) (*entity.GlobalNetworkSet, *ierror.Error) {
	gnsEntity, coreErr := ds.storage.GetGNSByName(ctx, name)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrNotFoundGlobalNetworkSet) {
			return nil, httpbase.ErrNotFound(ctx, "not found").SetSubError(coreErr)
		}
		return nil, httpbase.ErrDatabase(ctx, "get global network set failed").SetSubError(coreErr)
	}
	return gnsEntity, nil
}

func (ds *gns) List(ctx context.Context) ([]*entity.GlobalNetworkSet, *ierror.Error) {
	gnssEntity, coreErr := ds.storage.ListGNSs(ctx)
	if coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "list global network sets failed").SetSubError(coreErr)
	}
	return gnssEntity, nil
}

func (ds *gns) Delete(ctx context.Context, name string) *ierror.Error {
	if coreErr := ds.storage.DeleteGNSByName(ctx, name); coreErr != nil {
		return httpbase.ErrDatabase(ctx, "delete global network set failed").SetSubError(coreErr)
	}
	return nil
}
