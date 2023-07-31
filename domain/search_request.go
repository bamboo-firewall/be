package domain

type SearchRequest struct {
	Options []Option `json:"options" binding:"required"`
}
