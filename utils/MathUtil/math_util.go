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

func RandomInt(start, end int) int {
	Preconditions.CheckState(end > start, "start can not great than end!")
	return rand.Intn(end-start) + start
}
