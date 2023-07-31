package handler

import (
	"net/http"

	"github.com/bamboo-firewall/be/domain"
	models "github.com/bamboo-firewall/watcher/model"
	"github.com/gin-gonic/gin"
)

type HEPHandler struct {
	HEPUsecase domain.HEPUsecase
}

func (hh *HEPHandler) convertToCalicoObjectResponse(heps []models.HostEndPoint) []domain.CalicoObjectResponse {
	response := make([]domain.CalicoObjectResponse, len(heps))
	for i, hep := range heps {
		response[i] = domain.CalicoObjectResponse{
			Kind:       hep.Kind,
			ApiVersion: hep.ApiVersion,
			Metadata:   hep.Metadata,
			Spec:       hep.Spec,
		}

	}
	return response
}

func (hh *HEPHandler) Fetch(c *gin.Context) {
	heps, err := hh.HEPUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, hh.convertToCalicoObjectResponse(heps))
}

func (hh *HEPHandler) Search(c *gin.Context) {
	var request domain.SearchRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	heps, err := hh.HEPUsecase.Search(c, request.Options)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, hh.convertToCalicoObjectResponse(heps))
}
