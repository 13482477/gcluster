package endpoint

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
)

type EndPoint struct {
	registry *ServiceRegistry
	//services map[string]*Service

	services sync.Map
}

func NewEndPoint(registry *ServiceRegistry) *EndPoint {
	return &EndPoint{
		registry: registry,
		//services: make(map[string]*Service),
	}
}

func (e *EndPoint) RegisterService(service string, s *Service) *EndPoint {
	return e.RegisterServiceWithLb(service, s, LBStrategyRR)
}

//rr round_robin
func (e *EndPoint) RegisterServiceWithLb(service string, s *Service, lb string) *EndPoint {

	if _, ok := e.services.Load(service); ok {
		log.Debugf("RegisterService duplicate registry %v", service)
		return e
	}

	conn, err := e.getConn(Locator{Service: service, Lb: lb})
	if err != nil {
		log.Errorf("RegisterService no endpoint:%v, looping and try to connect", service)
		ticker := time.NewTicker(1 * time.Second)
		go func() {
			for range ticker.C {
				conn, err := e.getConn(Locator{Service: service})
				if err == nil {
					s, ok := e.services.Load(service)
					if ok {
						serviceIns := s.(*Service)
						if serviceIns != nil {
							serviceIns.setConn(conn)
						}
					}
					ticker.Stop()
				} else {
					log.Infof("RegisterService no endpoint:%v in retry", service)
				}
			}
		}()
	} else {
		s.setConn(conn)
	}

	e.services.Store(service, s)
	return e
}

func (e *EndPoint) GetService(service string) *Service {
	s, ok := e.services.Load(service)
	if ok {
		serviceIns := s.(*Service)
		return serviceIns
	}
	return nil
}

func (e *EndPoint) GetMethodDetail(loc Locator) (*Method, bool) {
	s := e.GetService(loc.Service)
	if s == nil {
		log.Errorf("can'r find service, %v", loc)
		return nil, false
	}
	return s.GetMethodDetail(loc)
}

func (e *EndPoint) Call(ctx context.Context, locator Locator, req interface{}) (interface{}, error) {
	if service := e.GetService(locator.Service); service != nil {
		return service.Call(ctx, locator, req)
	}

	return nil, fmt.Errorf("service not found: %v", locator)
}

func (e *EndPoint) getConn(locator Locator) (*grpc.ClientConn, error) {
	// TODO 如果一开始并没有服务注册上来这里会导致再也没有进行服务发现
	return e.registry.ConnectClient(locator.Service, locator.Lb)
}
