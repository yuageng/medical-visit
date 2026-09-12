package dto

import "medical-visit/internal/application"

// MonthlyVisitsResponse Dashboard 月度产品统计响应。
type MonthlyVisitsResponse struct {
	Month string                           `json:"month"`
	Items []application.MonthlyProductStat `json:"items"`
}
