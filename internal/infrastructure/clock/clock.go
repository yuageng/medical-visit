package clock

import "time"

// Clock 时间抽象接口，便于测试
type Clock interface {
	Now() time.Time
}

// RealClock 生产环境使用的真实时钟
type RealClock struct{}

// Now 返回当前 UTC 时间
func (RealClock) Now() time.Time {
	return time.Now().UTC()
}

// FixedClock 测试用固定时钟
type FixedClock struct {
	Time time.Time
}

// Now 返回固定时间
func (c FixedClock) Now() time.Time {
	return c.Time
}
