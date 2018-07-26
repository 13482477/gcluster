package endpoint

import (
	"fmt"
	"reflect"

	"context"

	"google.golang.org/grpc"
)

type Service struct {
	Name    string
	Client  interface{}
	Methods map[string]Method

	conn *grpc.ClientConn
}

func (s *Service) Call(ctx context.Context, locator Locator, req interface{}) (interface{}, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("no endpoint [%v]", locator)
	}

	if ep, ok := s.Methods[locator.Method]; ok {
		return ep.Call(ctx, s.getClient(s.conn), req)
	}

	return nil, fmt.Errorf("method not found [%v]", locator)
}

func (s *Service) GetMethodDetail(loc Locator) (*Method, bool) {
	m, ok := s.Methods[loc.Method]
	if !ok {
		return nil, false
	}

	return &m, ok
}

func (s *Service) getClient(conn *grpc.ClientConn) interface{} {
	args := []reflect.Value{
		reflect.ValueOf(conn),
	}
	result := reflect.ValueOf(s.Client).Call(args)
	return result[0].Interface()
}

func (s *Service) GetClient() (interface{}, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("no endpoint")
	}

	return s.getClient(s.conn), nil
}

func (s *Service) setConn(conn *grpc.ClientConn) {
	s.conn = conn
}
