package persistence

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"medical-visit/internal/domain/model"
	"medical-visit/internal/domain/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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

	query := regexp.QuoteMeta(`SELECT v.product_id, p.name AS product_name, COUNT(*) AS visit_count, COUNT(*) FILTER (WHERE v.compliance_status = $1) AS normal_count, COUNT(*) FILTER (WHERE v.compliance_status = $2) AS abnormal_count FROM visits AS v JOIN products AS p ON p.id = v.product_id WHERE v.checked_out_at >= $3 AND v.checked_out_at < $4 AND v.status IN ($5,$6) AND v.compliance_status IN ($7,$8) GROUP BY v.product_id, p.name ORDER BY visit_count DESC, p.name ASC`)
	rows := sqlmock.NewRows([]string{"product_id", "product_name", "visit_count", "normal_count", "abnormal_count"}).
		AddRow(productID, "产品 A", int64(3), int64(2), int64(1))
	mock.ExpectQuery(query).
		WithArgs(model.ComplianceStatusCompliant, model.ComplianceStatusNonCompliant, start, end, model.VisitStatusCheckedOut, model.VisitStatusCompleted, model.ComplianceStatusCompliant, model.ComplianceStatusNonCompliant).
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

func TestCheckOut_UsesStatusAndVersionCondition(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	visitID := "00000000-0000-0000-0000-000000000001"
	checkedOutAt := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	lat, lon, distance, duration := 31.2304, 121.4737, 0.0, 300
	visit := &model.Visit{
		Status: model.VisitStatusCheckedOut, CheckedOutAt: &checkedOutAt,
		CheckOutLatitude: &lat, CheckOutLongitude: &lon, CheckOutDistanceM: &distance,
		DurationSeconds: &duration, ComplianceStatus: model.ComplianceStatusCompliant,
		AnomalyReasons: model.AnomalyReasons{}, UpdatedAt: checkedOutAt,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "visits" SET .* WHERE id = \$[0-9]+ AND status = \$[0-9]+ AND version = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if err := (&visitRepositoryImpl{db: db}).CheckOut(uuid.MustParse(visitID), 2, visit); err != nil {
		t.Fatalf("CheckOut() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckOut_ReturnsConflictWhenConditionDoesNotMatch(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "visits" SET .* WHERE id = \$[0-9]+ AND status = \$[0-9]+ AND version = \$[0-9]+`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	err = (&visitRepositoryImpl{db: db}).CheckOut(uuid.New(), 2, &model.Visit{})
	if !errors.Is(err, repository.ErrVisitStateConflict) {
		t.Fatalf("error = %v, want ErrVisitStateConflict", err)
	}
}
