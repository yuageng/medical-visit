package dto

import (
	"time"

	"medical-visit/internal/domain/model"

	"github.com/google/uuid"
)

// CreateVisitRequest 创建拜访请求
type CreateVisitRequest struct {
	MRID           string    `json:"mr_id" binding:"required,uuid"`
	HCPID          string    `json:"hcp_id" binding:"required,uuid"`
	HospitalID     string    `json:"hospital_id" binding:"required,uuid"`
	DepartmentID   string    `json:"department_id" binding:"required,uuid"`
	ProductID      string    `json:"product_id" binding:"required,uuid"`
	PlannedStartAt time.Time `json:"planned_start_at" binding:"required"`
	PlanNote       string    `json:"plan_note" binding:"max=5000"`
}

// CheckInRequest 签到请求。
// 使用指针保留 0 经度/纬度，并能区分缺少字段和合法零值。
type CheckInRequest struct {
	Latitude    *float64  `json:"latitude" binding:"required"`
	Longitude   *float64  `json:"longitude" binding:"required"`
	CheckInTime time.Time `json:"check_in_time" binding:"required"`
}

// VisitResponse 拜访响应
type VisitResponse struct {
	ID               string                  `json:"id"`
	MRID             string                  `json:"mr_id"`
	MRName           string                  `json:"mr_name,omitempty"`
	HCPID            string                  `json:"hcp_id"`
	HCPName          string                  `json:"hcp_name,omitempty"`
	HCPTitle         string                  `json:"hcp_title,omitempty"`
	HospitalID       string                  `json:"hospital_id"`
	HospitalName     string                  `json:"hospital_name,omitempty"`
	DepartmentID     string                  `json:"department_id"`
	DepartmentName   string                  `json:"department_name,omitempty"`
	ProductID        string                  `json:"product_id"`
	ProductName      string                  `json:"product_name,omitempty"`
	PlannedStartAt   time.Time               `json:"planned_start_at"`
	PlanNote         string                  `json:"plan_note,omitempty"`
	Status           string                  `json:"status"`
	CheckedInAt      *time.Time              `json:"checked_in_at,omitempty"`
	CheckInLatitude  *float64                `json:"check_in_latitude,omitempty"`
	CheckInLongitude *float64                `json:"check_in_longitude,omitempty"`
	CheckInDistanceM *float64                `json:"check_in_distance_m,omitempty"`
	CheckedOutAt     *time.Time              `json:"checked_out_at,omitempty"`
	DurationSeconds  *int                    `json:"duration_seconds,omitempty"`
	ComplianceStatus string                  `json:"compliance_status"`
	AnomalyReasons   []AnomalyReasonResponse `json:"anomaly_reasons,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

// AnomalyReasonResponse 异常原因响应
type AnomalyReasonResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ToVisitResponse 将 Visit 实体转换为响应 DTO
func ToVisitResponse(visit *model.Visit) *VisitResponse {
	resp := &VisitResponse{
		ID:               visit.ID.String(),
		MRID:             visit.MRID.String(),
		HCPID:            visit.HCPID.String(),
		HospitalID:       visit.HospitalID.String(),
		DepartmentID:     visit.DepartmentID.String(),
		ProductID:        visit.ProductID.String(),
		PlannedStartAt:   visit.PlannedStartAt,
		PlanNote:         visit.PlanNote,
		Status:           visit.Status.String(),
		CheckedInAt:      visit.CheckedInAt,
		CheckInLatitude:  visit.CheckInLatitude,
		CheckInLongitude: visit.CheckInLongitude,
		CheckInDistanceM: visit.CheckInDistanceM,
		CheckedOutAt:     visit.CheckedOutAt,
		DurationSeconds:  visit.DurationSeconds,
		ComplianceStatus: visit.ComplianceStatus.String(),
		CreatedAt:        visit.CreatedAt,
		UpdatedAt:        visit.UpdatedAt,
	}

	// 转换异常原因
	if len(visit.AnomalyReasons) > 0 {
		resp.AnomalyReasons = make([]AnomalyReasonResponse, len(visit.AnomalyReasons))
		for i, reason := range visit.AnomalyReasons {
			resp.AnomalyReasons[i] = AnomalyReasonResponse{
				Code:    reason.String(),
				Message: reason.Message(),
			}
		}
	}

	// 如果关联数据已加载，填充名称
	if visit.MR != nil {
		resp.MRName = visit.MR.Name
	}
	if visit.HCP != nil {
		resp.HCPName = visit.HCP.Name
		resp.HCPTitle = visit.HCP.Title
	}
	if visit.Hospital != nil {
		resp.HospitalName = visit.Hospital.Name
	}
	if visit.Department != nil {
		resp.DepartmentName = visit.Department.Name
	}
	if visit.Product != nil {
		resp.ProductName = visit.Product.Name
	}

	return resp
}

// ParseUUID 解析 UUID 字符串
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
