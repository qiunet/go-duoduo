package Preconditions

import "errors"

// CheckState
// @param expression 表达式
// 检查状态. 如果expression 是false . 抛出异常.
func CheckState(expression bool, msg string) {
	if !expression {
		panic(errors.New(msg))
	}
}
