package MathUtil

import "testing"

func TestIsPowerOfTow(t *testing.T) {
	type args struct {
		val int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "2", args: args{2}, want: true},
		{name: "3", args: args{3}, want: false},
		{name: "16", args: args{16}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPowerOfTow(tt.args.val); got != tt.want {
				t.Errorf("IsPowerOfTow() = %v, want %v", got, tt.want)
			}
		})
	}
}
