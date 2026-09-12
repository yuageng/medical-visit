package model

import "testing"

func TestVisitStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status VisitStatus
		want   bool
	}{
		{"PLANNED is valid", VisitStatusPlanned, true},
		{"CHECKED_IN is valid", VisitStatusCheckedIn, true},
		{"CHECKED_OUT is valid", VisitStatusCheckedOut, true},
		{"COMPLETED is valid", VisitStatusCompleted, true},
		{"CANCELLED is valid", VisitStatusCancelled, true},
		{"invalid status", VisitStatus("INVALID"), false},
		{"empty status", VisitStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("VisitStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComplianceStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status ComplianceStatus
		want   bool
	}{
		{"PENDING is valid", ComplianceStatusPending, true},
		{"COMPLIANT is valid", ComplianceStatusCompliant, true},
		{"NON_COMPLIANT is valid", ComplianceStatusNonCompliant, true},
		{"invalid status", ComplianceStatus("INVALID"), false},
		{"empty status", ComplianceStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("ComplianceStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnomalyReason_Message(t *testing.T) {
	tests := []struct {
		name   string
		reason AnomalyReason
		want   string
	}{
		{
			"INVALID_TIME_SEQUENCE message",
			ReasonInvalidTimeSequence,
			"签退时间早于签到时间",
		},
		{
			"DURATION_TOO_SHORT message",
			ReasonDurationTooShort,
			"停留时长不足 5 分钟",
		},
		{
			"CHECK_IN_TOO_FAR message",
			ReasonCheckInTooFar,
			"签到位置距离医院超过 500 米",
		},
		{
			"CHECK_OUT_TOO_FAR message",
			ReasonCheckOutTooFar,
			"签退位置距离医院超过 500 米",
		},
		{
			"unknown reason",
			AnomalyReason("UNKNOWN"),
			"未知异常",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reason.Message(); got != tt.want {
				t.Errorf("AnomalyReason.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}
