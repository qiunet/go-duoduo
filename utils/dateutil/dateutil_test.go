package dateutil

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDateUtil_IsBetweenDays(t *testing.T) {
	t.Run("介于区间内", func(t *testing.T) {
		date1, _ := StringToDate("2016-05-16 00:00:00")
		date2, _ := StringToDate("2016-05-26 00:00:00")
		dt1, _ := StringToDate("2016-05-20 00:00:00")
		if !IsBetweenDays(dt1, date1, date2) {
			t.Fatal("should be between")
		}
	})
	t.Run("等于右边界不算", func(t *testing.T) {
		date1, _ := StringToDate("2016-05-16 00:00:00")
		date2, _ := StringToDate("2016-05-26 00:00:00")
		dt2, _ := StringToDate("2016-05-26 00:00:01")
		if IsBetweenDays(dt2, date1, date2) {
			t.Fatal("should not be between")
		}
	})
}

func TestDateUtil_FormatParseRoundTrip(t *testing.T) {
	t.Run("格式化与解析往返", func(t *testing.T) {
		val1 := "2016-05-20 00:00:00"
		val2 := "2016-05-26 00:00:01"
		dt1, err := StringToDate(val1)
		if err != nil {
			t.Fatal(err)
		}
		dt2, err := StringToDate(val2)
		if err != nil {
			t.Fatal(err)
		}
		time1 := GetMilliByTime(dt1)
		time2 := GetMilliByTime(dt2)
		if DateToStringMillis(time1) != val1 || DateToStringMillis(time2) != val2 {
			t.Fatal("format mismatch")
		}
		p1, _ := StringToDate(val1)
		p2, _ := StringToDate(val2)
		if GetMilliByTime(p1) != time1 || GetMilliByTime(p2) != time2 {
			t.Fatal("parse mismatch")
		}
	})
}

func TestDateUtil_ConcurrentFormat(t *testing.T) {
	t.Run("并发格式化解析", func(t *testing.T) {
		val1 := "2016-05-20 00:00:00"
		val2 := "2016-05-26 00:00:01"
		dt1, _ := StringToDate(val1)
		dt2, _ := StringToDate(val2)
		time1 := GetMilliByTime(dt1)
		time2 := GetMilliByTime(dt2)

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					if DateToStringMillis(time1) != val1 || DateToStringMillis(time2) != val2 {
						t.Error("format fail")
						return
					}
					p1, _ := StringToDate(val1)
					p2, _ := StringToDate(val2)
					if GetMilliByTime(p1) != time1 || GetMilliByTime(p2) != time2 {
						t.Error("parse fail")
						return
					}
				}
			}()
		}
		wg.Wait()
	})
}

func TestDateUtil_SameDay(t *testing.T) {
	t.Run("跨日与同日", func(t *testing.T) {
		date1, _ := StringToDate("2016-05-25 23:59:00")
		date2, _ := StringToDate("2016-05-26 00:01:00")
		date3, _ := StringToDate("2016-05-26 17:00:00")
		if IsSameDay(date1, date2) {
			t.Fatal("different day")
		}
		if !IsSameDay(date2, date3) {
			t.Fatal("same day")
		}
		if IsSameDayMillis(GetMilliByTime(date1), GetMilliByTime(date2)) {
			t.Fatal("millis different day")
		}
	})
}

func TestDateUtil_TimeOffset(t *testing.T) {
	t.Cleanup(ClearTimeOffset)

	t.Run("向前偏移", func(t *testing.T) {
		ClearTimeOffset()
		before := CurrentTimeMillis()
		if err := SetTimeOffset(10 * time.Minute); err != nil {
			t.Fatal(err)
		}
		after := CurrentTimeMillis()
		if after-before < int64((10*time.Minute)/time.Millisecond) {
			t.Fatalf("offset too small: %d", after-before)
		}
	})

	t.Run("拒绝非正偏移", func(t *testing.T) {
		ClearTimeOffset()
		if err := SetTimeOffset(0); err == nil {
			t.Fatal("expected error for 0")
		}
		if err := SetTimeOffset(-time.Second); err == nil {
			t.Fatal("expected error for negative")
		}
	})
}

func TestDateUtil_TimeChangeListener(t *testing.T) {
	t.Cleanup(ClearTimeOffset)

	t.Run("注册回调与取消", func(t *testing.T) {
		ClearTimeOffset()
		var n atomic.Int32
		cancel := AddTimeChangeListener(func() {
			n.Add(1)
		})
		t.Cleanup(cancel)

		if err := SetTimeOffset(time.Second); err != nil {
			t.Fatal(err)
		}
		if n.Load() != 1 {
			t.Fatalf("n=%d want 1", n.Load())
		}

		cancel()
		if err := SetTimeOffset(time.Second); err != nil {
			t.Fatal(err)
		}
		if n.Load() != 1 {
			t.Fatalf("after cancel n=%d want 1", n.Load())
		}
	})
}

func TestDateUtil_AddHour(t *testing.T) {
	t.Run("加减小时毫秒差", func(t *testing.T) {
		// 对齐 Java TestDateUtil.testAddHour
		const hours = 5
		now := Now()
		delta := int64(hours * 60 * 60 * 1000)
		if GetMilliByTime(now)-delta != GetMilliByTime(now.Add(-hours*time.Hour)) {
			t.Fatal("minus hours")
		}
		if GetMilliByTime(now)+delta != GetMilliByTime(now.Add(hours*time.Hour)) {
			t.Fatal("plus hours")
		}
	})
}

func TestDateUtil_ZoneRoundTrip(t *testing.T) {
	t.Run("默认时区与UTC毫秒接近", func(t *testing.T) {
		nowLocal := Now()
		nowUTC := NowIn(time.UTC)
		milliLocal := GetMilliByTime(nowLocal)
		milliUTC := GetMilliByTimeIn(nowUTC, time.UTC)
		diff := milliLocal - milliUTC
		if diff < 0 {
			diff = -diff
		}
		if diff >= 10 {
			t.Fatalf("diff=%d", diff)
		}
	})

	t.Run("格式化再解析后毫秒相等", func(t *testing.T) {
		// 对齐 Java TestDateUtil.testDateTime 后半段
		nowLocal := Now()
		nowUTC := NowIn(time.UTC)
		localParsed, err := StringToDate(DateToString(nowLocal))
		if err != nil {
			t.Fatal(err)
		}
		utcParsed, err := StringToDate(DateToString(nowUTC))
		if err != nil {
			t.Fatal(err)
		}
		milli1 := GetMilliByTime(localParsed)
		milli2 := GetMilliByTimeIn(utcParsed, time.UTC)
		if milli1 != milli2 {
			t.Fatalf("milli1=%d milli2=%d", milli1, milli2)
		}
	})

	t.Run("同一时间戳不同时区格式不同", func(t *testing.T) {
		nowMilli := CurrentTimeMillis()
		local8 := MilliToTimeIn(nowMilli, time.FixedZone("CST", 8*3600))
		utc := MilliToTimeIn(nowMilli, time.UTC)
		s1 := DateToString(local8)
		s2 := DateToString(utc)
		s3 := DateToString(utc.Add(8 * time.Hour))
		if s1 == s2 {
			t.Fatal("expected different string")
		}
		if s1 != s3 {
			t.Fatalf("s1=%s s3=%s", s1, s3)
		}
	})
}

func TestDateUtil_DayHelpers(t *testing.T) {
	t.Run("日周月年辅助", func(t *testing.T) {
		tm, _ := StringToDate("2016-05-16 00:00:00")
		ms := GetMilliByTime(tm)
		if GetDayOfMonth(ms) != 16 {
			t.Fatal("day of month")
		}
		if GetDayOfYear(ms) != tm.YearDay() {
			t.Fatal("day of year")
		}
		if GetDayOfWeek(ms) != 1 { // 2016-05-16 周一
			t.Fatalf("day of week=%d", GetDayOfWeek(ms))
		}
	})

	t.Run("闰年与月末", func(t *testing.T) {
		if !IsLeap(2020) || IsLeap(2019) {
			t.Fatal("leap")
		}
		if GetLastDayOfMonth(2020, time.February) != 29 {
			t.Fatal("feb 2020")
		}
		if GetLastDayOfMonth(2019, time.February) != 28 {
			t.Fatal("feb 2019")
		}
	})
}

func TestDateUtil_SamePeriod(t *testing.T) {
	t.Run("同月同年", func(t *testing.T) {
		a, _ := StringToDate("2016-05-01 00:00:00")
		b, _ := StringToDate("2016-05-20 12:00:00")
		c, _ := StringToDate("2016-06-01 00:00:00")
		ma, mb, mc := GetMilliByTime(a), GetMilliByTime(b), GetMilliByTime(c)
		if !IsSameMonth(ma, mb) || IsSameMonth(ma, mc) {
			t.Fatal("month")
		}
		if !IsSameYear(ma, mc) {
			t.Fatal("year")
		}
	})
}

func TestDateUtil_ConsumeMillis(t *testing.T) {
	t.Run("统计耗时非负", func(t *testing.T) {
		cost := CalConsumeMillisSeconds(func() {
			time.Sleep(5 * time.Millisecond)
		})
		if cost < 0 {
			t.Fatalf("cost=%d", cost)
		}
	})
}
