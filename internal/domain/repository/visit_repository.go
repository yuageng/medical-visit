package repository

import (
	"errors"
	"time"

	"medical-visit/internal/domain/model"

	"github.com/google/uuid"
)

var (
	// ErrVisitStateConflict 表示拜访状态或版本已被其他请求修改。
	ErrVisitStateConflict = errors.New("visit state conflict")
	// ErrMaterialStateConflict 表示学术资料在保存报告前已失效。
	ErrMaterialStateConflict = errors.New("academic material state conflict")
)

// MonthlyProductStat 是数据库聚合后的产品月度统计结果。
type MonthlyProductStat struct {
	ProductID     uuid.UUID
	ProductName   string
	VisitCount    int64
	NormalCount   int64
	AbnormalCount int64
}

// VisitRepository 拜访仓储接口
type VisitRepository interface {
	// Create 创建拜访
	Create(visit *model.Visit) error

	// CheckIn 在期望状态和版本下原子写入签到信息。
	CheckIn(visitID uuid.UUID, expectedVersion int, visit *model.Visit) error

	// CheckOut 在期望状态和版本下原子写入签退及最终合规信息。
	CheckOut(visitID uuid.UUID, expectedVersion int, visit *model.Visit) error

	// SaveCallReport 在期望版本下创建或更新拜访记录，并同步资料关联及拜访状态。
	SaveCallReport(visitID uuid.UUID, expectedVersion int, report *model.CallReport, materialIDs []uuid.UUID) error

	// FindByID 根据 ID 查询拜访
	FindByID(id uuid.UUID) (*model.Visit, error)

	// FindByMRAndTimeRange 查询指定 MR 时间范围内的拜访。
	FindByMRAndTimeRange(mrID uuid.UUID, startTime, endTime string) ([]*model.Visit, error)

	// MonthlyProductStats 按签退月份查询产品拜访聚合统计。
	MonthlyProductStats(start, end time.Time) ([]MonthlyProductStat, error)

	// Update 更新拜访
	Update(visit *model.Visit) error
}

// MasterDataRepository 主数据仓储接口
type MasterDataRepository interface {
	// MR 相关
	FindMRByID(id uuid.UUID) (*model.MedicalRepresentative, error)
	FindAllMRs() ([]*model.MedicalRepresentative, error)

	// HCP 相关
	FindHCPByID(id uuid.UUID) (*model.HCP, error)
	FindHCPsByHospital(hospitalID uuid.UUID) ([]*model.HCP, error)

	// Hospital 相关
	FindHospitalByID(id uuid.UUID) (*model.Hospital, error)
	FindAllHospitals() ([]*model.Hospital, error)

	// Department 相关
	FindDepartmentByID(id uuid.UUID) (*model.Department, error)
	FindDepartmentsByHospital(hospitalID uuid.UUID) ([]*model.Department, error)

	// Product 相关
	FindProductByID(id uuid.UUID) (*model.Product, error)
	FindAllProducts() ([]*model.Product, error)

	// AcademicMaterial 相关
	FindAcademicMaterialsByIDs(ids []uuid.UUID) ([]*model.AcademicMaterial, error)
}
