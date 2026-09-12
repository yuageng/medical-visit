package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MedicalRepresentative 医药代表
type MedicalRepresentative struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EmployeeNo string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name       string    `gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time `gorm:"not null;default:now()"`
	UpdatedAt  time.Time `gorm:"not null;default:now()"`
}

// TableName 指定表名
func (MedicalRepresentative) TableName() string {
	return "medical_representatives"
}

// Hospital 医院
type Hospital struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(200);not null"`
	Address   string    `gorm:"type:varchar(500)"`
	Latitude  *float64  `gorm:"type:numeric(9,6)"`
	Longitude *float64  `gorm:"type:numeric(9,6)"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName 指定表名
func (Hospital) TableName() string {
	return "hospitals"
}

// Department 科室
type Department struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	HospitalID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name       string    `gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time `gorm:"not null;default:now()"`
	UpdatedAt  time.Time `gorm:"not null;default:now()"`

	Hospital *Hospital `gorm:"foreignKey:HospitalID;constraint:OnDelete:RESTRICT"`
}

// TableName 指定表名
func (Department) TableName() string {
	return "departments"
}

// HCP 医生（Healthcare Professional）
type HCP struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Title        string    `gorm:"type:varchar(100)"`
	HospitalID   uuid.UUID `gorm:"type:uuid;not null;index"`
	DepartmentID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`

	Hospital   *Hospital   `gorm:"foreignKey:HospitalID;constraint:OnDelete:RESTRICT"`
	Department *Department `gorm:"foreignKey:DepartmentID;constraint:OnDelete:RESTRICT"`
}

// TableName 指定表名
func (HCP) TableName() string {
	return "hcps"
}

// Product 产品
type Product struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code      string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name      string    `gorm:"type:varchar(200);not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// AnomalyReasons 自定义类型用于存储异常原因数组
type AnomalyReasons []AnomalyReason

// Scan 实现 sql.Scanner 接口
func (a *AnomalyReasons) Scan(value interface{}) error {
	if value == nil {
		*a = []AnomalyReason{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	var reasons []string
	if err := json.Unmarshal(bytes, &reasons); err != nil {
		return err
	}

	result := make([]AnomalyReason, len(reasons))
	for i, r := range reasons {
		result[i] = AnomalyReason(r)
	}
	*a = result
	return nil
}

// Value 实现 driver.Valuer 接口
func (a AnomalyReasons) Value() (driver.Value, error) {
	if len(a) == 0 {
		return json.Marshal([]string{})
	}

	reasons := make([]string, len(a))
	for i, r := range a {
		reasons[i] = string(r)
	}
	return json.Marshal(reasons)
}

// Visit 拜访
type Visit struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MRID         uuid.UUID `gorm:"type:uuid;not null;index:idx_visits_mr_id"`
	HCPID        uuid.UUID `gorm:"type:uuid;not null;index"`
	HospitalID   uuid.UUID `gorm:"type:uuid;not null;index"`
	DepartmentID uuid.UUID `gorm:"type:uuid;not null"`
	ProductID    uuid.UUID `gorm:"type:uuid;not null;index:idx_visits_product_checked_out"`

	// 计划信息
	PlannedStartAt time.Time `gorm:"not null;index:idx_visits_mr_id,priority:2"`
	PlanNote       string    `gorm:"type:text"`

	// 状态
	Status VisitStatus `gorm:"type:varchar(30);not null;default:'PLANNED';index"`

	// 签到信息
	CheckedInAt      *time.Time `gorm:"index:idx_visits_product_checked_out,priority:2,sort:desc"`
	CheckInLatitude  *float64   `gorm:"type:numeric(9,6)"`
	CheckInLongitude *float64   `gorm:"type:numeric(9,6)"`
	CheckInDistanceM *float64   `gorm:"type:numeric(10,2)"`

	// 签退信息
	CheckedOutAt      *time.Time `gorm:"index:idx_visits_product_checked_out,priority:2,sort:desc"`
	CheckOutLatitude  *float64   `gorm:"type:numeric(9,6)"`
	CheckOutLongitude *float64   `gorm:"type:numeric(9,6)"`
	CheckOutDistanceM *float64   `gorm:"type:numeric(10,2)"`

	// 时长
	DurationSeconds *int `gorm:"type:integer"`

	// 合规结果
	ComplianceStatus ComplianceStatus `gorm:"type:varchar(20);not null;default:'PENDING';index"`
	AnomalyReasons   AnomalyReasons   `gorm:"type:jsonb;not null;default:'[]'"`

	// 审计与乐观锁
	Version   int       `gorm:"not null;default:1"`
	CreatedAt time.Time `gorm:"not null;default:now();index:,sort:desc"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`

	// 关联
	MR         *MedicalRepresentative `gorm:"foreignKey:MRID;constraint:OnDelete:RESTRICT"`
	HCP        *HCP                   `gorm:"foreignKey:HCPID;constraint:OnDelete:RESTRICT"`
	Hospital   *Hospital              `gorm:"foreignKey:HospitalID;constraint:OnDelete:RESTRICT"`
	Department *Department            `gorm:"foreignKey:DepartmentID;constraint:OnDelete:RESTRICT"`
	Product    *Product               `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT"`
	CallReport *CallReport            `gorm:"foreignKey:VisitID"`
}

// TableName 指定表名
func (Visit) TableName() string {
	return "visits"
}

// AcademicMaterial 系统中的合规学术资料。
type AcademicMaterial struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code      string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Name      string     `gorm:"type:varchar(255);not null"`
	ProductID *uuid.UUID `gorm:"type:uuid;index"`
	Active    bool       `gorm:"not null;default:true;index"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
}

// TableName 指定表名。
func (AcademicMaterial) TableName() string { return "academic_materials" }

// CallReport 拜访记录。
type CallReport struct {
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	VisitID              uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	DiscussionSummary    string    `gorm:"type:text;not null"`
	HCPFeedback          string    `gorm:"type:text"`
	MaterialsDistributed bool      `gorm:"not null;default:false"`
	MaterialNote         string    `gorm:"type:text"`
	AdditionalNotes      string    `gorm:"type:text"`
	CreatedAt            time.Time `gorm:"not null;default:now()"`
	UpdatedAt            time.Time `gorm:"not null;default:now()"`

	Visit     *Visit              `gorm:"foreignKey:VisitID;constraint:OnDelete:RESTRICT"`
	Materials []*AcademicMaterial `gorm:"many2many:call_report_materials;"`
}

// TableName 指定表名
func (CallReport) TableName() string {
	return "call_reports"
}
