package application

import (
	"errors"
	"fmt"
	"time"

	"medical-visit/internal/domain/compliance"
	"medical-visit/internal/domain/model"
	"medical-visit/internal/domain/repository"
	"medical-visit/internal/infrastructure/clock"

	"github.com/google/uuid"
)

var (
	// ErrMRNotFound MR 不存在
	ErrMRNotFound = errors.New("medical representative not found")
	// ErrHCPNotFound 医生不存在
	ErrHCPNotFound = errors.New("HCP not found")
	// ErrHospitalNotFound 医院不存在
	ErrHospitalNotFound = errors.New("hospital not found")
	// ErrDepartmentNotFound 科室不存在
	ErrDepartmentNotFound = errors.New("department not found")
	// ErrProductNotFound 产品不存在
	ErrProductNotFound = errors.New("product not found")
	// ErrDepartmentNotBelongToHospital 科室不属于该医院
	ErrDepartmentNotBelongToHospital = errors.New("department does not belong to the specified hospital")
	// ErrHCPNotInHospital HCP 不在该医院
	ErrHCPNotInHospital = errors.New("HCP is not in the specified hospital")
	// ErrHCPNotInDepartment HCP 不在该科室
	ErrHCPNotInDepartment = errors.New("HCP is not in the specified department")
	// ErrPlannedTimeInPast 计划时间在过去
	ErrPlannedTimeInPast = errors.New("planned start time is in the past")
	// ErrVisitNotFound 拜访不存在
	ErrVisitNotFound = errors.New("visit not found")
	// ErrVisitCannotCheckIn 拜访当前状态不允许签到
	ErrVisitCannotCheckIn = errors.New("visit status does not allow check-in")
	// ErrHospitalCoordinatesMissing 医院没有完整 GPS 坐标
	ErrHospitalCoordinatesMissing = errors.New("hospital coordinates are missing")
	// ErrInvalidCheckInTime 签到时间非法
	ErrInvalidCheckInTime = errors.New("check-in time is invalid")
	// ErrInvalidCheckInGPS 签到 GPS 非法
	ErrInvalidCheckInGPS = errors.New("check-in GPS is invalid")
	// ErrVisitCannotCheckOut 拜访当前状态不允许签退
	ErrVisitCannotCheckOut = errors.New("visit status does not allow check-out")
	// ErrInvalidCheckOutTime 签退时间为空或早于签到时间
	ErrInvalidCheckOutTime = errors.New("check-out time is invalid")
	// ErrInvalidCheckOutGPS 签退 GPS 非法
	ErrInvalidCheckOutGPS = errors.New("check-out GPS is invalid")
)

// VisitService 拜访应用服务
type VisitService struct {
	visitRepo      repository.VisitRepository
	masterDataRepo repository.MasterDataRepository
	clock          clock.Clock
	compliance     *compliance.Service
}

// NewVisitService 创建拜访应用服务实例
func NewVisitService(
	visitRepo repository.VisitRepository,
	masterDataRepo repository.MasterDataRepository,
	clk clock.Clock,
	complianceServices ...*compliance.Service,
) *VisitService {
	complianceService, err := compliance.NewService(compliance.DefaultRules())
	if len(complianceServices) > 0 && complianceServices[0] != nil {
		complianceService = complianceServices[0]
	}
	if err != nil {
		panic(err)
	}
	return &VisitService{
		visitRepo:      visitRepo,
		masterDataRepo: masterDataRepo,
		clock:          clk,
		compliance:     complianceService,
	}
}

// CreateVisitInput 创建拜访输入
type CreateVisitInput struct {
	MRID           uuid.UUID
	HCPID          uuid.UUID
	HospitalID     uuid.UUID
	DepartmentID   uuid.UUID
	ProductID      uuid.UUID
	PlannedStartAt time.Time
	PlanNote       string
}

// CreateVisit 创建拜访
func (s *VisitService) CreateVisit(input CreateVisitInput) (*model.Visit, error) {
	// 1. 校验所有基础数据存在
	mr, err := s.masterDataRepo.FindMRByID(input.MRID)
	if err != nil {
		return nil, fmt.Errorf("failed to find MR: %w", err)
	}
	if mr == nil {
		return nil, ErrMRNotFound
	}

	hcp, err := s.masterDataRepo.FindHCPByID(input.HCPID)
	if err != nil {
		return nil, fmt.Errorf("failed to find HCP: %w", err)
	}
	if hcp == nil {
		return nil, ErrHCPNotFound
	}

	hospital, err := s.masterDataRepo.FindHospitalByID(input.HospitalID)
	if err != nil {
		return nil, fmt.Errorf("failed to find hospital: %w", err)
	}
	if hospital == nil {
		return nil, ErrHospitalNotFound
	}

	department, err := s.masterDataRepo.FindDepartmentByID(input.DepartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find department: %w", err)
	}
	if department == nil {
		return nil, ErrDepartmentNotFound
	}

	product, err := s.masterDataRepo.FindProductByID(input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// 2. 校验科室属于医院
	if department.HospitalID != input.HospitalID {
		return nil, ErrDepartmentNotBelongToHospital
	}

	// 3. 校验 HCP 与医院、科室一致
	if hcp.HospitalID != input.HospitalID {
		return nil, ErrHCPNotInHospital
	}
	if hcp.DepartmentID != input.DepartmentID {
		return nil, ErrHCPNotInDepartment
	}

	// 4. 校验计划时间（允许当前时间或未来时间，不允许明显早于当前时间）
	now := s.clock.Now()
	// 允许 5 分钟的时间误差，避免客户端时间略有偏差导致创建失败
	if input.PlannedStartAt.Before(now.Add(-5 * time.Minute)) {
		return nil, ErrPlannedTimeInPast
	}

	// 5. 创建拜访实体
	visit := &model.Visit{
		ID:               uuid.New(),
		MRID:             input.MRID,
		HCPID:            input.HCPID,
		HospitalID:       input.HospitalID,
		DepartmentID:     input.DepartmentID,
		ProductID:        input.ProductID,
		PlannedStartAt:   input.PlannedStartAt,
		PlanNote:         input.PlanNote,
		Status:           model.VisitStatusPlanned,
		ComplianceStatus: model.ComplianceStatusPending,
		AnomalyReasons:   model.AnomalyReasons{},
		Version:          1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// 6. 持久化
	if err := s.visitRepo.Create(visit); err != nil {
		return nil, fmt.Errorf("failed to create visit: %w", err)
	}

	// 7. 重新加载以获取关联数据
	createdVisit, err := s.visitRepo.FindByID(visit.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload visit: %w", err)
	}

	return createdVisit, nil
}

// CheckInInput 签到输入。
type CheckInInput struct {
	VisitID     uuid.UUID
	Latitude    float64
	Longitude   float64
	CheckInTime time.Time
}

// CheckIn 执行拜访签到。
func (s *VisitService) CheckIn(input CheckInInput) (*model.Visit, error) {
	visit, err := s.visitRepo.FindByID(input.VisitID)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}
	if visit.Status != model.VisitStatusPlanned {
		return nil, ErrVisitCannotCheckIn
	}
	if input.CheckInTime.IsZero() {
		return nil, ErrInvalidCheckInTime
	}
	if err := (compliance.Coordinate{Latitude: input.Latitude, Longitude: input.Longitude}).Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCheckInGPS, err)
	}
	if visit.Hospital == nil || visit.Hospital.Latitude == nil || visit.Hospital.Longitude == nil {
		return nil, ErrHospitalCoordinatesMissing
	}
	hospitalCoordinate := compliance.Coordinate{Latitude: *visit.Hospital.Latitude, Longitude: *visit.Hospital.Longitude}
	checkInCoordinate := compliance.Coordinate{Latitude: input.Latitude, Longitude: input.Longitude}
	distance, err := compliance.HaversineDistanceMeters(checkInCoordinate, hospitalCoordinate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate check-in distance: %w", err)
	}
	result, err := s.compliance.ValidateCheckIn(compliance.CheckInInput{
		CheckInTime:      input.CheckInTime,
		CheckInDistanceM: distance,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate check-in: %w", err)
	}
	checkedInAt := input.CheckInTime.UTC()
	visit.CheckedInAt = &checkedInAt
	visit.CheckInLatitude = &input.Latitude
	visit.CheckInLongitude = &input.Longitude
	visit.CheckInDistanceM = &distance
	visit.Status = model.VisitStatusCheckedIn
	visit.ComplianceStatus = result.Status
	visit.AnomalyReasons = result.AnomalyReasons
	visit.UpdatedAt = s.clock.Now().UTC()
	if err := s.visitRepo.CheckIn(visit.ID, visit.Version, visit); err != nil {
		if errors.Is(err, repository.ErrVisitStateConflict) {
			return nil, ErrVisitCannotCheckIn
		}
		return nil, err
	}
	return s.visitRepo.FindByID(visit.ID)
}

// CheckOutInput 签退输入。
type CheckOutInput struct {
	VisitID      uuid.UUID
	Latitude     float64
	Longitude    float64
	CheckOutTime time.Time
}

// CheckOut 仅允许已签到拜访签退，并原子保存最终合规结果。
func (s *VisitService) CheckOut(input CheckOutInput) (*model.Visit, error) {
	visit, err := s.visitRepo.FindByID(input.VisitID)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}
	if visit.Status != model.VisitStatusCheckedIn || visit.CheckedInAt == nil {
		return nil, ErrVisitCannotCheckOut
	}
	if input.CheckOutTime.IsZero() || input.CheckOutTime.Before(*visit.CheckedInAt) {
		return nil, ErrInvalidCheckOutTime
	}
	coordinate := compliance.Coordinate{Latitude: input.Latitude, Longitude: input.Longitude}
	if err := coordinate.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCheckOutGPS, err)
	}
	if visit.Hospital == nil || visit.Hospital.Latitude == nil || visit.Hospital.Longitude == nil {
		return nil, ErrHospitalCoordinatesMissing
	}
	distance, err := compliance.HaversineDistanceMeters(coordinate, compliance.Coordinate{
		Latitude: *visit.Hospital.Latitude, Longitude: *visit.Hospital.Longitude,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to calculate check-out distance: %w", err)
	}
	duration := int(input.CheckOutTime.Sub(*visit.CheckedInAt).Seconds())
	result, err := s.compliance.ValidateCheckOut(compliance.CheckOutInput{
		DurationSeconds: duration, CheckOutDistanceM: distance,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate check-out: %w", err)
	}
	checkedOutAt := input.CheckOutTime.UTC()
	visit.CheckedOutAt = &checkedOutAt
	visit.CheckOutLatitude = &input.Latitude
	visit.CheckOutLongitude = &input.Longitude
	visit.CheckOutDistanceM = &distance
	visit.DurationSeconds = &duration
	visit.Status = model.VisitStatusCheckedOut
	visit.AnomalyReasons = mergeAnomalyReasons(visit.AnomalyReasons, result.AnomalyReasons)
	if len(visit.AnomalyReasons) > 0 {
		visit.ComplianceStatus = model.ComplianceStatusNonCompliant
	} else {
		visit.ComplianceStatus = model.ComplianceStatusCompliant
	}
	visit.UpdatedAt = s.clock.Now().UTC()
	if err := s.visitRepo.CheckOut(visit.ID, visit.Version, visit); err != nil {
		if errors.Is(err, repository.ErrVisitStateConflict) {
			return nil, ErrVisitCannotCheckOut
		}
		return nil, err
	}
	return s.visitRepo.FindByID(visit.ID)
}

func mergeAnomalyReasons(existing, additions model.AnomalyReasons) model.AnomalyReasons {
	result := append(model.AnomalyReasons{}, existing...)
	seen := make(map[model.AnomalyReason]struct{}, len(result))
	for _, reason := range result {
		seen[reason] = struct{}{}
	}
	for _, reason := range additions {
		if _, exists := seen[reason]; !exists {
			result = append(result, reason)
			seen[reason] = struct{}{}
		}
	}
	return result
}

// GetVisitByID 根据 ID 获取拜访详情
func (s *VisitService) GetVisitByID(id uuid.UUID) (*model.Visit, error) {
	visit, err := s.visitRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}
	return visit, nil
}
