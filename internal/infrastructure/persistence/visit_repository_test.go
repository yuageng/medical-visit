package persistence

import (
	"regexp"
	"testing"
	"time"

	"medical-visit/internal/domain/model"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMonthlyProductStats_UsesDatabaseAggregation(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer sqlDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	productID := "00000000-0000-0000-0000-000000000001"

	query := regexp.QuoteMeta(`SELECT v.product_id, p.name AS product_name, COUNT(*) AS visit_count, COUNT(*) FILTER (WHERE v.compliance_status = $1) AS normal_count, COUNT(*) FILTER (WHERE v.compliance_status = $2) AS abnormal_count FROM visits AS v JOIN products AS p ON p.id = v.product_id WHERE v.checked_out_at >= $3 AND v.checked_out_at < $4 AND v.status IN ($5,$6) GROUP BY v.product_id, p.name ORDER BY visit_count DESC, p.name ASC`)
	rows := sqlmock.NewRows([]string{"product_id", "product_name", "visit_count", "normal_count", "abnormal_count"}).
		AddRow(productID, "产品 A", int64(3), int64(2), int64(1))
	mock.ExpectQuery(query).
		WithArgs(model.ComplianceStatusCompliant, model.ComplianceStatusNonCompliant, start, end, model.VisitStatusCheckedOut, model.VisitStatusCompleted).
		WillReturnRows(rows)

	repo := &visitRepositoryImpl{db: db}
	stats, err := repo.MonthlyProductStats(start, end)
	if err != nil {
		t.Fatalf("MonthlyProductStats() error = %v", err)
	}
	if len(stats) != 1 || stats[0].ProductName != "产品 A" || stats[0].VisitCount != 3 || stats[0].NormalCount != 2 || stats[0].AbnormalCount != 1 {
		t.Fatalf("stats = %+v", stats)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL expectations were not met: %v", err)
	}
}
