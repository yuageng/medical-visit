package compliance

import (
	"errors"
	"fmt"
	"math"
	"time"

	"medical-visit/internal/domain/model"
)

// Rules 合规规则配置。
type Rules struct {
	MinDurationSeconds int
	MaxDistanceMeters  float64
}

// DefaultRules 返回拜访合规校验的默认规则。
func DefaultRules() Rules {
	return Rules{MinDurationSeconds: 300, MaxDistanceMeters: 500}
}

// Validate 校验规则配置。
func (r Rules) Validate() error {
	if r.MinDurationSeconds < 0 {
		return errors.New("min duration must be non-negative")
	}
	if math.IsNaN(r.MaxDistanceMeters) || math.IsInf(r.MaxDistanceMeters, 0) || r.MaxDistanceMeters < 0 {
		return errors.New("max distance must be a finite non-negative number")
	}
	return nil
}

// CheckInInput 签到阶段的合规输入。
type CheckInInput struct {
	CheckInTime      time.Time
	CheckInDistanceM float64
}

// CheckInResult 签到阶段的合规结果。
type CheckInResult struct {
	Status         model.ComplianceStatus
	AnomalyReasons model.AnomalyReasons
}

// Service 独立的合规领域服务。
type Service struct {
	rules Rules
}

// NewService 创建合规服务。
func NewService(rules Rules) (*Service, error) {
	if err := rules.Validate(); err != nil {
		return nil, fmt.Errorf("invalid compliance rules: %w", err)
	}
	return &Service{rules: rules}, nil
}

// ValidateCheckIn 评估签到位置，不依赖数据库、HTTP 或系统时钟。
func (s *Service) ValidateCheckIn(input CheckInInput) (CheckInResult, error) {
	if input.CheckInTime.IsZero() {
		return CheckInResult{}, errors.New("check-in time is required")
	}
	if math.IsNaN(input.CheckInDistanceM) || math.IsInf(input.CheckInDistanceM, 0) || input.CheckInDistanceM < 0 {
		return CheckInResult{}, errors.New("check-in distance must be a finite non-negative number")
	}

	result := CheckInResult{
		Status:         model.ComplianceStatusCompliant,
		AnomalyReasons: model.AnomalyReasons{},
	}
	if input.CheckInDistanceM > s.rules.MaxDistanceMeters {
		result.Status = model.ComplianceStatusNonCompliant
		result.AnomalyReasons = append(result.AnomalyReasons, model.ReasonCheckInTooFar)
	}
	return result, nil
}

// CheckOutInput 签退阶段的合规输入。
type CheckOutInput struct {
	DurationSeconds   int
	CheckOutDistanceM float64
}

// ValidateCheckOut 评估签退时长与位置；恰好达到阈值视为合规。
func (s *Service) ValidateCheckOut(input CheckOutInput) (CheckInResult, error) {
	if input.DurationSeconds < 0 {
		return CheckInResult{}, errors.New("duration must be non-negative")
	}
	if math.IsNaN(input.CheckOutDistanceM) || math.IsInf(input.CheckOutDistanceM, 0) || input.CheckOutDistanceM < 0 {
		return CheckInResult{}, errors.New("check-out distance must be a finite non-negative number")
	}
	result := CheckInResult{Status: model.ComplianceStatusCompliant, AnomalyReasons: model.AnomalyReasons{}}
	if input.DurationSeconds < s.rules.MinDurationSeconds {
		result.Status = model.ComplianceStatusNonCompliant
		result.AnomalyReasons = append(result.AnomalyReasons, model.ReasonDurationTooShort)
	}
	if input.CheckOutDistanceM > s.rules.MaxDistanceMeters {
		result.Status = model.ComplianceStatusNonCompliant
		result.AnomalyReasons = append(result.AnomalyReasons, model.ReasonCheckOutTooFar)
	}
	return result, nil
}
