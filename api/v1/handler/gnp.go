package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/api/v1/mapper"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
	"github.com/bamboo-firewall/be/domain/model"
)

type gnpService interface {
	Create(ctx context.Context, input *model.CreateGlobalNetworkPolicyInput) (*entity.GlobalNetworkPolicy, *ierror.Error)
	Delete(ctx context.Context, name string) *ierror.Error
}

func NewGNP(s gnpService) *gnp {
	return &gnp{
		service: s,
	}
}

type gnp struct {
	service gnpService
}

func (h *gnp) Create(c *gin.Context) {
	in := new(dto.CreateGlobalNetworkPolicyInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	gnsEntity, ierr := h.service.Create(c.Request.Context(), mapper.ToCreateGlobalNetworkPolicyInput(in))
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToGlobalNetworkPolicyDTO(gnsEntity))
}

func (h *gnp) Delete(c *gin.Context) {
	in := new(dto.DeleteGlobalNetworkSetInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	if err := h.service.Delete(c.Request.Context(), in.Metadata.Name); err != nil {
		httpbase.ReturnErrorResponse(c, err)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, nil)
}