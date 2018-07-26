package utils

import (
	"testing"
)

type TestStruct struct {
	String string
	Int    int
	IntA   int8
	IntB   int16
	IntC   int32
	IntD   int64
	Bool   bool
	FloatA float32
	FloatB float64
	Empty  int
}

func TestSmartFillMap(t *testing.T) {
	target := &TestStruct{}

	data := map[string]string{
		"string":  "string",
		"int":     "1",
		"int_a":   "2",
		"int_b":   "3",
		"int_c":   "4",
		"int_d":   "5",
		"bool":    "true",
		"float_a": "1.1",
		"float_b": "2.2",
	}

	err := SmartFillMap(target, data)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(target)
}
