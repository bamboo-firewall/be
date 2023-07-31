package handler

import (
	"net/http"

	"github.com/bamboo-firewall/be/domain"
	"github.com/gin-gonic/gin"
)

type OptionHandler struct {
	HEPUsecase    domain.HEPUsecase
	GNSUsecase    domain.GNSUsecase
	PolicyUsecase domain.PolicyUsecase
}

func (hh *OptionHandler) FetchByType(c *gin.Context) {
	var request domain.FetchOptionRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	if !(request.Type == domain.CollectionHEP || request.Type == domain.CollectionPolicy || request.Type == domain.CollectionGNS) {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Type is not in valid!"})
	}
	var options []domain.Option

	switch request.Type {
	case domain.CollectionHEP:
		options, err = hh.HEPUsecase.GetOptions(c, request.Filter, request.Label)
	case domain.CollectionGNS:
		options, err = hh.GNSUsecase.GetOptions(c, request.Filter, request.Label)
	case domain.CollectionPolicy:
		options, err = hh.PolicyUsecase.GetOptions(c, request.Filter, request.Label)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, options)
}
