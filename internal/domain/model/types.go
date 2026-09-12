package model

// VisitStatus 拜访流程状态
type VisitStatus string

const (
	// VisitStatusPlanned 已计划
	VisitStatusPlanned VisitStatus = "PLANNED"
	// VisitStatusCheckedIn 已签到
	VisitStatusCheckedIn VisitStatus = "CHECKED_IN"
	// VisitStatusCheckedOut 已签退
	VisitStatusCheckedOut VisitStatus = "CHECKED_OUT"
	// VisitStatusCompleted 已完成
	VisitStatusCompleted VisitStatus = "COMPLETED"
	// VisitStatusCancelled 已取消
	VisitStatusCancelled VisitStatus = "CANCELLED"
)

// String 返回状态字符串
func (s VisitStatus) String() string {
	return string(s)
}

// IsValid 校验状态是否合法
func (s VisitStatus) IsValid() bool {
	switch s {
	case VisitStatusPlanned, VisitStatusCheckedIn, VisitStatusCheckedOut,
		VisitStatusCompleted, VisitStatusCancelled:
		return true
	}
	return false
}

// ComplianceStatus 合规评估状态
type ComplianceStatus string

const (
	// ComplianceStatusPending 待评估
	ComplianceStatusPending ComplianceStatus = "PENDING"
	// ComplianceStatusCompliant 合规
	ComplianceStatusCompliant ComplianceStatus = "COMPLIANT"
	// ComplianceStatusNonCompliant 不合规
	ComplianceStatusNonCompliant ComplianceStatus = "NON_COMPLIANT"
)

// String 返回合规状态字符串
func (s ComplianceStatus) String() string {
	return string(s)
}

// IsValid 校验合规状态是否合法
func (s ComplianceStatus) IsValid() bool {
	switch s {
	case ComplianceStatusPending, ComplianceStatusCompliant, ComplianceStatusNonCompliant:
		return true
	}
	return false
}

// AnomalyReason 异常原因机器码
type AnomalyReason string

const (
	// ReasonInvalidTimeSequence 签退时间早于签到时间
	ReasonInvalidTimeSequence AnomalyReason = "INVALID_TIME_SEQUENCE"
	// ReasonDurationTooShort 停留时长不足
	ReasonDurationTooShort AnomalyReason = "DURATION_TOO_SHORT"
	// ReasonCheckInTooFar 签到位置距离医院过远
	ReasonCheckInTooFar AnomalyReason = "CHECK_IN_TOO_FAR"
	// ReasonCheckOutTooFar 签退位置距离医院过远
	ReasonCheckOutTooFar AnomalyReason = "CHECK_OUT_TOO_FAR"
)

// String 返回机器码字符串
func (r AnomalyReason) String() string {
	return string(r)
}

// Message 返回人类可读的中文说明
func (r AnomalyReason) Message() string {
	switch r {
	case ReasonInvalidTimeSequence:
		return "签退时间早于签到时间"
	case ReasonDurationTooShort:
		return "停留时长不足 5 分钟"
	case ReasonCheckInTooFar:
		return "签到位置距离医院超过 500 米"
	case ReasonCheckOutTooFar:
		return "签退位置距离医院超过 500 米"
	default:
		return "未知异常"
	}
}
