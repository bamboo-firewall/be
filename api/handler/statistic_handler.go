package handler

import (
	"net/http"

	"github.com/bamboo-firewall/be/domain"
	"github.com/gin-gonic/gin"
)

type StatisticHandler struct {
	StatisticUsecase domain.StatisticUsecase
}

func (hh *StatisticHandler) GetSummary(c *gin.Context) {
	summary, err := hh.StatisticUsecase.GetSummary(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SummaryResponse{
		Summary: summary,
	})
}

func (hh *StatisticHandler) GetProjectSummary(c *gin.Context) {
	projectSummary, err := hh.StatisticUsecase.GetProjectSummary(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.ProjectSummaryResponse{
		ProjectSummary: projectSummary,
	})
}
