package application

import (
	"errors"
	"testing"
	"time"

	"medical-visit/internal/domain/repository"

	"github.com/google/uuid"
)

func TestMonthlyVisitsByProduct_UsesUTCMonthAndMapsStats(t *testing.T) {
	productID := uuid.New()
	var gotStart, gotEnd time.Time
	repo := &mockVisitRepository{
		monthlyProductStatsFunc: func(start, end time.Time) ([]repository.MonthlyProductStat, error) {
			gotStart, gotEnd = start, end
			return []repository.MonthlyProductStat{{
				ProductID: productID, ProductName: "产品 A", VisitCount: 3, NormalCount: 2, AbnormalCount: 1,
			}}, nil
		},
	}

	result, err := NewDashboardService(repo).MonthlyVisitsByProduct("2026-09")
	if err != nil {
		t.Fatalf("MonthlyVisitsByProduct() error = %v", err)
	}
	if !gotStart.Equal(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("start = %v", gotStart)
	}
	if !gotEnd.Equal(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("end = %v", gotEnd)
	}
	if len(result) != 1 || result[0].ProductID != productID.String() || result[0].VisitCount != 3 || result[0].NormalCount != 2 || result[0].AbnormalCount != 1 {
		t.Fatalf("result = %+v", result)
	}
}

func TestMonthlyVisitsByProduct_InvalidMonth(t *testing.T) {
	called := false
	repo := &mockVisitRepository{monthlyProductStatsFunc: func(start, end time.Time) ([]repository.MonthlyProductStat, error) {
		called = true
		return nil, nil
	}}

	_, err := NewDashboardService(repo).MonthlyVisitsByProduct("2026/09")
	if !errors.Is(err, ErrInvalidDashboardMonth) {
		t.Fatalf("error = %v, want ErrInvalidDashboardMonth", err)
	}
	if called {
		t.Fatal("repository should not be called for invalid month")
	}
}

func TestMonthlyVisitsByProduct_RepositoryError(t *testing.T) {
	wantErr := errors.New("database unavailable")
	repo := &mockVisitRepository{monthlyProductStatsFunc: func(start, end time.Time) ([]repository.MonthlyProductStat, error) {
		return nil, wantErr
	}}

	_, err := NewDashboardService(repo).MonthlyVisitsByProduct("2026-09")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
