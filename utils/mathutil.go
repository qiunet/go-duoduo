package utils

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// 游戏万分比基数。
const baseRate int64 = 10000

// RandomType 随机区间类型，对应 Java MathUtil.RandomType。
type RandomType int

const (
	// RandomTypeK 前闭后开 [start, end)
	RandomTypeK RandomType = iota
	// RandomTypeB 闭区间 [start, end]
	RandomTypeB
)

// WeightObj 带权重对象，对应 Java IWeightObj。
type WeightObj interface {
	// Weight 返回权重值。
	Weight() int
}

// RandomIntN 返回 [0, n) 的随机整数。对应 Java MathUtil.random(int)。
func RandomIntN(n int) int {
	if n <= 0 {
		return 0
	}
	return rand.IntN(n)
}

// RandomIntRange 返回 [start, end) 的随机整数。对应 Java MathUtil.random(start, end)。
func RandomIntRange(start, end int) int {
	return RandomIntRangeType(start, end, RandomTypeK)
}

// RandomIntRangeType 按区间类型返回随机整数。对应 Java MathUtil.random(start, end, type)。
func RandomIntRangeType(start, end int, randomType RandomType) int {
	if end-start < 0 {
		start, end = end, start
	} else if end-start == 0 {
		return start
	}
	if randomType == RandomTypeB {
		end++
	}
	return start + RandomIntN(end-start)
}

// RandomFloat64Range 返回 [start, end) 的随机浮点数。对应 Java MathUtil.random(double, double)。
func RandomFloat64Range(start, end float64) float64 {
	if end <= start {
		return start
	}
	return start + rand.Float64()*(end-start)
}

// RandomFloat32 返回 [0, i) 的随机 float32。对应 Java MathUtil.random(float)。
func RandomFloat32(i float32) float32 {
	return i * rand.Float32()
}

// IsPowerOfTwo 判断是否为 2 的次幂（与 Java 一致：0 视为 true）。
// 对应 Java MathUtil.isPowerOfTwo。
func IsPowerOfTwo(val int) bool {
	return (val & -val) == val
}

// GetHexVal 返回十六进制大写字符串。对应 Java MathUtil.getHexVal(long)。
func GetHexVal(val int64) string {
	return fmt.Sprintf("%X", uint64(val))
}

// GetHexValPad 返回固定位数的十六进制大写字符串（左侧 0 补齐）。
// 对应 Java MathUtil.getHexVal(long, count)。
func GetHexValPad(val int64, count int) string {
	return fmt.Sprintf("%0*X", count, uint64(val))
}

// GetBinaryVal 返回二进制字符串（按 32 位无符号，对齐 Java Integer.toBinaryString）。
// 对应 Java MathUtil.getBinaryVal。
func GetBinaryVal(val int) string {
	return fmt.Sprintf("%b", uint32(val))
}

// RandByWeight 按列表权重随机一个元素；列表为空或总权重为 0 时返回零值与 false。
// 对应 Java MathUtil.randByWeight(List)。
func RandByWeight[T WeightObj](list []T) (T, bool) {
	total := 0
	for _, obj := range list {
		total += obj.Weight()
	}
	return RandByWeightTotal(list, total)
}

// RandByWeightTotal 在 [0, totalRandNum) 上随机后按权重区间命中对象。
// 对应 Java MathUtil.randByWeight(List, totalRandNum)。
func RandByWeightTotal[T WeightObj](list []T, totalRandNum int) (T, bool) {
	var zero T
	if totalRandNum <= 0 || len(list) == 0 {
		return zero, false
	}
	randVal := RandomIntN(totalRandNum)
	start := 0
	for _, obj := range list {
		w := obj.Weight()
		if randVal >= start && randVal < start+w {
			return obj, true
		}
		start += w
	}
	return zero, false
}

// GetByRate 计算 val 的万分之 rate：(val * rate) / 10000。
// 对应 Java MathUtil.getByRate。
func GetByRate(val, rate int64) int64 {
	return (val * rate) / baseRate
}

// IsHit 万分比是否命中：在 [0, 10000) 随机后判断 <= rate。
// 对应 Java MathUtil.isHit。
func IsHit(rate int) bool {
	return RandomIntN(10000) <= rate
}

// ToInt 安全地将 int64 转为 int；超出 Java int（32 位）范围返回错误。
// 对应 Java MathUtil.toInt。
func ToInt(val int64) (int, error) {
	if val > math.MaxInt32 || val < math.MinInt32 {
		return 0, fmt.Errorf("utils: value %d out of int range", val)
	}
	return int(val), nil
}

// BuildVal 将 x、y 拼成一个 int64（高 32 位 x，低 32 位 y）。
// 对应 Java MathUtil.buildVal。
func BuildVal(x, y int) int64 {
	return (int64(x) << 32) | (int64(y) & 0xffffffff)
}

// CalX 从 BuildVal 结果中取出 x。对应 Java MathUtil.calX。
func CalX(val int64) int {
	return int(val >> 32)
}

// CalY 从 BuildVal 结果中取出 y（低 32 位有符号）。对应 Java MathUtil.calY。
func CalY(val int64) int {
	return int(int32(val))
}
