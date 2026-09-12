package persistence

import (
	"errors"
	"fmt"
	"time"

	"medical-visit/internal/domain/model"
	"medical-visit/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// visitRepositoryImpl 拜访仓储实现
type visitRepositoryImpl struct {
	db *gorm.DB
}

// NewVisitRepository 创建拜访仓储实例
func NewVisitRepository(db *gorm.DB) repository.VisitRepository {
	return &visitRepositoryImpl{db: db}
}

// Create 创建拜访
func (r *visitRepositoryImpl) Create(visit *model.Visit) error {
	if err := r.db.Create(visit).Error; err != nil {
		return fmt.Errorf("failed to create visit: %w", err)
	}
	return nil
}

// CheckIn 使用状态和版本条件更新，避免并发重复签到覆盖首次签到事实。
func (r *visitRepositoryImpl) CheckIn(visitID uuid.UUID, expectedVersion int, visit *model.Visit) error {
	updates := map[string]interface{}{
		"status":              visit.Status,
		"checked_in_at":       visit.CheckedInAt,
		"check_in_latitude":   visit.CheckInLatitude,
		"check_in_longitude":  visit.CheckInLongitude,
		"check_in_distance_m": visit.CheckInDistanceM,
		"compliance_status":   visit.ComplianceStatus,
		"anomaly_reasons":     visit.AnomalyReasons,
		"version":             expectedVersion + 1,
		"updated_at":          visit.UpdatedAt,
	}
	result := r.db.Model(&model.Visit{}).
		Where("id = ? AND status = ? AND version = ?", visitID, model.VisitStatusPlanned, expectedVersion).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to check in visit: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return repository.ErrVisitStateConflict
	}
	return nil
}

// SaveCallReport 在单个事务中保存报告、资料关联并推进拜访状态。
func (r *visitRepositoryImpl) SaveCallReport(visitID uuid.UUID, report *model.CallReport, materialIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing model.CallReport
		err := tx.Where("visit_id = ?", visitID).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			report.VisitID = visitID
			if err := tx.Create(report).Error; err != nil {
				return fmt.Errorf("failed to create call report: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to find call report: %w", err)
		} else {
			report.ID = existing.ID
			report.VisitID = visitID
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"discussion_summary":    report.DiscussionSummary,
				"hcp_feedback":          report.HCPFeedback,
				"materials_distributed": report.MaterialsDistributed,
				"material_note":         report.MaterialNote,
				"additional_notes":      report.AdditionalNotes,
				"updated_at":            report.UpdatedAt,
			}).Error; err != nil {
				return fmt.Errorf("failed to update call report: %w", err)
			}
		}

		if err := tx.Exec("DELETE FROM call_report_materials WHERE call_report_id = ?", report.ID).Error; err != nil {
			return fmt.Errorf("failed to clear report materials: %w", err)
		}
		for _, materialID := range materialIDs {
			if err := tx.Exec("INSERT INTO call_report_materials (call_report_id, material_id) VALUES (?, ?)", report.ID, materialID).Error; err != nil {
				return fmt.Errorf("failed to save report material: %w", err)
			}
		}

		updates := map[string]interface{}{"status": model.VisitStatusCompleted, "version": gorm.Expr("version + 1"), "updated_at": report.UpdatedAt}
		result := tx.Model(&model.Visit{}).Where("id = ? AND status IN ?", visitID, []model.VisitStatus{model.VisitStatusCheckedOut, model.VisitStatusCompleted}).Updates(updates)
		if result.Error != nil {
			return fmt.Errorf("failed to complete visit: %w", result.Error)
		}
		if result.RowsAffected != 1 {
			return repository.ErrVisitStateConflict
		}
		return nil
	})
}

// MonthlyProductStats 在数据库中完成产品维度的月度条件聚合。
func (r *visitRepositoryImpl) MonthlyProductStats(start, end time.Time) ([]repository.MonthlyProductStat, error) {
	var stats []repository.MonthlyProductStat
	err := r.db.Table("visits AS v").
		Select("v.product_id, p.name AS product_name, COUNT(*) AS visit_count, COUNT(*) FILTER (WHERE v.compliance_status = ?) AS normal_count, COUNT(*) FILTER (WHERE v.compliance_status = ?) AS abnormal_count", model.ComplianceStatusCompliant, model.ComplianceStatusNonCompliant).
		Joins("JOIN products AS p ON p.id = v.product_id").
		Where("v.checked_out_at >= ? AND v.checked_out_at < ? AND v.status IN ?", start, end, []model.VisitStatus{model.VisitStatusCheckedOut, model.VisitStatusCompleted}).
		Group("v.product_id, p.name").
		Order("visit_count DESC, p.name ASC").
		Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly product stats: %w", err)
	}
	return stats, nil
}

// FindByID 根据 ID 查询拜访
func (r *visitRepositoryImpl) FindByID(id uuid.UUID) (*model.Visit, error) {
	var visit model.Visit
	err := r.db.
		Preload("MR").
		Preload("HCP").
		Preload("Hospital").
		Preload("Department").
		Preload("Product").
		Preload("CallReport").
		First(&visit, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}

	return &visit, nil
}

// Update 更新拜访
func (r *visitRepositoryImpl) Update(visit *model.Visit) error {
	if err := r.db.Save(visit).Error; err != nil {
		return fmt.Errorf("failed to update visit: %w", err)
	}
	return nil
}

// FindByMRAndTimeRange 查询某 MR 在指定时间范围内的拜访
func (r *visitRepositoryImpl) FindByMRAndTimeRange(mrID uuid.UUID, startTime, endTime string) ([]*model.Visit, error) {
	var visits []*model.Visit
	err := r.db.
		Where("mr_id = ? AND planned_start_at >= ? AND planned_start_at < ?", mrID, startTime, endTime).
		Order("planned_start_at ASC").
		Find(&visits).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find visits by MR and time range: %w", err)
	}

	return visits, nil
}
