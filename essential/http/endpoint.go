package http

import (
	"reflect"
	"gcluster/essential/manager"
	"github.com/go-kit/kit/endpoint"
	log "github.com/sirupsen/logrus"
	"context"
)

func MakeEndpoint(manager manager.GClusterManager, method string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Debugf("manager method invoke, manager=%s, method=%s", reflect.TypeOf(manager), method)

		target := reflect.ValueOf(manager)

		params := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(request),
		}

		results := target.MethodByName(method).Call(params)
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		} else {
			return results[0].Interface(), nil
		}
	}
}
