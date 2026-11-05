package jsonutil

import (
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	jsoniter.RegisterTypeEncoder("time.Time", unixTimeEncoder{})
	jsoniter.RegisterTypeDecoder("time.Time", unixTimeDecoder{})
}

type unixTimeEncoder struct{}

func (unixTimeEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return (*time.Time)(ptr).IsZero()
}

func (unixTimeEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	t := *(*time.Time)(ptr)
	if t.IsZero() {
		stream.WriteNil()
		return
	}
	stream.WriteInt64(t.Unix())
}

type unixTimeDecoder struct{}

func (unixTimeDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	switch iter.WhatIsNext() {
	case jsoniter.NilValue:
		iter.ReadNil()
		*(*time.Time)(ptr) = time.Time{}
	case jsoniter.NumberValue:
		sec := iter.ReadInt64()
		*(*time.Time)(ptr) = time.Unix(sec, 0)
	case jsoniter.StringValue:
		s := iter.ReadString()
		if s == "" {
			*(*time.Time)(ptr) = time.Time{}
			return
		}
		if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
			*(*time.Time)(ptr) = t
			return
		}
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			*(*time.Time)(ptr) = t
			return
		}
		iter.ReportError("decode time.Time", "invalid time string: "+s)
	default:
		iter.Skip()
		iter.ReportError("decode time.Time", "expect number, string or null")
	}
}

// Marshal marshals an object to JSON
func Marshal(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return jsonAPI.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	return jsonAPI.Unmarshal(data, v)
}

// ToJSONString 对象转 JSON 字符串；nil 不转，返回空字符串。
// time.Time 序列化为秒级时间戳；零值时间为 null。
func ToJSONString(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	b, err := jsonAPI.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSONString JSON 字符串转对象；空串 / "null" 不转。
// 时间戳（秒）或 RFC3339 字符串均可反序列化为 time.Time。
func FromJSONString(data string, v interface{}) error {
	if data == "" || data == "null" {
		return nil
	}
	return jsonAPI.Unmarshal([]byte(data), v)
}
