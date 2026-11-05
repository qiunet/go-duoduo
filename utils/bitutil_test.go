package utils

import "testing"

func TestBitUtil(t *testing.T) {
	t.Run("写入并按位段读取", func(t *testing.T) {
		ss1 := WriteBit(0, 13, 0) // 13 -> 1101
		ss2 := WriteBit(ss1, 142342, 4)
		if got := ReadBit(ss2, 0, 4); got != 13 {
			t.Fatalf("ReadBit low=%d", got)
		}
		if got := ReadBit(ss2, 4, 28); got != 142342 {
			t.Fatalf("ReadBit high=%d", got)
		}
	})

	t.Run("置位与判断", func(t *testing.T) {
		ss3 := SetBit(0, 12)
		if !IsBitSet(ss3, 12) {
			t.Fatal("bit 12 should be set")
		}
	})

	t.Run("清位", func(t *testing.T) {
		ss3 := SetBit(0, 12)
		ss4 := SetBit(ss3, 1)
		ss5 := SetBit(ss4, 2)
		ss6 := RemoveBit(ss5, 12)
		if ss6 != 6 {
			t.Fatalf("ss6=%d want 6", ss6)
		}
		if IsBitSet(ss6, 12) {
			t.Fatal("bit 12 should be cleared")
		}
	})
}
