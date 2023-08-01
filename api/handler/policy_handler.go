package handler

import (
	"net/http"

	"github.com/bamboo-firewall/be/domain"
	models "github.com/bamboo-firewall/watcher/model"
	"github.com/gin-gonic/gin"
)

type PolicyHandler struct {
	PolicyUsecase domain.PolicyUsecase
}

func (h *PolicyHandler) convertToCalicoObjectResponse(policies []models.GlobalNetworkPolicies) []domain.CalicoObjectResponse {
	response := make([]domain.CalicoObjectResponse, len(policies))
	for i, item := range policies {
		response[i] = domain.CalicoObjectResponse{
			Kind:       item.Kind,
			ApiVersion: item.ApiVersion,
			Metadata:   item.Metadata,
			Spec:       item.Spec,
		}

	}
	return response
}

func (h *PolicyHandler) Fetch(c *gin.Context) {
	items, err := h.PolicyUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, h.convertToCalicoObjectResponse(items))
}

func (hh *PolicyHandler) Search(c *gin.Context) {
	var request domain.SearchRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	policies, err := hh.PolicyUsecase.Search(c, request.Options)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, hh.convertToCalicoObjectResponse(policies))
}
