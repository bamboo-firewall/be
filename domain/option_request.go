package domain

type FetchOptionRequest struct {
	Type   string   `json:"type" binding:"required"`
	Label  string   `json:"label" binding:"required"`
	Filter []Option `json:"filter" binding:"required"`
}
