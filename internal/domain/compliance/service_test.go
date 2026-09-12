package compliance

import (
	"testing"
	"time"

	"medical-visit/internal/domain/model"
)

func TestValidateCheckIn_DistanceBoundary(t *testing.T) {
	service, err := NewService(Rules{MinDurationSeconds: 300, MaxDistanceMeters: 500})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name     string
		distance float64
		status   model.ComplianceStatus
	}{
		{name: "below", distance: 499.999, status: model.ComplianceStatusCompliant},
		{name: "exact", distance: 500, status: model.ComplianceStatusCompliant},
		{name: "above", distance: 500.001, status: model.ComplianceStatusNonCompliant},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ValidateCheckIn(CheckInInput{CheckInTime: time.Now(), CheckInDistanceM: tt.distance})
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != tt.status {
				t.Fatalf("status = %s, want %s", result.Status, tt.status)
			}
		})
	}
}

func TestValidateCheckOut_DurationAndDistanceBoundaries(t *testing.T) {
	service, err := NewService(Rules{MinDurationSeconds: 300, MaxDistanceMeters: 500})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name     string
		duration int
		distance float64
		status   model.ComplianceStatus
		reasons  model.AnomalyReasons
	}{
		{name: "299 seconds", duration: 299, distance: 500, status: model.ComplianceStatusNonCompliant, reasons: model.AnomalyReasons{model.ReasonDurationTooShort}},
		{name: "300 seconds and 500 meters", duration: 300, distance: 500, status: model.ComplianceStatusCompliant},
		{name: "301 seconds", duration: 301, distance: 499.999, status: model.ComplianceStatusCompliant},
		{name: "both anomalies", duration: 299, distance: 500.001, status: model.ComplianceStatusNonCompliant, reasons: model.AnomalyReasons{model.ReasonDurationTooShort, model.ReasonCheckOutTooFar}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ValidateCheckOut(CheckOutInput{DurationSeconds: tt.duration, CheckOutDistanceM: tt.distance})
			if err != nil {
				t.Fatal(err)
			}
			if result.Status != tt.status {
				t.Fatalf("status = %s, want %s", result.Status, tt.status)
			}
			if len(result.AnomalyReasons) != len(tt.reasons) {
				t.Fatalf("reasons = %v, want %v", result.AnomalyReasons, tt.reasons)
			}
			for i := range tt.reasons {
				if result.AnomalyReasons[i] != tt.reasons[i] {
					t.Fatalf("reasons = %v, want %v", result.AnomalyReasons, tt.reasons)
				}
			}
		})
	}
}
