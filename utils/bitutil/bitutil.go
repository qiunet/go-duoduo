package bitutil

// WriteBit 将 val 写入 ori 的 index 起始位（左移后按位或）。
// 对应 Java BitUtil.writeBit。
func WriteBit(ori, val, index int) int {
	return ori | (val << index)
}

// ReadBit 从 ori 的 index 起始位读取长度为 len 的位段。
// 对应 Java BitUtil.readBit。
func ReadBit(ori, index, length int) int {
	flag := 0
	for i := 0; i < length; i++ {
		flag |= 1 << i
	}
	ori >>= index
	return ori & flag
}

// SetBit 将 ori 的 index 位置为 1。
// 对应 Java BitUtil.setBit。
func SetBit(ori, index int) int {
	return ori | (1 << index)
}

// RemoveBit 将 ori 的 index 位清为 0。
// 对应 Java BitUtil.removeBit。
func RemoveBit(ori, index int) int {
	return ori & ^(1 << index)
}

// IsBitSet 判断 ori 的 index 位是否为 1。
// 对应 Java BitUtil.isBitSet。
func IsBitSet(ori, index int) bool {
	return (ori & (1 << index)) != 0
}
