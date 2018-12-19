package rpc

import (
	"sync"
	"io"
	"strings"
	"net/url"
	"time"
	"fmt"
	"context"
	"errors"
	"net/http"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/sd/consul"
	"github.com/opentracing/opentracing-go"
	transportHttp "github.com/go-kit/kit/transport/http"
	goKitOpenTracing "github.com/go-kit/kit/tracing/opentracing"
	"gcluster/essential/log"
)

var rpcManager *GClusterRpcManager
var rpcManagerOnce sync.Once

type GClusterRpcManager struct {
	EndpointMap map[string]map[string]*GClusterRpcOption
	InstanceMap map[string]sd.Instancer
	Client      consul.Client
	Tracer      opentracing.Tracer
	Tags        []string
	PassingOnly bool
	MaxAttempts int
	MaxTime     time.Duration
}

type GClusterRpcOption struct {
	ServiceName    string
	Path           string
	HttpMethod     string
	CreateReq      func() interface{}
	CreateResp     func() interface{}
	EncodeRequest  func(context context.Context, req *http.Request, request interface{}) error
	DecodeResponse func(context context.Context, resp *http.Response) (interface{}, error)
	Endpoint       endpoint.Endpoint
}

func GetGClusterRpcManager(client consul.Client, tracer opentracing.Tracer) *GClusterRpcManager {
	rpcManagerOnce.Do(func() {
		rpcManager = &GClusterRpcManager{
			EndpointMap: make(map[string]map[string]*GClusterRpcOption),
			InstanceMap: make(map[string]sd.Instancer),
			Client:      client,
			Tracer:      tracer,
			Tags:        make([]string, 0),
			PassingOnly: true,
			MaxAttempts: 3,
			MaxTime:     250 * time.Millisecond,
		}
	})
	return rpcManager
}

func GetRpcManager() *GClusterRpcManager {
	return rpcManager
}

func (manager *GClusterRpcManager) Subscript(rpcService *GClusterRpcOption) {
	manager.MakeRpcEndpoint(rpcService)

	var subEndpointMap map[string]*GClusterRpcOption
	if v, ok := manager.EndpointMap[rpcService.ServiceName]; ok {
		subEndpointMap = v
	} else {
		subEndpointMap = make(map[string]*GClusterRpcOption)
		manager.EndpointMap[rpcService.ServiceName] = subEndpointMap
	}

	subEndpointMap[rpcService.Path] = rpcService
}

func (manager *GClusterRpcManager) MakeRpcEndpoint(rpcOption *GClusterRpcOption) {
	var instance sd.Instancer
	if v, ok := manager.InstanceMap[rpcOption.ServiceName]; ok {
		instance = v
	} else {
		instance = consul.NewInstancer(manager.Client, applog.GetConsulLogger(), rpcOption.ServiceName, manager.Tags, manager.PassingOnly)
		manager.InstanceMap[rpcOption.ServiceName] = instance
	}

	defaultEndpoint := sd.NewEndpointer(instance, func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}
		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = rpcOption.Path

		return transportHttp.NewClient(rpcOption.HttpMethod, tgt, EncodeJSONRequest, MakeDecodeJsonResponse(rpcOption.CreateResp()), transportHttp.ClientBefore(goKitOpenTracing.ContextToHTTP(manager.Tracer, applog.GetOpenTracingLogger()))).Endpoint(), nil, nil
	}, applog.GetEndpointLogger())

	balancer := lb.NewRoundRobin(defaultEndpoint)
	ep := lb.Retry(manager.MaxAttempts, manager.MaxTime, balancer)

	if manager.Tracer != nil {
		ep = goKitOpenTracing.TraceClient(manager.Tracer, fmt.Sprintf("%s.%s", rpcOption.ServiceName, rpcOption.Path))(ep)
	}
	rpcOption.Endpoint = ep
}

func (manager *GClusterRpcManager) Call(ctx context.Context, server string, path string, req interface{}) (interface{}, error) {
	if subMap, ok := manager.EndpointMap[server]; ok {
		if rpcService, ok := subMap[path]; ok {
			if resp, err := rpcService.Endpoint(ctx, req); err != nil {
				return nil, errors.New(fmt.Sprintf("rpc call failed, serviceName=%s, path=%s, error=%v", server, path, err))
			} else {
				return resp, nil
			}
		} else {
			return nil, errors.New(fmt.Sprintf("path no found, path not be registered, server=%s, path=%s", server, path))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("server no found, server not be registered, server=%s, path=%s", server, path))
	}
}
