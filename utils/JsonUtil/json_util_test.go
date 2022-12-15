package JsonUtil

import (
	"testing"
)

type user struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	Level   int
}

func TestToJsonString(t *testing.T) {
	type args struct {
		obj any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"user", args{obj: &user{"123456", "qiu", 10}}, `{"account":"123456","name":"qiu","Level":10}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJsonString(tt.args.obj); got != tt.want {
				t.Errorf("ToJsonString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGenerateOjb(t *testing.T) {
	data := `{"account":"123456","name":"qiu","Level":10}`
	p := new(user)
	GetGenerateOjb(data, p)
	if p.Level != 10 || p.Account != "123456" || p.Name != "qiu" {
		t.Fail()
	}
}
