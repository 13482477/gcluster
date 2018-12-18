package http

import (
	"github.com/gorilla/mux"
	"sync"
	transport "github.com/go-kit/kit/transport/http"
	"gcluster/essential/manager"
	"github.com/go-kit/kit/tracing/opentracing"
	opentracinggo "github.com/opentracing/opentracing-go"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/log"
	"gcluster/essential/log"
	"gcluster/essential/metric"
)

var httpServer *MCloudHttpServer
var httpServerOnce sync.Once

type MCloudHttpServer struct {
	Router *mux.Router
	Tracer opentracinggo.Tracer
	Metric *metric.MCloudMetric
}

type MCloudHttpEndpointOption struct {
	Path        string
	HttpMethod  string
	Method      string
	CreateReq   func() interface{}
	CreateResp  func() interface{}
	ReqDecoder  transport.DecodeRequestFunc
	RespEncoder transport.EncodeResponseFunc
}

func GetHttpServer() *MCloudHttpServer {
	httpServerOnce.Do(func() {
		httpServer = &MCloudHttpServer{
			Router: mux.NewRouter(),
		}
	})
	return httpServer
}

func (httpServer *MCloudHttpServer) Register(manager manager.MCloudManager, endpointOption *MCloudHttpEndpointOption) *MCloudHttpServer {
	if endpointOption.ReqDecoder == nil && endpointOption.CreateReq == nil {
		log.Panic("You must filling the ReqDecoder or CreateReq at least one.")
	}

	endpoint := MakeEndpoint(manager, endpointOption.Method)
	reqDecoder := MakeRequestDecoder(endpointOption)
	respEncoder := MakeResponseEncoder(endpointOption)

	if httpServer.Metric != nil {
		endpoint = metric.MCloudMetricServer(endpointOption.Method, httpServer.Metric)(endpoint)
	}

	serverOptions := make([]transport.ServerOption, 0)
	serverOptions = append(serverOptions, transport.ServerErrorEncoder(ErrorEncoder))

	if httpServer.Tracer != nil {
		endpoint = opentracing.TraceServer(httpServer.Tracer, endpointOption.Method)(endpoint)
		serverOptions = append(serverOptions, transport.ServerBefore(opentracing.HTTPToContext(httpServer.Tracer, endpointOption.Path, applog.GetOpenTracingLogger())))
	}

	httpServer.Router.Handle(endpointOption.Path, transport.NewServer(
		endpoint,
		reqDecoder,
		respEncoder,
		serverOptions...
	)).Methods(endpointOption.HttpMethod)

	return httpServer
}
