package commonutil

// ExistInSlice 检查元素是否在可变参数列表中。对应 Java CommonUtil.existInList(T, T...)。
func ExistInSlice[T comparable](element T, arrays ...T) bool {
	for _, v := range arrays {
		if v == element {
			return true
		}
	}
	return false
}

// ExistInList 检查元素是否在切片中。对应 Java CommonUtil.existInList(T, Collection)。
func ExistInList[T comparable](element T, list []T) bool {
	if len(list) == 0 {
		return false
	}
	for _, v := range list {
		if v == element {
			return true
		}
	}
	return false
}

// GetSubListPage 截取切片的一部分（从 skip 开始取 count 个）。
// list 为 nil/空时返回 nil。对应 Java CommonUtil.getSubListPage。
func GetSubListPage[T any](list []T, skip, count int) []T {
	if len(list) == 0 {
		return nil
	}
	if skip < 0 {
		skip = 0
	}
	if skip >= len(list) || count <= 0 {
		return []T{}
	}
	end := skip + count
	if end > len(list) {
		end = len(list)
	}
	out := make([]T, end-skip)
	copy(out, list[skip:end])
	return out
}

// ReverseBytes 翻转字节数组。对应 Java CommonUtil.reverse(byte[])。
func ReverseBytes(array []byte) {
	ReverseBytesStep(array, 1)
}

// ReverseBytesStep 按 step 翻转字节数组。对应 Java CommonUtil.reverse(byte[], step)。
func ReverseBytesStep(array []byte, step int) {
	if len(array) <= 1 {
		return
	}
	if step <= 0 {
		step = 1
	}
	loopNum := len(array) / 2
	for i := 0; i < loopNum; i += step {
		last := len(array) - 1 - i
		array[i], array[last] = array[last], array[i]
	}
}

// ReverseInts 翻转整型切片。对应 Java CommonUtil.reverse(int[])。
func ReverseInts(array []int) {
	if len(array) <= 1 {
		return
	}
	loopNum := len(array) / 2
	for i := 0; i < loopNum; i++ {
		last := len(array) - 1 - i
		array[i], array[last] = array[last], array[i]
	}
}

// IsEmptyMap 判断 map 是否为空（含 nil）。对应 Java CommonUtil.isEmptyMap。
func IsEmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

// IsNotEmptyMap 判断 map 是否非空。对应 Java CommonUtil.isNotEmptyMap。
func IsNotEmptyMap[K comparable, V any](m map[K]V) bool {
	return !IsEmptyMap(m)
}

// IsEmptySlice 判断切片是否为空（含 nil）。对应 Java CommonUtil.isEmpty。
func IsEmptySlice[T any](collection []T) bool {
	return len(collection) == 0
}

// IsNotEmptySlice 判断切片是否非空。对应 Java CommonUtil.isNotEmpty。
func IsNotEmptySlice[T any](collection []T) bool {
	return !IsEmptySlice(collection)
}
