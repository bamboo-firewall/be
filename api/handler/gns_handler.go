package handler

import (
	"net/http"

	"github.com/bamboo-firewall/be/domain"
	models "github.com/bamboo-firewall/watcher/model"
	"github.com/gin-gonic/gin"
)

type GNSHandler struct {
	GNSUsecase domain.GNSUsecase
}

func (hh *GNSHandler) convertToCalicoObjectResponse(gns []models.GlobalNetworkSet) []domain.CalicoObjectResponse {
	response := make([]domain.CalicoObjectResponse, len(gns))
	for i, item := range gns {
		response[i] = domain.CalicoObjectResponse{
			Kind:       item.Kind,
			ApiVersion: item.ApiVersion,
			Metadata:   item.Metadata,
			Spec:       item.Spec,
		}

	}
	return response
}

func (hh *GNSHandler) Fetch(c *gin.Context) {
	gns, err := hh.GNSUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, hh.convertToCalicoObjectResponse(gns))
}

func (hh *GNSHandler) Search(c *gin.Context) {
	var request domain.SearchRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	gns, err := hh.GNSUsecase.Search(c, request.Options)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, hh.convertToCalicoObjectResponse(gns))
}
