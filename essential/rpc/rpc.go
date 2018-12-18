package rpc

import (
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/endpoint"
	"sync"
	"io"
	"strings"
	"net/url"
	"github.com/go-kit/kit/sd/lb"
	transportHttp "github.com/go-kit/kit/transport/http"
	consulsd "github.com/go-kit/kit/sd/consul"
	opentracinggo "github.com/opentracing/opentracing-go"
	"time"
	"github.com/go-kit/kit/tracing/opentracing"
	"fmt"
	"context"
	"errors"
	"net/http"
	applog "gcluster/essential/log"
)

var rpcManager *MCloudRpcManager
var rpcManagerOnce sync.Once

type MCloudRpcManager struct {
	EndpointMap map[string]map[string]*MCloudRpcOption
	InstanceMap map[string]sd.Instancer
	Client      consulsd.Client
	Tracer      opentracinggo.Tracer
	Tags        []string
	PassingOnly bool
	MaxAttempts int           // per request, before giving up
	MaxTime     time.Duration // wallclock time, before giving up
}

type MCloudRpcOption struct {
	ServiceName    string
	Path           string
	HttpMethod     string
	CreateReq      func() interface{}
	CreateResp     func() interface{}
	EncodeRequest  func(context context.Context, req *http.Request, request interface{}) error
	DecodeResponse func(context context.Context, resp *http.Response) (interface{}, error)
	Endpoint       endpoint.Endpoint
}

func GetMCloudRpcManager(client consulsd.Client, tracer opentracinggo.Tracer) *MCloudRpcManager {
	rpcManagerOnce.Do(func() {
		rpcManager = &MCloudRpcManager{
			EndpointMap: make(map[string]map[string]*MCloudRpcOption),
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

func GetRpcManager() *MCloudRpcManager {
	return rpcManager
}

func (manager *MCloudRpcManager) Subscript(rpcService *MCloudRpcOption) {
	manager.MakeRpcEndpoint(rpcService)

	var subEndpointMap map[string]*MCloudRpcOption
	if v, ok := manager.EndpointMap[rpcService.ServiceName]; ok {
		subEndpointMap = v
	} else {
		subEndpointMap = make(map[string]*MCloudRpcOption)
		manager.EndpointMap[rpcService.ServiceName] = subEndpointMap
	}

	subEndpointMap[rpcService.Path] = rpcService
}

func (manager *MCloudRpcManager) MakeRpcEndpoint(rpcOption *MCloudRpcOption) {
	var instance sd.Instancer
	if v, ok := manager.InstanceMap[rpcOption.ServiceName]; ok {
		instance = v
	} else {
		instance = consulsd.NewInstancer(manager.Client, applog.GetConsulLogger(), rpcOption.ServiceName, manager.Tags, manager.PassingOnly)
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

		return transportHttp.NewClient(rpcOption.HttpMethod, tgt, EncodeJSONRequest, MakeDecodeJsonResponse(rpcOption.CreateResp()), transportHttp.ClientBefore(opentracing.ContextToHTTP(manager.Tracer, applog.GetOpenTracingLogger()))).Endpoint(), nil, nil
	}, applog.GetEndpointLogger())

	balancer := lb.NewRoundRobin(defaultEndpoint)
	ep := lb.Retry(manager.MaxAttempts, manager.MaxTime, balancer)

	if manager.Tracer != nil {
		ep = opentracing.TraceClient(manager.Tracer, fmt.Sprintf("%s.%s", rpcOption.ServiceName, rpcOption.Path))(ep)
	}
	rpcOption.Endpoint = ep
}

func (manager *MCloudRpcManager) Call(ctx context.Context, server string, path string, req interface{}) (interface{}, error) {
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
