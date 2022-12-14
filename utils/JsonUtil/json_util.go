package JsonUtil

import (
	"encoding/json"
	"log"
)

// ToJsonString
// 转为json string
func ToJsonString(obj any) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// GetGenerateOjb 获得一个对象
func GetGenerateOjb(data string, v any) {
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		log.Fatal(err)
	}
}
