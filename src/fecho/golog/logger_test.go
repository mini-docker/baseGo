package golog

import (
	"testing"
	"time"
)

type StructTest struct {
	FieldA int               `json:"fieldA"`
	FieldB string            `json:"fieldB"`
	FieldC *StructTest       `json:"fieldC"`
	FieldD map[string]string `json:"fieldD"`
}

func TestLoggingInterface_Info(t *testing.T) {
	s := &StructTest{
		FieldA: 0,
		FieldB: "fieldb",
		FieldC: &StructTest{
			FieldA: 0,
			FieldB: "b",
			FieldC: nil,
			FieldD: map[string]string{
				"a": "a",
				"b": "b",
			},
		},
		FieldD: map[string]string{
			"aaa": "aaa",
			"bbb": "bbb",
		},
	}

	Info("xxx", "xxx", "struct", s)
	time.Sleep(2 * time.Second)
}
