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
	Create(ctx context.Context, input *model.CreateHostEndpointInput) (*entity.HostEndpoint, *ierror.Error)
	List(ctx context.Context, input *model.ListHostEndpointsInput) ([]*entity.HostEndpoint, *ierror.Error)
	Get(ctx context.Context, input *model.GetHostEndpointInput) (*entity.HostEndpoint, *ierror.Error)
	Delete(ctx context.Context, input *model.DeleteHostEndpointInput) *ierror.Error
	FetchPolicies(ctx context.Context, input *model.ListHostEndpointsInput) ([]*model.HostEndpointPolicy, *ierror.Error)
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

	hepEntity, ierr := h.service.Create(c.Request.Context(), mapper.ToCreateHostEndpointInput(in))
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToHostEndpointDTO(hepEntity))
}

func (h *hep) List(c *gin.Context) {
	in := new(dto.ListHostEndpointsInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	gnpsEntity, ierr := h.service.List(c.Request.Context(), mapper.ToListHostEndpointsInput(in))
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToListHostEndpointDTOs(gnpsEntity))
}

func (h *hep) Get(c *gin.Context) {
	in := new(dto.GetHostEndpointInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	hepEntity, ierr := h.service.Get(c.Request.Context(), mapper.ToGetHostEndpointInput(in))
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToHostEndpointDTO(hepEntity))
}

func (h *hep) Delete(c *gin.Context) {
	in := new(dto.DeleteHostEndpointInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}

	if err := h.service.Delete(c.Request.Context(), &model.DeleteHostEndpointInput{
		TenantID: in.Spec.TenantID,
		IP:       in.Spec.IP,
		IPs:      in.Spec.IPs,
	}); err != nil {
		httpbase.ReturnErrorResponse(c, err)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, nil)
}

func (h *hep) FetchPolicies(c *gin.Context) {
	in := new(dto.FetchHostEndpointPoliciesInput)
	if ierr := httpbase.BindInput(c, in); ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	hostEndpointPolicies, ierr := h.service.FetchPolicies(c.Request.Context(), mapper.ToFetchHostEndpointPolicyInput(in))
	if ierr != nil {
		httpbase.ReturnErrorResponse(c, ierr)
		return
	}
	httpbase.ReturnSuccessResponse(c, http.StatusOK, mapper.ToFetchHEPPoliciesOutput(hostEndpointPolicies))
}
