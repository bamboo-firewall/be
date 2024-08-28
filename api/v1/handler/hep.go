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

type hepService interface {
	CreateHEP(ctx context.Context, input *model.CreateHostEndpointInput) (*entity.HostEndpoint, *ierror.Error)
	DeleteHEP(ctx context.Context, name string) *ierror.Error
}

func NewHEP(s hepService) *hep {
	return &hep{
		service: s,
	}
}

type hep struct {
	service hepService
}

func (h *hep) Create(c *gin.Context) {
	in := new(dto.CreateHostEndpointInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	hepEntity, ierr := h.service.CreateHEP(c.Request.Context(), mapper.ToCreateHostEndpointInput(in))
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToHostEndpointDTO(hepEntity))
}

func (h *hep) Delete(c *gin.Context) {
	in := new(dto.CreateHostEndpointInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	if err := h.service.DeleteHEP(c.Request.Context(), in.Metadata.Name); err != nil {
		httpbase.ReturnErrorResponse(c, err)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, nil)
}
