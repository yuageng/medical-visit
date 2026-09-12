package application

import (
	"errors"
	"testing"
	"time"

	"medical-visit/internal/domain/model"
	"medical-visit/internal/domain/repository"
	"medical-visit/internal/infrastructure/clock"

	"github.com/google/uuid"
)

func newCheckOutFixture(status model.VisitStatus, checkedInAt *time.Time, existing model.AnomalyReasons) (*VisitService, *mockVisitRepository, uuid.UUID, time.Time) {
	visitID := uuid.New()
	lat, lon := 31.2304, 121.4737
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	visit := &model.Visit{ID: visitID, Status: status, Version: 2, CheckedInAt: checkedInAt, AnomalyReasons: existing, Hospital: &model.Hospital{Latitude: &lat, Longitude: &lon}}
	repo := &mockVisitRepository{}
	repo.findByIDFunc = func(id uuid.UUID) (*model.Visit, error) {
		if id != visitID {
			return nil, nil
		}
		return visit, nil
	}
	repo.checkOutFunc = func(id uuid.UUID, version int, updated *model.Visit) error {
		visit = updated
		return nil
	}
	return NewVisitService(repo, &mockMasterDataRepository{}, clock.FixedClock{Time: now}), repo, visitID, now
}

func TestCheckOut_SuccessAtFiveMinuteBoundary(t *testing.T) {
	checkedIn := time.Date(2026, 9, 12, 9, 55, 0, 0, time.UTC)
	service, _, visitID, now := newCheckOutFixture(model.VisitStatusCheckedIn, &checkedIn, nil)
	result, err := service.CheckOut(CheckOutInput{VisitID: visitID, Latitude: 31.2304, Longitude: 121.4737, CheckOutTime: now})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != model.VisitStatusCheckedOut || result.DurationSeconds == nil || *result.DurationSeconds != 300 {
		t.Fatalf("result = %+v", result)
	}
	if result.ComplianceStatus != model.ComplianceStatusCompliant || len(result.AnomalyReasons) != 0 {
		t.Fatalf("compliance = %s, reasons = %v", result.ComplianceStatus, result.AnomalyReasons)
	}
}

func TestCheckOut_PreservesCheckInAnomalyAndAddsMultipleAnomalies(t *testing.T) {
	checkedIn := time.Date(2026, 9, 12, 9, 56, 0, 0, time.UTC)
	service, _, visitID, now := newCheckOutFixture(model.VisitStatusCheckedIn, &checkedIn, model.AnomalyReasons{model.ReasonCheckInTooFar})
	result, err := service.CheckOut(CheckOutInput{VisitID: visitID, Latitude: 31.2400, Longitude: 121.4737, CheckOutTime: now})
	if err != nil {
		t.Fatal(err)
	}
	want := []model.AnomalyReason{model.ReasonCheckInTooFar, model.ReasonDurationTooShort, model.ReasonCheckOutTooFar}
	if len(result.AnomalyReasons) != len(want) {
		t.Fatalf("reasons = %v, want %v", result.AnomalyReasons, want)
	}
	for i := range want {
		if result.AnomalyReasons[i] != want[i] {
			t.Fatalf("reasons = %v, want %v", result.AnomalyReasons, want)
		}
	}
}

func TestCheckOut_RejectsTimeBeforeCheckIn(t *testing.T) {
	checkedIn := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	service, _, visitID, _ := newCheckOutFixture(model.VisitStatusCheckedIn, &checkedIn, nil)
	_, err := service.CheckOut(CheckOutInput{VisitID: visitID, Latitude: 31.2304, Longitude: 121.4737, CheckOutTime: checkedIn.Add(-time.Second)})
	if !errors.Is(err, ErrInvalidCheckOutTime) {
		t.Fatalf("error = %v, want ErrInvalidCheckOutTime", err)
	}
}

func TestCheckOut_RejectsIllegalStates(t *testing.T) {
	checkedIn := time.Date(2026, 9, 12, 9, 55, 0, 0, time.UTC)
	for _, status := range []model.VisitStatus{model.VisitStatusPlanned, model.VisitStatusCheckedOut, model.VisitStatusCompleted, model.VisitStatusCancelled} {
		t.Run(status.String(), func(t *testing.T) {
			service, _, visitID, now := newCheckOutFixture(status, &checkedIn, nil)
			_, err := service.CheckOut(CheckOutInput{VisitID: visitID, Latitude: 31.2304, Longitude: 121.4737, CheckOutTime: now})
			if !errors.Is(err, ErrVisitCannotCheckOut) {
				t.Fatalf("error = %v, want ErrVisitCannotCheckOut", err)
			}
		})
	}
}

func TestCheckOut_MapsConcurrentStateConflict(t *testing.T) {
	checkedIn := time.Date(2026, 9, 12, 9, 55, 0, 0, time.UTC)
	service, repo, visitID, now := newCheckOutFixture(model.VisitStatusCheckedIn, &checkedIn, nil)
	repo.checkOutFunc = func(uuid.UUID, int, *model.Visit) error { return repository.ErrVisitStateConflict }
	_, err := service.CheckOut(CheckOutInput{VisitID: visitID, Latitude: 31.2304, Longitude: 121.4737, CheckOutTime: now})
	if !errors.Is(err, ErrVisitCannotCheckOut) {
		t.Fatalf("error = %v, want ErrVisitCannotCheckOut", err)
	}
}
