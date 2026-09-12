package application

import (
	"testing"
	"time"

	"medical-visit/internal/domain/model"
	"medical-visit/internal/domain/repository"
	"medical-visit/internal/infrastructure/clock"

	"github.com/google/uuid"
)

// mockVisitRepository 模拟拜访仓储
type mockVisitRepository struct {
	createFunc               func(*model.Visit) error
	checkInFunc              func(uuid.UUID, int, *model.Visit) error
	checkOutFunc             func(uuid.UUID, int, *model.Visit) error
	saveCallReportFunc       func(uuid.UUID, int, *model.CallReport, []uuid.UUID) error
	findByIDFunc             func(uuid.UUID) (*model.Visit, error)
	updateFunc               func(*model.Visit) error
	findByMRAndTimeRangeFunc func(uuid.UUID, string, string) ([]*model.Visit, error)
	monthlyProductStatsFunc  func(time.Time, time.Time) ([]repository.MonthlyProductStat, error)
}

func (m *mockVisitRepository) Create(visit *model.Visit) error {
	if m.createFunc != nil {
		return m.createFunc(visit)
	}
	return nil
}

func (m *mockVisitRepository) CheckIn(visitID uuid.UUID, expectedVersion int, visit *model.Visit) error {
	if m.checkInFunc != nil {
		return m.checkInFunc(visitID, expectedVersion, visit)
	}
	if m.updateFunc != nil {
		return m.updateFunc(visit)
	}
	return nil
}

func (m *mockVisitRepository) CheckOut(visitID uuid.UUID, expectedVersion int, visit *model.Visit) error {
	if m.checkOutFunc != nil {
		return m.checkOutFunc(visitID, expectedVersion, visit)
	}
	if m.updateFunc != nil {
		return m.updateFunc(visit)
	}
	return nil
}

func (m *mockVisitRepository) SaveCallReport(visitID uuid.UUID, expectedVersion int, report *model.CallReport, materialIDs []uuid.UUID) error {
	if m.saveCallReportFunc != nil {
		return m.saveCallReportFunc(visitID, expectedVersion, report, materialIDs)
	}
	return nil
}

func (m *mockVisitRepository) FindByID(id uuid.UUID) (*model.Visit, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return nil, nil
}

func (m *mockVisitRepository) Update(visit *model.Visit) error {
	if m.updateFunc != nil {
		return m.updateFunc(visit)
	}
	return nil
}

func (m *mockVisitRepository) MonthlyProductStats(start, end time.Time) ([]repository.MonthlyProductStat, error) {
	if m.monthlyProductStatsFunc != nil {
		return m.monthlyProductStatsFunc(start, end)
	}
	return nil, nil
}

func (m *mockVisitRepository) FindByMRAndTimeRange(mrID uuid.UUID, startTime, endTime string) ([]*model.Visit, error) {
	if m.findByMRAndTimeRangeFunc != nil {
		return m.findByMRAndTimeRangeFunc(mrID, startTime, endTime)
	}
	return nil, nil
}

// mockMasterDataRepository 模拟主数据仓储
type mockMasterDataRepository struct {
	findMRByIDFunc                 func(uuid.UUID) (*model.MedicalRepresentative, error)
	findHCPByIDFunc                func(uuid.UUID) (*model.HCP, error)
	findHospitalByIDFunc           func(uuid.UUID) (*model.Hospital, error)
	findDepartmentByIDFunc         func(uuid.UUID) (*model.Department, error)
	findProductByIDFunc            func(uuid.UUID) (*model.Product, error)
	findAllMRsFunc                 func() ([]*model.MedicalRepresentative, error)
	findHCPsByHospitalFunc         func(uuid.UUID) ([]*model.HCP, error)
	findAllHospitalsFunc           func() ([]*model.Hospital, error)
	findDepartmentsByHospitalFunc  func(uuid.UUID) ([]*model.Department, error)
	findAllProductsFunc            func() ([]*model.Product, error)
	findAcademicMaterialsByIDsFunc func([]uuid.UUID) ([]*model.AcademicMaterial, error)
}

func (m *mockMasterDataRepository) FindMRByID(id uuid.UUID) (*model.MedicalRepresentative, error) {
	if m.findMRByIDFunc != nil {
		return m.findMRByIDFunc(id)
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindAllMRs() ([]*model.MedicalRepresentative, error) {
	if m.findAllMRsFunc != nil {
		return m.findAllMRsFunc()
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindHCPByID(id uuid.UUID) (*model.HCP, error) {
	if m.findHCPByIDFunc != nil {
		return m.findHCPByIDFunc(id)
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindHCPsByHospital(hospitalID uuid.UUID) ([]*model.HCP, error) {
	if m.findHCPsByHospitalFunc != nil {
		return m.findHCPsByHospitalFunc(hospitalID)
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindHospitalByID(id uuid.UUID) (*model.Hospital, error) {
	if m.findHospitalByIDFunc != nil {
		return m.findHospitalByIDFunc(id)
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindAllHospitals() ([]*model.Hospital, error) {
	if m.findAllHospitalsFunc != nil {
		return m.findAllHospitalsFunc()
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindDepartmentByID(id uuid.UUID) (*model.Department, error) {
	if m.findDepartmentByIDFunc != nil {
		return m.findDepartmentByIDFunc(id)
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindDepartmentsByHospital(hospitalID uuid.UUID) ([]*model.Department, error) {
	if m.findDepartmentsByHospitalFunc != nil {
		return m.findDepartmentsByHospitalFunc(hospitalID)
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindProductByID(id uuid.UUID) (*model.Product, error) {
	if m.findProductByIDFunc != nil {
		return m.findProductByIDFunc(id)
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindAcademicMaterialsByIDs(ids []uuid.UUID) ([]*model.AcademicMaterial, error) {
	if m.findAcademicMaterialsByIDsFunc != nil {
		return m.findAcademicMaterialsByIDsFunc(ids)
	}
	return nil, nil
}

func (m *mockMasterDataRepository) FindAllProducts() ([]*model.Product, error) {
	if m.findAllProductsFunc != nil {
		return m.findAllProductsFunc()
	}
	return nil, nil
}

func TestCreateVisit_Success(t *testing.T) {
	// 准备测试数据
	mrID := uuid.New()
	hcpID := uuid.New()
	hospitalID := uuid.New()
	departmentID := uuid.New()
	productID := uuid.New()
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	plannedTime := now.Add(1 * time.Hour)

	// 创建 mock 仓储
	visitRepo := &mockVisitRepository{
		createFunc: func(visit *model.Visit) error {
			return nil
		},
		findByIDFunc: func(id uuid.UUID) (*model.Visit, error) {
			return &model.Visit{
				ID:               id,
				MRID:             mrID,
				HCPID:            hcpID,
				HospitalID:       hospitalID,
				DepartmentID:     departmentID,
				ProductID:        productID,
				PlannedStartAt:   plannedTime,
				Status:           model.VisitStatusPlanned,
				ComplianceStatus: model.ComplianceStatusPending,
				CreatedAt:        now,
				UpdatedAt:        now,
			}, nil
		},
	}

	masterDataRepo := &mockMasterDataRepository{
		findMRByIDFunc: func(id uuid.UUID) (*model.MedicalRepresentative, error) {
			return &model.MedicalRepresentative{ID: mrID, Name: "张三"}, nil
		},
		findHCPByIDFunc: func(id uuid.UUID) (*model.HCP, error) {
			return &model.HCP{
				ID:           hcpID,
				HospitalID:   hospitalID,
				DepartmentID: departmentID,
				Name:         "李医生",
			}, nil
		},
		findHospitalByIDFunc: func(id uuid.UUID) (*model.Hospital, error) {
			return &model.Hospital{ID: hospitalID, Name: "市人民医院"}, nil
		},
		findDepartmentByIDFunc: func(id uuid.UUID) (*model.Department, error) {
			return &model.Department{
				ID:         departmentID,
				HospitalID: hospitalID,
				Name:       "心内科",
			}, nil
		},
		findProductByIDFunc: func(id uuid.UUID) (*model.Product, error) {
			return &model.Product{ID: productID, Name: "产品 A"}, nil
		},
	}

	clk := clock.FixedClock{Time: now}

	// 创建服务
	service := NewVisitService(visitRepo, masterDataRepo, clk)

	// 执行测试
	input := CreateVisitInput{
		MRID:           mrID,
		HCPID:          hcpID,
		HospitalID:     hospitalID,
		DepartmentID:   departmentID,
		ProductID:      productID,
		PlannedStartAt: plannedTime,
		PlanNote:       "测试拜访",
	}

	visit, err := service.CreateVisit(input)

	// 验证结果
	if err != nil {
		t.Fatalf("CreateVisit failed: %v", err)
	}

	if visit.Status != model.VisitStatusPlanned {
		t.Errorf("Expected status PLANNED, got %s", visit.Status)
	}

	if visit.ComplianceStatus != model.ComplianceStatusPending {
		t.Errorf("Expected compliance status PENDING, got %s", visit.ComplianceStatus)
	}
}

func TestCreateVisit_MRNotFound(t *testing.T) {
	visitRepo := &mockVisitRepository{}
	masterDataRepo := &mockMasterDataRepository{
		findMRByIDFunc: func(id uuid.UUID) (*model.MedicalRepresentative, error) {
			return nil, nil // MR 不存在
		},
	}
	clk := clock.FixedClock{Time: time.Now()}
	service := NewVisitService(visitRepo, masterDataRepo, clk)

	input := CreateVisitInput{
		MRID:           uuid.New(),
		HCPID:          uuid.New(),
		HospitalID:     uuid.New(),
		DepartmentID:   uuid.New(),
		ProductID:      uuid.New(),
		PlannedStartAt: time.Now().Add(1 * time.Hour),
	}

	_, err := service.CreateVisit(input)

	if err != ErrMRNotFound {
		t.Errorf("Expected ErrMRNotFound, got %v", err)
	}
}

func TestCreateVisit_DepartmentNotBelongToHospital(t *testing.T) {
	mrID := uuid.New()
	hcpID := uuid.New()
	hospitalID := uuid.New()
	departmentID := uuid.New()
	productID := uuid.New()
	otherHospitalID := uuid.New() // 不同的医院

	visitRepo := &mockVisitRepository{}
	masterDataRepo := &mockMasterDataRepository{
		findMRByIDFunc: func(id uuid.UUID) (*model.MedicalRepresentative, error) {
			return &model.MedicalRepresentative{ID: mrID}, nil
		},
		findHCPByIDFunc: func(id uuid.UUID) (*model.HCP, error) {
			return &model.HCP{
				ID:           hcpID,
				HospitalID:   hospitalID,
				DepartmentID: departmentID,
			}, nil
		},
		findHospitalByIDFunc: func(id uuid.UUID) (*model.Hospital, error) {
			return &model.Hospital{ID: hospitalID}, nil
		},
		findDepartmentByIDFunc: func(id uuid.UUID) (*model.Department, error) {
			// 科室属于另一个医院
			return &model.Department{
				ID:         departmentID,
				HospitalID: otherHospitalID,
			}, nil
		},
		findProductByIDFunc: func(id uuid.UUID) (*model.Product, error) {
			return &model.Product{ID: productID}, nil
		},
	}

	clk := clock.FixedClock{Time: time.Now()}
	service := NewVisitService(visitRepo, masterDataRepo, clk)

	input := CreateVisitInput{
		MRID:           mrID,
		HCPID:          hcpID,
		HospitalID:     hospitalID,
		DepartmentID:   departmentID,
		ProductID:      productID,
		PlannedStartAt: time.Now().Add(1 * time.Hour),
	}

	_, err := service.CreateVisit(input)

	if err != ErrDepartmentNotBelongToHospital {
		t.Errorf("Expected ErrDepartmentNotBelongToHospital, got %v", err)
	}
}

func TestCreateVisit_HCPNotInHospital(t *testing.T) {
	mrID := uuid.New()
	hcpID := uuid.New()
	hospitalID := uuid.New()
	departmentID := uuid.New()
	productID := uuid.New()
	otherHospitalID := uuid.New()

	visitRepo := &mockVisitRepository{}
	masterDataRepo := &mockMasterDataRepository{
		findMRByIDFunc: func(id uuid.UUID) (*model.MedicalRepresentative, error) {
			return &model.MedicalRepresentative{ID: mrID}, nil
		},
		findHCPByIDFunc: func(id uuid.UUID) (*model.HCP, error) {
			// HCP 属于另一个医院
			return &model.HCP{
				ID:           hcpID,
				HospitalID:   otherHospitalID,
				DepartmentID: departmentID,
			}, nil
		},
		findHospitalByIDFunc: func(id uuid.UUID) (*model.Hospital, error) {
			return &model.Hospital{ID: hospitalID}, nil
		},
		findDepartmentByIDFunc: func(id uuid.UUID) (*model.Department, error) {
			return &model.Department{
				ID:         departmentID,
				HospitalID: hospitalID,
			}, nil
		},
		findProductByIDFunc: func(id uuid.UUID) (*model.Product, error) {
			return &model.Product{ID: productID}, nil
		},
	}

	clk := clock.FixedClock{Time: time.Now()}
	service := NewVisitService(visitRepo, masterDataRepo, clk)

	input := CreateVisitInput{
		MRID:           mrID,
		HCPID:          hcpID,
		HospitalID:     hospitalID,
		DepartmentID:   departmentID,
		ProductID:      productID,
		PlannedStartAt: time.Now().Add(1 * time.Hour),
	}

	_, err := service.CreateVisit(input)

	if err != ErrHCPNotInHospital {
		t.Errorf("Expected ErrHCPNotInHospital, got %v", err)
	}
}

func TestCreateVisit_PlannedTimeInPast(t *testing.T) {
	mrID := uuid.New()
	hcpID := uuid.New()
	hospitalID := uuid.New()
	departmentID := uuid.New()
	productID := uuid.New()
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	pastTime := now.Add(-10 * time.Minute) // 10 分钟前

	visitRepo := &mockVisitRepository{}
	masterDataRepo := &mockMasterDataRepository{
		findMRByIDFunc: func(id uuid.UUID) (*model.MedicalRepresentative, error) {
			return &model.MedicalRepresentative{ID: mrID}, nil
		},
		findHCPByIDFunc: func(id uuid.UUID) (*model.HCP, error) {
			return &model.HCP{
				ID:           hcpID,
				HospitalID:   hospitalID,
				DepartmentID: departmentID,
			}, nil
		},
		findHospitalByIDFunc: func(id uuid.UUID) (*model.Hospital, error) {
			return &model.Hospital{ID: hospitalID}, nil
		},
		findDepartmentByIDFunc: func(id uuid.UUID) (*model.Department, error) {
			return &model.Department{
				ID:         departmentID,
				HospitalID: hospitalID,
			}, nil
		},
		findProductByIDFunc: func(id uuid.UUID) (*model.Product, error) {
			return &model.Product{ID: productID}, nil
		},
	}

	clk := clock.FixedClock{Time: now}
	service := NewVisitService(visitRepo, masterDataRepo, clk)

	input := CreateVisitInput{
		MRID:           mrID,
		HCPID:          hcpID,
		HospitalID:     hospitalID,
		DepartmentID:   departmentID,
		ProductID:      productID,
		PlannedStartAt: pastTime,
	}

	_, err := service.CreateVisit(input)

	if err != ErrPlannedTimeInPast {
		t.Errorf("Expected ErrPlannedTimeInPast, got %v", err)
	}
}
