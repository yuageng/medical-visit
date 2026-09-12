package clock

import (
	"testing"
	"time"
)

func TestRealClock_Now(t *testing.T) {
	clock := RealClock{}
	now := clock.Now()

	if now.IsZero() {
		t.Error("RealClock.Now() returned zero time")
	}

	// 验证返回的是 UTC 时间
	if now.Location() != time.UTC {
		t.Errorf("RealClock.Now() location = %v, want UTC", now.Location())
	}

	// 验证时间在合理范围内（不是未来时间）
	if now.After(time.Now().UTC().Add(time.Second)) {
		t.Error("RealClock.Now() returned future time")
	}
}

func TestFixedClock_Now(t *testing.T) {
	fixedTime := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	clock := FixedClock{Time: fixedTime}

	now := clock.Now()

	if !now.Equal(fixedTime) {
		t.Errorf("FixedClock.Now() = %v, want %v", now, fixedTime)
	}
}
