package dateutil

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// DefaultDateTimeFormat 默认日期时间格式（Go layout），对应 Java yyyy-MM-dd HH:mm:ss。
	DefaultDateTimeFormat = "2006-01-02 15:04:05"
	// DefaultDateFormat 默认日期格式（Go layout），对应 Java yyyy-MM-dd。
	DefaultDateFormat = "2006-01-02"

	// DaySeconds 一天的秒数。
	DaySeconds = 24 * 3600
	// DayMS 一天的毫秒数。
	DayMS = DaySeconds * 1000
	// WeekSeconds 一周的秒数。
	WeekSeconds = 7 * DaySeconds
	// WeekMS 一周的毫秒数。
	WeekMS = WeekSeconds * 1000
)

var (
	defaultLoc atomic.Pointer[time.Location]
	timeOffset atomic.Int64 // 逻辑时间向前偏移（毫秒）

	listenerMu  sync.Mutex
	listeners   []timeChangeListener
	listenerSeq atomic.Uint64

	formatterMu sync.RWMutex
	formatters  = map[string]string{
		DefaultDateTimeFormat: DefaultDateTimeFormat,
		DefaultDateFormat:     DefaultDateFormat,
		"yyyy-MM-dd HH:mm:ss": DefaultDateTimeFormat,
		"yyyy-MM-dd":          DefaultDateFormat,
	}
)

func init() {
	defaultLoc.Store(time.Local)
}

// SetDefaultZoneID 设置默认时区。对应 Java DateUtil.setDefaultZoneId。
func SetDefaultZoneID(loc *time.Location) {
	if loc == nil {
		loc = time.Local
	}
	defaultLoc.Store(loc)
}

// GetDefaultZoneID 获取默认时区。对应 Java DateUtil.getDefaultZoneId。
func GetDefaultZoneID() *time.Location {
	return defaultLoc.Load()
}

// SetTimeOffset 对全局逻辑时间向前偏移（必须 >0）。对应 Java DateUtil.setTimeOffset。
func SetTimeOffset(offset time.Duration) error {
	if offset <= 0 {
		return fmt.Errorf("utils: 逻辑时间只允许向前调整, offset必须大于0, 传入: %v", offset)
	}
	before := CurrentTimeMillis()
	timeOffset.Add(offset.Milliseconds())
	after := CurrentTimeMillis()
	if after <= before {
		return fmt.Errorf("utils: 逻辑时间调整后未前进, before: %d, after: %d", before, after)
	}
	notifyTimeChanged()
	return nil
}

// ClearTimeOffset 清除时间偏移，恢复系统时钟。对应 Java DateUtil.clearTimeOffset。
func ClearTimeOffset() {
	timeOffset.Store(0)
	notifyTimeChanged()
}

// timeChangeListener 带 ID 的时间变化监听项。
type timeChangeListener struct {
	id uint64
	fn func()
}

// AddTimeChangeListener 注册逻辑时间变化监听，返回取消函数（可安全多次调用）。
// 对应 Java DateUtil.addTimeChangeListener；移除请调用返回的 cancel，不再按函数指针比对。
func AddTimeChangeListener(listener func()) (cancel func()) {
	if listener == nil {
		return func() {}
	}
	id := listenerSeq.Add(1)
	listenerMu.Lock()
	listeners = append(listeners, timeChangeListener{id: id, fn: listener})
	listenerMu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			listenerMu.Lock()
			defer listenerMu.Unlock()
			dst := listeners[:0]
			for _, item := range listeners {
				if item.id != id {
					dst = append(dst, item)
				}
			}
			listeners = dst
		})
	}
}

// notifyTimeChanged 通知所有时间变化监听。
func notifyTimeChanged() {
	listenerMu.Lock()
	cp := append([]timeChangeListener{}, listeners...)
	listenerMu.Unlock()
	for _, item := range cp {
		func() {
			defer func() { _ = recover() }()
			item.fn()
		}()
	}
}

// CurrentTimeMillis 当前逻辑时间毫秒戳。对应 Java DateUtil.currentTimeMillis。
func CurrentTimeMillis() int64 {
	return time.Now().UnixMilli() + timeOffset.Load()
}

// CurrSeconds 当前逻辑时间秒戳。对应 Java DateUtil.currSeconds。
func CurrSeconds() int64 {
	return CurrentTimeMillis() / 1000
}

// Now 当前逻辑时间（默认时区）。对应 Java DateUtil.nowLocalDateTime()。
func Now() time.Time {
	return NowIn(GetDefaultZoneID())
}

// NowIn 当前逻辑时间（指定时区）。对应 Java DateUtil.nowLocalDateTime(ZoneId)。
func NowIn(loc *time.Location) time.Time {
	if loc == nil {
		loc = GetDefaultZoneID()
	}
	return time.UnixMilli(CurrentTimeMillis()).In(loc)
}

// CalConsumeMillisSeconds 计算 fn 执行耗时毫秒（用系统时钟，不受逻辑偏移影响）。
// 对应 Java DateUtil.calConsumeMillisSeconds。
func CalConsumeMillisSeconds(fn func()) int64 {
	start := time.Now()
	fn()
	return time.Since(start).Milliseconds()
}

// DateToStringMillis 毫秒戳转默认格式字符串。对应 Java DateUtil.dateToString(long)。
func DateToStringMillis(millis int64) string {
	return DateToStringMillisFormat(millis, DefaultDateTimeFormat)
}

// DateToStringMillisFormat 毫秒戳按格式转字符串。对应 Java DateUtil.dateToString(long, format)。
func DateToStringMillisFormat(millis int64, format string) string {
	return DateToStringFormat(MilliToTime(millis), format)
}

// DateToString 时间转默认格式字符串。对应 Java DateUtil.dateToString(LocalDateTime)。
func DateToString(t time.Time) string {
	return DateToStringFormat(t, DefaultDateTimeFormat)
}

// DateToStringFormat 时间按格式转字符串。format 支持 Go layout，也兼容常见 Java pattern。
// 对应 Java DateUtil.dateToString(LocalDateTime, format)。
func DateToStringFormat(t time.Time, format string) string {
	return t.Format(resolveFormat(format))
}

// StringToDate 按默认格式解析时间。对应 Java DateUtil.stringToDate(String)。
func StringToDate(s string) (time.Time, error) {
	return StringToDateFormat(s, DefaultDateTimeFormat)
}

// StringToDateFormat 按格式解析时间（默认时区）。对应 Java DateUtil.stringToDate(String, format)。
func StringToDateFormat(s, format string) (time.Time, error) {
	return time.ParseInLocation(resolveFormat(format), s, GetDefaultZoneID())
}

// GetMilliByTime 时间转毫秒（默认时区解释）。对应 Java DateUtil.getMilliByTime(LocalDateTime)。
func GetMilliByTime(t time.Time) int64 {
	return GetMilliByTimeIn(t, GetDefaultZoneID())
}

// GetMilliByTimeIn 时间转毫秒（指定时区）。对应 Java DateUtil.getMilliByTime(LocalDateTime, ZoneId)。
func GetMilliByTimeIn(t time.Time, loc *time.Location) int64 {
	if loc == nil {
		loc = GetDefaultZoneID()
	}
	// 取墙钟字段，按 loc 重新锚定，对齐 Java LocalDateTime.atZone(zoneId)
	y, m, d := t.Date()
	hh, mm, ss := t.Clock()
	ns := t.Nanosecond()
	return time.Date(y, m, d, hh, mm, ss, ns, loc).UnixMilli()
}

// GetSecondsByTime 时间转秒。对应 Java DateUtil.getSecondsByTime。
func GetSecondsByTime(t time.Time) int64 {
	return GetMilliByTime(t) / 1000
}

// MilliToTime 毫秒戳转时间（默认时区）。对应 Java DateUtil.getLocalDateTime(long)。
func MilliToTime(millis int64) time.Time {
	return MilliToTimeIn(millis, GetDefaultZoneID())
}

// MilliToTimeIn 毫秒戳转时间（指定时区）。对应 Java DateUtil.getLocalDateTime(long, ZoneId)。
func MilliToTimeIn(millis int64, loc *time.Location) time.Time {
	if loc == nil {
		loc = GetDefaultZoneID()
	}
	return time.UnixMilli(millis).In(loc)
}

// IsBetweenDays 判断 date 是否满足 [before, last)。对应 Java DateUtil.isBetweenDays。
func IsBetweenDays(date, before, last time.Time) bool {
	d := GetMilliByTime(date)
	d1 := GetMilliByTime(before)
	d2 := GetMilliByTime(last)
	return d >= d1 && d < d2
}

// IsSameDay 判断两个时间是否同一天。对应 Java DateUtil.isSameDay。
func IsSameDay(t1, t2 time.Time) bool {
	loc := GetDefaultZoneID()
	a := t1.In(loc)
	b := t2.In(loc)
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameDayMillis 判断两个毫秒戳是否同一天。对应 Java DateUtil.isSameDay(long, long)。
func IsSameDayMillis(ms1, ms2 int64) bool {
	return IsSameDay(MilliToTime(ms1), MilliToTime(ms2))
}

// IsSameWeek 判断是否同一周（ISO 周）。对应 Java DateUtil.isSameWeek。
func IsSameWeek(ms1, ms2 int64) bool {
	t1 := MilliToTime(ms1)
	t2 := MilliToTime(ms2)
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()
	return y1 == y2 && w1 == w2
}

// IsSameMonth 判断是否同一月。对应 Java DateUtil.isSameMonth。
func IsSameMonth(ms1, ms2 int64) bool {
	t1 := MilliToTime(ms1)
	t2 := MilliToTime(ms2)
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// IsSameQuarter 判断是否同一季度（与 Java 一致：month/4）。对应 Java DateUtil.isSameQuarter。
func IsSameQuarter(ms1, ms2 int64) bool {
	t1 := MilliToTime(ms1)
	t2 := MilliToTime(ms2)
	return t1.Year() == t2.Year() && int(t1.Month())/4 == int(t2.Month())/4
}

// IsSameYear 判断是否同一年。对应 Java DateUtil.isSameYear。
func IsSameYear(ms1, ms2 int64) bool {
	return MilliToTime(ms1).Year() == MilliToTime(ms2).Year()
}

// GetDayOfMonth 获取月份中的日（1-31）。对应 Java DateUtil.getDayOfMonth(long)。
func GetDayOfMonth(millis int64) int {
	return MilliToTime(millis).Day()
}

// GetDayOfMonthNow 当前逻辑时间的日。对应 Java DateUtil.getDayOfMonth()。
func GetDayOfMonthNow() int {
	return GetDayOfMonth(CurrentTimeMillis())
}

// GetDayOfYear 获取一年中的日。对应 Java DateUtil.getDayOfYear(long)。
func GetDayOfYear(millis int64) int {
	return MilliToTime(millis).YearDay()
}

// GetDayOfYearNow 当前逻辑时间的年日。对应 Java DateUtil.getDayOfYear()。
func GetDayOfYearNow() int {
	return GetDayOfYear(CurrentTimeMillis())
}

// GetDayOfWeek 获取星期（1=周一 … 7=周日）。对应 Java DateUtil.getDayOfWeek(long)。
func GetDayOfWeek(millis int64) int {
	w := MilliToTime(millis).Weekday()
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

// GetDayOfWeekNow 当前逻辑时间的星期。对应 Java DateUtil.getDayOfWeek()。
func GetDayOfWeekNow() int {
	return GetDayOfWeek(CurrentTimeMillis())
}

// IsLeap 是否闰年。对应 Java DateUtil.isLeap。
func IsLeap(year int) bool {
	return time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay() == 366
}

// GetLastDayOfMonth 某年某月最后一天。对应 Java DateUtil.getLastDayOfMonth。
func GetLastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// resolveFormat 将 format 解析为 Go layout（缓存 + Java 常用 pattern 映射）。
func resolveFormat(format string) string {
	formatterMu.RLock()
	layout, ok := formatters[format]
	formatterMu.RUnlock()
	if ok {
		return layout
	}

	formatterMu.Lock()
	defer formatterMu.Unlock()
	if layout, ok = formatters[format]; ok {
		return layout
	}
	// 未知格式按 Go layout 原样使用并缓存
	formatters[format] = format
	return format
}
