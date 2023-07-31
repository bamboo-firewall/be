package domain

import "context"

type Summary struct {
	TotalGlobalNetworkSet int64 `json:"total_global_network_set"`
	TotalPolicy           int64 `json:"total_policy"`
	TotalHostEndpoint     int64 `json:"total_host_endpoint"`
	TotalUser             int64 `json:"total_user"`
}

type ProjectSummary struct {
	ProjectName string `json:"project_name"`
	Total       int64  `json:"total"`
}

type SummaryResponse struct {
	Summary Summary `json:"summary"`
}

type ProjectSummaryResponse struct {
	ProjectSummary []ProjectSummary `json:"project_summary"`
}

type StatisticUsecase interface {
	GetSummary(c context.Context) (Summary, error)
	GetProjectSummary(c context.Context) ([]ProjectSummary, error)
}
