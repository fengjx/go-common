package timex

import "time"

// some time layout or time
const (
	DatetimeLayout = "2006-01-02 15:04:05"
	LayoutWithMs3  = "2006-01-02 15:04:05.000"
	LayoutWithMs6  = "2006-01-02 15:04:05.000000"
	DateLayout     = "2006-01-02"
	TimeLayout     = "15:04:05"

	// ZeroUnix zero unix timestamp
	ZeroUnix int64 = -62135596800

	Second   = 1
	MinSec   = 60
	HourSec  = 3600
	DaySec   = 86400
	WeekSec  = 7 * 86400
	MonthSec = 30 * 86400

	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
)

var (
	// ZeroTime zero time instance
	ZeroTime = time.Time{}
)

// DayStart 对应时间当天0点
func DayStart(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// DayEnd 对应时间当天最后的时间
func DayEnd(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// TodayStart 当天0点
func TodayStart() time.Time {
	return DayStart(time.Now())
}

// TodayEnd 当天最后的时间
func TodayEnd() time.Time {
	return DayEnd(time.Now())
}

// Between 检查时间是否在开始和结束范围内
func Between(dst, start, end time.Time) bool {
	if start.IsZero() && end.IsZero() {
		return false
	}

	if start.IsZero() {
		return dst.Before(end)
	}
	if end.IsZero() {
		return dst.After(start)
	}

	return dst.After(start) && dst.Before(end)
}
