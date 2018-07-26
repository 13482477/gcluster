package endpoint

import (
	"context"
	"reflect"
)

type Method struct {
	Method         interface{}
	CreateRequest  func() interface{}
	CreateResponse func() interface{}
}

func (ep *Method) Call(ctx context.Context, client interface{}, req interface{}) (interface{}, error) {
	args := []reflect.Value{
		reflect.ValueOf(client),
		reflect.ValueOf(ctx),
		reflect.ValueOf(req),
	}
	result := reflect.ValueOf(ep.Method).Call(args)

	resultError := result[1].Interface()
	if resultError != nil {
		return result[0].Interface(), resultError.(error)
	} else {
		return result[0].Interface(), nil
	}
}
