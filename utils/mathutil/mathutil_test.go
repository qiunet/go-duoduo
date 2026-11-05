package mathutil

import (
	"math"
	"testing"
)

type weightItem struct {
	name string
	w    int
}

func (w weightItem) Weight() int { return w.w }

func TestIsPowerOfTwo(t *testing.T) {
	t.Run("是2的次幂", func(t *testing.T) {
		for _, v := range []int{0, 1, 2, 4, 8, 16, 32, 64, 128, 1024} {
			if !IsPowerOfTwo(v) {
				t.Fatalf("%d should be power of two", v)
			}
		}
	})
	t.Run("不是2的次幂", func(t *testing.T) {
		for _, v := range []int{5, 3, 1023} {
			if IsPowerOfTwo(v) {
				t.Fatalf("%d should not be power of two", v)
			}
		}
	})
}

func TestRandomIntRange(t *testing.T) {
	t.Run("半开区间落在范围内", func(t *testing.T) {
		for i := 1; i < 20; i++ {
			start, end := i, i*2
			for j := 0; j < 50; j++ {
				r := RandomIntRange(start, end)
				if r < start || r >= end {
					t.Fatalf("random(%d,%d)=%d out of range", start, end, r)
				}
			}
		}
	})
}

func TestRandomIntRangeTypeClosed(t *testing.T) {
	t.Run("闭区间包含端点", func(t *testing.T) {
		for j := 0; j < 100; j++ {
			r := RandomIntRangeType(1, 3, RandomTypeB)
			if r < 1 || r > 3 {
				t.Fatalf("closed random=%d", r)
			}
		}
	})
}

func TestHexBinary(t *testing.T) {
	t.Run("十六进制大写", func(t *testing.T) {
		if got := GetHexVal(255); got != "FF" {
			t.Fatalf("hex=%s", got)
		}
	})
	t.Run("十六进制左侧补零", func(t *testing.T) {
		if got := GetHexValPad(255, 4); got != "00FF" {
			t.Fatalf("hexPad=%s", got)
		}
	})
	t.Run("二进制字符串", func(t *testing.T) {
		if got := GetBinaryVal(6); got != "110" {
			t.Fatalf("bin=%s", got)
		}
	})
}

func TestGetByRateAndHit(t *testing.T) {
	t.Run("万分比取值", func(t *testing.T) {
		if got := GetByRate(10000, 100); got != 100 {
			t.Fatalf("rate=%d", got)
		}
	})
	t.Run("万分比命中边界", func(t *testing.T) {
		if IsHit(0) {
			t.Fatal("rate=0 should never hit")
		}
		if !IsHit(10000) {
			t.Fatal("rate=10000 should always hit")
		}
		for i := 0; i < 200; i++ {
			if IsHit(0) {
				t.Fatal("rate=0 should never hit")
			}
			if !IsHit(10000) {
				t.Fatal("rate=10000 should always hit")
			}
		}
	})
}

func TestBuildValXY(t *testing.T) {
	cases := []struct {
		name string
		x, y int
	}{
		{name: "正负混合", x: 12, y: -5},
		{name: "负x正y", x: -1, y: 7},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			v := BuildVal(tt.x, tt.y)
			if CalX(v) != tt.x || CalY(v) != tt.y {
				t.Fatalf("x=%d y=%d", CalX(v), CalY(v))
			}
		})
	}
}

func TestToInt(t *testing.T) {
	t.Run("正常转换", func(t *testing.T) {
		i, err := ToInt(42)
		if err != nil || i != 42 {
			t.Fatalf("i=%d err=%v", i, err)
		}
	})
	t.Run("超出32位范围报错", func(t *testing.T) {
		big := int64(math.MaxInt32) + 1
		if _, err := ToInt(big); err == nil {
			t.Fatal("expected overflow error")
		}
	})
}

func TestRandByWeight(t *testing.T) {
	t.Run("权重全落在唯一非零项", func(t *testing.T) {
		list := []weightItem{
			{name: "a", w: 0},
			{name: "b", w: 10},
			{name: "c", w: 0},
		}
		hit := map[string]int{}
		for i := 0; i < 200; i++ {
			obj, ok := RandByWeight(list)
			if !ok {
				t.Fatal("expected hit")
			}
			hit[obj.name]++
		}
		if hit["b"] != 200 {
			t.Fatalf("hit=%v", hit)
		}
	})
}
