package persistence

import (
	"errors"
	"fmt"

	"medical-visit/internal/domain/model"
	"medical-visit/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// masterDataRepositoryImpl 主数据仓储实现
type masterDataRepositoryImpl struct {
	db *gorm.DB
}

// NewMasterDataRepository 创建主数据仓储实例
func NewMasterDataRepository(db *gorm.DB) repository.MasterDataRepository {
	return &masterDataRepositoryImpl{db: db}
}

// FindMRByID 根据 ID 查询医药代表
func (r *masterDataRepositoryImpl) FindMRByID(id uuid.UUID) (*model.MedicalRepresentative, error) {
	var mr model.MedicalRepresentative
	err := r.db.First(&mr, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find MR: %w", err)
	}
	return &mr, nil
}

// FindAllMRs 查询所有医药代表
func (r *masterDataRepositoryImpl) FindAllMRs() ([]*model.MedicalRepresentative, error) {
	var mrs []*model.MedicalRepresentative
	err := r.db.Order("name ASC").Find(&mrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all MRs: %w", err)
	}
	return mrs, nil
}

// FindHCPByID 根据 ID 查询医生
func (r *masterDataRepositoryImpl) FindHCPByID(id uuid.UUID) (*model.HCP, error) {
	var hcp model.HCP
	err := r.db.Preload("Hospital").Preload("Department").First(&hcp, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find HCP: %w", err)
	}
	return &hcp, nil
}

// FindHCPsByHospital 查询某医院的所有医生
func (r *masterDataRepositoryImpl) FindHCPsByHospital(hospitalID uuid.UUID) ([]*model.HCP, error) {
	var hcps []*model.HCP
	err := r.db.
		Preload("Department").
		Where("hospital_id = ?", hospitalID).
		Order("name ASC").
		Find(&hcps).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find HCPs by hospital: %w", err)
	}
	return hcps, nil
}

// FindHospitalByID 根据 ID 查询医院
func (r *masterDataRepositoryImpl) FindHospitalByID(id uuid.UUID) (*model.Hospital, error) {
	var hospital model.Hospital
	err := r.db.First(&hospital, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find hospital: %w", err)
	}
	return &hospital, nil
}

// FindAllHospitals 查询所有医院
func (r *masterDataRepositoryImpl) FindAllHospitals() ([]*model.Hospital, error) {
	var hospitals []*model.Hospital
	err := r.db.Order("name ASC").Find(&hospitals).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all hospitals: %w", err)
	}
	return hospitals, nil
}

// FindDepartmentByID 根据 ID 查询科室
func (r *masterDataRepositoryImpl) FindDepartmentByID(id uuid.UUID) (*model.Department, error) {
	var department model.Department
	err := r.db.Preload("Hospital").First(&department, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find department: %w", err)
	}
	return &department, nil
}

// FindDepartmentsByHospital 查询某医院的所有科室
func (r *masterDataRepositoryImpl) FindDepartmentsByHospital(hospitalID uuid.UUID) ([]*model.Department, error) {
	var departments []*model.Department
	err := r.db.
		Where("hospital_id = ?", hospitalID).
		Order("name ASC").
		Find(&departments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find departments by hospital: %w", err)
	}
	return departments, nil
}

// FindProductByID 根据 ID 查询产品
func (r *masterDataRepositoryImpl) FindProductByID(id uuid.UUID) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	return &product, nil
}

// FindAcademicMaterialsByIDs 根据 ID 批量查询学术资料，由 Service 校验 active 状态。
func (r *masterDataRepositoryImpl) FindAcademicMaterialsByIDs(ids []uuid.UUID) ([]*model.AcademicMaterial, error) {
	if len(ids) == 0 {
		return []*model.AcademicMaterial{}, nil
	}
	var materials []*model.AcademicMaterial
	if err := r.db.Where("id IN ?", ids).Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("failed to find academic materials: %w", err)
	}
	return materials, nil
}

// FindAllProducts 查询所有产品
func (r *masterDataRepositoryImpl) FindAllProducts() ([]*model.Product, error) {
	var products []*model.Product
	err := r.db.Order("name ASC").Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all products: %w", err)
	}
	return products, nil
}
