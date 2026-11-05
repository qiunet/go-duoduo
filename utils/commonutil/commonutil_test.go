package commonutil

import "testing"

func TestCommonUtil_ExistAndSlice(t *testing.T) {
	t.Run("存在于可变参数", func(t *testing.T) {
		// 对齐 Java TestCommonUtil.testExistInArray
		arr := []int{1, 2, 3, 4, 5}
		if !ExistInList(3, arr) || ExistInList(8, arr) {
			t.Fatal("existInArray")
		}
		if !ExistInSlice("a", "x", "a", "y") || ExistInSlice("z", "x", "a") {
			t.Fatal("existInVararg")
		}
	})
	t.Run("存在于切片", func(t *testing.T) {
		// 对齐 Java TestCommonUtil.testExistInCollection
		list := []int{0, 1, 2, 3, 4}
		if !ExistInList(3, list) || ExistInList(8, list) {
			t.Fatal("existInCollection")
		}
		if ExistInList(9, []int{}) {
			t.Fatal("empty")
		}
	})
	t.Run("分页截取", func(t *testing.T) {
		// 对齐 Java TestCommonUtil.testSubList：0..9 skip=1 count=4 → [1,2,3,4]
		list := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		got := GetSubListPage(list, 1, 4)
		if len(got) != 4 || got[0] != 1 || got[len(got)-1] != 4 {
			t.Fatalf("got=%v", got)
		}
		if GetSubListPage([]int(nil), 0, 1) != nil {
			t.Fatal("nil list")
		}
		if len(GetSubListPage(list, 100, 2)) != 0 {
			t.Fatal("skip beyond")
		}
	})
}

func TestCommonUtil_Reverse(t *testing.T) {
	t.Run("翻转字节", func(t *testing.T) {
		// 对齐 Java TestCommonUtil.testReverse（含奇数长度）
		arr1 := []byte{1, 2, 3, 4}
		arr2 := []byte{1, 2, 3, 4, 5}
		ReverseBytes(arr1)
		ReverseBytes(arr2)
		want1 := []byte{4, 3, 2, 1}
		want2 := []byte{5, 4, 3, 2, 1}
		for i := range want1 {
			if arr1[i] != want1[i] {
				t.Fatalf("arr1=%v", arr1)
			}
		}
		for i := range want2 {
			if arr2[i] != want2[i] {
				t.Fatalf("arr2=%v", arr2)
			}
		}
	})
	t.Run("翻转整型", func(t *testing.T) {
		a := []int{1, 2, 3}
		ReverseInts(a)
		if a[0] != 3 || a[2] != 1 {
			t.Fatalf("%v", a)
		}
	})
}

func TestCommonUtil_Empty(t *testing.T) {
	t.Run("空判断", func(t *testing.T) {
		if !IsEmptyMap[string, int](nil) || IsNotEmptyMap(map[string]int{"a": 1}) == false {
			t.Fatal("map")
		}
		if !IsEmptySlice[int](nil) || !IsNotEmptySlice([]int{1}) {
			t.Fatal("slice")
		}
	})
}
