package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/api/v1/mapper"
	"github.com/bamboo-firewall/be/domain/model"
	"github.com/bamboo-firewall/be/pkg/entity"
	"github.com/bamboo-firewall/be/pkg/httpbase"
	"github.com/bamboo-firewall/be/pkg/httpbase/ierror"
)

type gnsService interface {
	Create(ctx context.Context, input *model.CreateGlobalNetworkSetInput) (*entity.GlobalNetworkSet, *ierror.Error)
	List(ctx context.Context) ([]*entity.GlobalNetworkSet, *ierror.Error)
	Get(ctx context.Context, name string) (*entity.GlobalNetworkSet, *ierror.Error)
	Delete(ctx context.Context, name string) *ierror.Error
	Validate(ctx context.Context, input *model.CreateGlobalNetworkSetInput) (*model.ValidateGlobalNetworkSetOutput, *ierror.Error)
}

func NewGNS(s gnsService) *gns {
	return &gns{
		service: s,
	}
}

type gns struct {
	service gnsService
}

func (h *gns) Create(c *gin.Context) {
	in := new(dto.CreateGlobalNetworkSetInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	gnsEntity, ierr := h.service.Create(c.Request.Context(), mapper.ToCreateGlobalNetworkSetInput(in))
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToGlobalNetworkSetDTO(gnsEntity))
}

func (h *gns) List(c *gin.Context) {
	gnpsEntity, ierr := h.service.List(c.Request.Context())
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToListGlobalNetworkSetDTOs(gnpsEntity))
}

func (h *gns) Get(c *gin.Context) {
	in := new(dto.GetGNSInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	gnsEntity, ierr := h.service.Get(c.Request.Context(), in.Name)
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToGlobalNetworkSetDTO(gnsEntity))
}

func (h *gns) Delete(c *gin.Context) {
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

func (h *gns) Validate(c *gin.Context) {
	in := new(dto.CreateGlobalNetworkSetInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	validateGlobalNetworkSetOutput, ierr := h.service.Validate(c.Request.Context(), mapper.ToCreateGlobalNetworkSetInput(in))
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToValidateGlobalNetworkSetOutput(validateGlobalNetworkSetOutput))
}
