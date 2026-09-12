package application

import (
	"errors"
	"testing"
	"time"

	"medical-visit/internal/domain/model"
	"medical-visit/internal/infrastructure/clock"

	"github.com/google/uuid"
)

func newCheckInFixture(status model.VisitStatus, hospitalLatitude, hospitalLongitude *float64) (*VisitService, *mockVisitRepository, uuid.UUID, time.Time) {
	visitID := uuid.New()
	hospitalID := uuid.New()
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	visit := &model.Visit{
		ID:               visitID,
		HospitalID:       hospitalID,
		Status:           status,
		Version:          1,
		ComplianceStatus: model.ComplianceStatusPending,
		Hospital:         &model.Hospital{ID: hospitalID, Latitude: hospitalLatitude, Longitude: hospitalLongitude},
	}

	visitRepo := &mockVisitRepository{
		findByIDFunc: func(id uuid.UUID) (*model.Visit, error) {
			if id != visitID {
				return nil, nil
			}
			return visit, nil
		},
		updateFunc: func(updated *model.Visit) error {
			visit = updated
			return nil
		},
	}
	masterRepo := &mockMasterDataRepository{}
	service := NewVisitService(visitRepo, masterRepo, clock.FixedClock{Time: now})
	return service, visitRepo, visitID, now
}

func TestCheckIn_SuccessWithinHospitalRadius(t *testing.T) {
	lat, lon := 31.2304, 121.4737
	service, repo, visitID, now := newCheckInFixture(model.VisitStatusPlanned, &lat, &lon)

	result, err := service.CheckIn(CheckInInput{
		VisitID: visitID, Latitude: lat, Longitude: lon, CheckInTime: now,
	})
	if err != nil {
		t.Fatalf("CheckIn() error = %v", err)
	}
	if result.Status != model.VisitStatusCheckedIn {
		t.Fatalf("status = %s, want CHECKED_IN", result.Status)
	}
	if result.ComplianceStatus != model.ComplianceStatusCompliant {
		t.Fatalf("compliance status = %s, want COMPLIANT", result.ComplianceStatus)
	}
	if result.CheckInDistanceM == nil || *result.CheckInDistanceM != 0 {
		t.Fatalf("check-in distance = %v, want 0", result.CheckInDistanceM)
	}
	if result.CheckedInAt == nil || !result.CheckedInAt.Equal(now) {
		t.Fatalf("checked-in time = %v, want %v", result.CheckedInAt, now)
	}
	if result.CheckInLatitude == nil || result.CheckInLongitude == nil {
		t.Fatal("check-in coordinates were not saved")
	}
	_ = repo
}

func TestCheckIn_Over500MetersMarksAnomaly(t *testing.T) {
	hospitalLat, hospitalLon := 31.2304, 121.4737
	service, _, visitID, now := newCheckInFixture(model.VisitStatusPlanned, &hospitalLat, &hospitalLon)

	result, err := service.CheckIn(CheckInInput{
		VisitID: visitID, Latitude: 31.2400, Longitude: 121.4737, CheckInTime: now,
	})
	if err != nil {
		t.Fatalf("CheckIn() error = %v", err)
	}
	if result.ComplianceStatus != model.ComplianceStatusNonCompliant {
		t.Fatalf("compliance status = %s, want NON_COMPLIANT", result.ComplianceStatus)
	}
	if len(result.AnomalyReasons) != 1 || result.AnomalyReasons[0] != model.ReasonCheckInTooFar {
		t.Fatalf("anomaly reasons = %v, want CHECK_IN_TOO_FAR", result.AnomalyReasons)
	}
	if result.CheckInDistanceM == nil || *result.CheckInDistanceM <= 500 {
		t.Fatalf("distance = %v, want > 500", result.CheckInDistanceM)
	}
}

func TestCheckIn_AlreadyCheckedIn(t *testing.T) {
	lat, lon := 31.2304, 121.4737
	service, _, visitID, now := newCheckInFixture(model.VisitStatusCheckedIn, &lat, &lon)

	_, err := service.CheckIn(CheckInInput{VisitID: visitID, Latitude: lat, Longitude: lon, CheckInTime: now})
	if !errors.Is(err, ErrVisitCannotCheckIn) {
		t.Fatalf("error = %v, want ErrVisitCannotCheckIn", err)
	}
}

func TestCheckIn_InvalidGPS(t *testing.T) {
	lat, lon := 31.2304, 121.4737
	service, _, visitID, now := newCheckInFixture(model.VisitStatusPlanned, &lat, &lon)

	_, err := service.CheckIn(CheckInInput{VisitID: visitID, Latitude: 91, Longitude: lon, CheckInTime: now})
	if !errors.Is(err, ErrInvalidCheckInGPS) {
		t.Fatalf("error = %v, want ErrInvalidCheckInGPS", err)
	}
}

func TestCheckIn_HospitalCoordinatesMissing(t *testing.T) {
	lon := 121.4737
	service, _, visitID, now := newCheckInFixture(model.VisitStatusPlanned, nil, &lon)

	_, err := service.CheckIn(CheckInInput{VisitID: visitID, Latitude: 31.2304, Longitude: lon, CheckInTime: now})
	if !errors.Is(err, ErrHospitalCoordinatesMissing) {
		t.Fatalf("error = %v, want ErrHospitalCoordinatesMissing", err)
	}
}

func TestCheckIn_VisitNotFound(t *testing.T) {
	service, _, _, now := newCheckInFixture(model.VisitStatusPlanned, nil, nil)

	_, err := service.CheckIn(CheckInInput{VisitID: uuid.New(), Latitude: 0, Longitude: 0, CheckInTime: now})
	if !errors.Is(err, ErrVisitNotFound) {
		t.Fatalf("error = %v, want ErrVisitNotFound", err)
	}
}

func TestCheckIn_InvalidStatus(t *testing.T) {
	lat, lon := 31.2304, 121.4737
	service, _, visitID, now := newCheckInFixture(model.VisitStatusCancelled, &lat, &lon)

	_, err := service.CheckIn(CheckInInput{VisitID: visitID, Latitude: lat, Longitude: lon, CheckInTime: now})
	if !errors.Is(err, ErrVisitCannotCheckIn) {
		t.Fatalf("error = %v, want ErrVisitCannotCheckIn", err)
	}
}
