package MathUtil

import (
	"github.com/qiunet/go-duoduo/utils/Preconditions"
	"math/rand"
	"runtime"
)

func init() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	rand.Seed(int64(stats.Frees))
}

// RandomInt
// 前开后闭的随机组
func RandomInt(start, end int) int {
	Preconditions.CheckState(end > start, "start can not great than end!")
	return rand.Intn(end-start) + start
}

// IsPowerOfTow
// 是否是2的次幂
func IsPowerOfTow(val int) bool {
	return (val & -val) == val
}
