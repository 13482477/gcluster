package http

import (
	"testing"
	"reflect"
	"fmt"
	"context"
)

type Sample struct {

}

func (s *Sample) Call(ctx context.Context, b float64) {

}

func TestMethodParamsType(t *testing.T) {
	sample := &Sample{}

	paramName1 := reflect.ValueOf(sample).MethodByName("Call").Type().In(0)

	fmt.Print(paramName1)
}