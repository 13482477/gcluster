package endpoint

import (
	"github.com/coreos/etcd/clientv3"
	etcd "github.com/coreos/etcd/clientv3"

	"context"

	"time"

	"net"

	"net/http"

	"os"

	"strings"

	"reflect"

	"encoding/json"

	"fmt"

	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	 grpclb "poseidon/essential/balance"
	"google.golang.org/grpc/codes"
)

type ServiceRegistryOption struct {
	Endpoints   []string
	DialTimeout time.Duration
}

type ServiceRegistry struct {
	client *clientv3.Client
}

func StartServiceRegistry(option ServiceRegistryOption) (*ServiceRegistry, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   option.Endpoints,
		DialTimeout: time.Second * option.DialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceRegistry{
		client: etcdClient,
	}, nil
}

func (service *ServiceRegistry) GetConfig(ctx context.Context, key string, config interface{}) error {
	response, err := service.client.Get(ctx, key)
	if err != nil {
		return err
	}

	if len(response.Kvs) == 0 {
		return errors.Errorf("no config set for key:%s", key)
	}

	return json.Unmarshal(response.Kvs[0].Value, config)
}

func (service *ServiceRegistry) ConnectClient(name string, lb string) (*grpc.ClientConn, error) {
	resolver := etcdnaming.GRPCResolver{Client: service.client}

	serviceName := "/services/"+name
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	balancer := grpc.RoundRobin(&resolver)
	if lb == LBStrategyHash {
		serviceName = ""
		etcdConfg := etcd.Config{
			Endpoints: service.client.Endpoints(),
		}
		r := grpclb.NewResolver("/services",  name, etcdConfg)
		balancer = grpclb.NewBalancer(r, grpclb.NewKetamaSelector(grpclb.DefaultKetamaKey))
	}

	client, err := grpc.DialContext(
		ctx,
		serviceName,
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithBalancer(balancer),

		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_opentracing.UnaryClientInterceptor(),
				grpc_prometheus.UnaryClientInterceptor,
			)),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				grpc_opentracing.StreamClientInterceptor(),
				grpc_prometheus.StreamClientInterceptor,
			)),
	)

	if err != nil {
		return nil, errors.Wrapf(err, "ServiceRegistry ConnectClient fail")
	}

	return client, nil
}

func (service *ServiceRegistry) NewServer(opt ...grpc.ServerOption) *grpc.Server {

	options := make([]grpc.ServerOption, 0)
	for _, o := range opt {
		options = append(options, o)
	}

	options = append(options, grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_logrus.StreamServerInterceptor(log.NewEntry(log.New())),
			grpc_recovery.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
		)))
	options = append(options, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(log.NewEntry(log.New()), grpc_logrus.WithLevels(func(code codes.Code) log.Level {
				if code == codes.OK {
					return log.DebugLevel
				}
				return grpc_logrus.DefaultCodeToLevel(code)
			})),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		)))

	server := grpc.NewServer(options...)
	grpc_prometheus.Register(server)
	return server
}

func (service *ServiceRegistry) RegisterServerAndRun(name string, server *grpc.Server, listener net.Listener) error {
	addr := service.convertTcpAddress(listener.Addr())

	if err := service.RegisterServiceToEtcd(name, addr); err != nil {
		return err
	}

	if err := service.RegisterServiceToPrometheus(server); err != nil {
		return err
	}

	log.Info("======================================")
	log.Infof("==  service start on : %s", addr)
	log.Info("======================================")
	return server.Serve(listener)
}

func (service *ServiceRegistry) RegisterServiceToPrometheus(server *grpc.Server) error {
	grpc_prometheus.Register(server)
	return nil
}

func (service *ServiceRegistry) RegisterServiceToEtcd(name string, addr net.Addr) error {
	lease := clientv3.NewLease(service.client)

	leaseResp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		return errors.Wrapf(err, "RegisterServiceToEtcd failed")
	}

	serviceName := "/services/" + name
	namingUpdate := naming.Update{
		Op:   naming.Add,
		Addr: addr.String(),
	}

	resolver := etcdnaming.GRPCResolver{Client: service.client}
	err = resolver.Update(
		context.TODO(),
		serviceName,
		namingUpdate,
		clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return errors.Wrapf(err, "RegisterServiceToEtcd resolver failed")
	}

	log.Infof("ServiceRegistry RegisterServiceToEtcd [%s] [%s]", serviceName, namingUpdate)

	lease.KeepAlive(context.TODO(), leaseResp.ID)

	return nil
}

func (service *ServiceRegistry) RunPrometheusMatrix(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return errors.Wrapf(err, "RunPrometheusMatrix listen failed")
	}

	grpc_prometheus.EnableHandlingTimeHistogram()

	go func() {
		log.Info("======================================")
		log.Infof("==  Metrics start on : %s", listener.Addr())
		log.Info("======================================")

		http.Handle("/metrics", promhttp.Handler())
		panic(http.Serve(listener, nil))
	}()

	return service.RegisterServiceToEtcd("prometheus_metrics", service.convertTcpAddress(listener.Addr()))
}

func (service *ServiceRegistry) convertTcpAddress(addr net.Addr) net.Addr {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return addr
	}

	if !tcpAddr.IP.IsUnspecified() {
		return addr
	}

	host, err := os.Hostname()
	if err != nil {
		log.Errorf("ServiceRegistry convertTcpAddr Hostname fail, err:%s", err)
		return addr
	}

	ips := make([]net.IP, 0)

	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorf("ServiceRegistry convertTcpAddr InterfaceAddrs fail, err:%s", err)
		return addr
	}

	log.Debugf("ServiceRegistry convertTcpAddr InterfaceAddrs host:[%s], addr:[%s]", host, interfaceAddrs)

	ips = make([]net.IP, 0)
	for _, addr := range interfaceAddrs {
		if ipAddr, ok := addr.(*net.IPNet); ok {
			ips = append(ips, ipAddr.IP)
		} else {
			log.Debugf("bot ip addr addr:[%s], type:[%s]", addr, reflect.TypeOf(addr))
		}
	}

	log.Debugf("ServiceRegistry convertTcpAddr LookupIP host:[%s], ip:[%s]", host, ips)

	found := false
	for _, ip := range ips {
		if ip.IsUnspecified() || ip.IsLoopback() {
			continue
		}

		if !strings.HasPrefix(ip.String(), "10.215.") &&
			!strings.HasPrefix(ip.String(), "192.168.") {
			continue
		}

		tcpAddr.IP = ip
		found = true
		break
	}

	if !found {
		log.Errorf("ServiceRegistry convertTcpAddr LookupIP no suitable ip found, host:%s", ips)
		return addr
	}

	return tcpAddr
}
