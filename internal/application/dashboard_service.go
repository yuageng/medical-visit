package application

import (
	"fmt"
	"time"

	"medical-visit/internal/domain/repository"
)

// MonthlyProductStat 月度产品拜访统计。
type MonthlyProductStat struct {
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	VisitCount    int64  `json:"visit_count"`
	NormalCount   int64  `json:"normal_count"`
	AbnormalCount int64  `json:"abnormal_count"`
}

// DashboardService Dashboard 应用服务。
type DashboardService struct {
	visitRepo repository.VisitRepository
}

func NewDashboardService(visitRepo repository.VisitRepository) *DashboardService {
	return &DashboardService{visitRepo: visitRepo}
}

// MonthlyVisitsByProduct 按 UTC 自然月查询产品拜访统计。
func (s *DashboardService) MonthlyVisitsByProduct(month string) ([]MonthlyProductStat, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("invalid month, expected YYYY-MM: %w", err)
	}
	end := start.AddDate(0, 1, 0)

	stats, err := s.visitRepo.MonthlyProductStats(start.UTC(), end.UTC())
	if err != nil {
		return nil, err
	}
	result := make([]MonthlyProductStat, len(stats))
	for i, stat := range stats {
		result[i] = MonthlyProductStat{
			ProductID: stat.ProductID.String(), ProductName: stat.ProductName,
			VisitCount: stat.VisitCount, NormalCount: stat.NormalCount, AbnormalCount: stat.AbnormalCount,
		}
	}
	return result, nil
}
