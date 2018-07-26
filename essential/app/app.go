package app

import (
	"fmt"

	"net"
	"os"
	"poseidon/essential/endpoint"
	"strings"

	"sync"

	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"

	"poseidon/essential/binder"

	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/spf13/viper"
)

var (
	appOnce     sync.Once
	appInstance *PoseidonApp
)

const (
	kWorkflowTypeSync  = 1
	kWorkflowTypeAsync = 2
)

type Workflow struct {
	Type    int
	Process func(app *PoseidonApp) error
}

type Manager interface {
	Start() error
}

type PoseidonApp struct {
	ServiceName     string
	Name            string
	Usage           string
	Version         string
	Config          interface{}
	ConfigWatcher   []ConfigWatcher
	Manager         *Manager
	ServiceRegistry *endpoint.ServiceRegistry
	Endpoint        *endpoint.EndPoint
	Workflow        []Workflow
}

func GetPoseidonApp() *PoseidonApp {
	appOnce.Do(func() {
		appInstance = &PoseidonApp{
			ConfigWatcher: make([]ConfigWatcher, 0),
		}
	})
	return appInstance
}

func (ppap *PoseidonApp) Run(workflow ...Workflow) error {
	ppap.Workflow = workflow

	app := &cli.App{
		Name:    ppap.Name,
		Usage:   ppap.Usage,
		Version: ppap.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config,c",
				Value: "config.toml",
				Usage: "config file path",
			},
			&cli.StringSliceFlag{
				Name:  "etcd,e",
				Value: cli.NewStringSlice("127.0.0.1:2379"),
				Usage: "ectd endpoints",
			},
			&cli.BoolFlag{
				Name:  "push_config,pc",
				Usage: "push config to etcd cluster",
			},
		},
		Action: func(ctx *cli.Context) error {

			// set log level to debug on start
			log.SetLevel(log.DebugLevel)

			configLoader := ConfigLoader{
				Name:          ppap.Name,
				FilePath:      ctx.String("config"),
				EtcdEndpoint:  ctx.StringSlice("etcd"),
				Config:        ppap.Config,
				ConfigWatcher: ppap.defaultConfigWatcher,
			}

			if ctx.Bool("push_config") {
				return configLoader.PushEtcdConfig(ctx.String("config"))
			}

			if err := configLoader.Load(); err != nil {
				log.Panic(err)
				return err
			}

			wg := sync.WaitGroup{}
			wg.Add(1)

			for _, wf := range ppap.Workflow {
				if wf.Type == kWorkflowTypeSync {
					if err := wf.Process(ppap); err != nil {
						log.Panic(err)
					}
				} else {
					localWorkflow := wf
					wg.Add(1)
					go func() {
						if err := localWorkflow.Process(ppap); err != nil {
							log.Panic(err)
						}
						wg.Done()
					}()
				}
			}

			wg.Done()
			wg.Wait()

			log.Info("RunPoseidon All Done")

			return nil
		},
	}

	return app.Run(os.Args)
}

func (ppap *PoseidonApp) startServiceRegistry() error {
	config := endpoint.ServiceRegistryOption{}
	if !viper.IsSet("ServiceRegistry") {
		log.Errorf("PoseidonApp startServiceRegistry ServiceRegistry not set")
		return errors.Errorf("PoseidonApp startServiceRegistry ServiceRegistry not set")
	}

	if err := viper.UnmarshalKey("ServiceRegistry", &config); err != nil {
		log.Errorf("unable to get serviceRegistryConfig")
		return errors.Wrapf(err, "startServiceRegistry unmarshal failed")
	}

	serverRegistry, err := endpoint.StartServiceRegistry(config)
	if err != nil {
		return err
	}

	ppap.ServiceRegistry = serverRegistry
	return nil
}

func (ppap *PoseidonApp) AddConfigWatcher(watcher ConfigWatcher) {
	ppap.ConfigWatcher = append(ppap.ConfigWatcher, watcher)
}

func (ppap *PoseidonApp) defaultConfigWatcher(config interface{}) error {
	for _, watcher := range ppap.ConfigWatcher {
		watcher(config)
	}
	return nil
}

func WithLogger() Workflow {
	setLogLevel := func(level string) error {
		l, err := log.ParseLevel(viper.GetString("LogLevel"))
		if err != nil {
			return err
		}
		log.SetLevel(l)
		return nil
	}

	return Workflow{
		Type: kWorkflowTypeSync,
		Process: func(ppap *PoseidonApp) error {

			ppap.AddConfigWatcher(func(_ interface{}) error {
				return setLogLevel(viper.GetString("LogLevel"))
			})

			return setLogLevel(viper.GetString("LogLevel"))
		},
	}
}

func WithLoggerFormatter(formatter log.Formatter) Workflow {
	return Workflow{
		Type: kWorkflowTypeSync,
		Process: func(ppap *PoseidonApp) error {
			if formatter != nil {
				log.SetFormatter(formatter)
			}
			return nil
		},
	}
}

func WithServiceRegistery() Workflow {
	return Workflow{
		Type: kWorkflowTypeSync,
		Process: func(ppap *PoseidonApp) error {
			if ppap.ServiceRegistry == nil {
				return ppap.startServiceRegistry()
			}
			return nil
		},
	}
}

func WithRegisterService(handler func(ppap *PoseidonApp) error) Workflow {
	return Workflow{
		Type: kWorkflowTypeSync,
		Process: func(ppap *PoseidonApp) error {
			if ppap.ServiceRegistry == nil {
				err := ppap.startServiceRegistry()
				if err != nil {
					return err
				}
			}
			ppap.Endpoint = endpoint.NewEndPoint(ppap.ServiceRegistry)
			return handler(ppap)
		},
	}
}

func WithPrometheusMetrics() Workflow {
	return Workflow{
		Type: kWorkflowTypeSync,
		Process: func(ppap *PoseidonApp) error {
			if ppap.ServiceRegistry == nil {
				err := ppap.startServiceRegistry()
				if err != nil {
					return err
				}
			}

			var port = 0
			if TCPServerConfig, ok := ppap.Config.(TCPServerConfig); ok {
				port = TCPServerConfig.GetTCPServerConfig().PrometheusPort
			}
			if err := ppap.ServiceRegistry.RunPrometheusMatrix(port); err != nil {
				return err
			}
			return nil
		},
	}
}

func WithManager(handler func(dbMap map[string]*xorm.Engine, poseidonApp *PoseidonApp) Manager) Workflow {
	return Workflow{
		Type: kWorkflowTypeSync,
		Process: func(poseidonApp *PoseidonApp) error {
			config, ok := poseidonApp.Config.(MysqlServerConfig)
			if !ok {
				log.Error("mysql server config is not MysqlServerConfig")
				return errors.Errorf("mysql server config is not MysqlServerConfig")
			}
			mysqlConfigMap := config.GetMysqlConfig()

			loggerLevelConfig := viper.GetString("LogLevel")
			dbMap := make(map[string]*xorm.Engine)
			for dbName, dbStr := range mysqlConfigMap.ConnectionString {
				dbEngine, err := xorm.NewEngine("mysql", dbStr)
				if err != nil {
					return err
				}

				if strings.ToUpper(loggerLevelConfig) == "DEBUG" {
					dbEngine.ShowSQL(true)
				}
				dbMap[dbName] = dbEngine
			}
			manager := handler(dbMap, poseidonApp)
			if manager == nil {
				return fmt.Errorf("build manager fail")
			}
			err := manager.Start()
			if err != nil {
				return err
			}
			poseidonApp.Manager = &manager
			return nil
		},
	}
}

func WithHttpServer(handler func(e *echo.Echo, ppap *PoseidonApp) error) Workflow {
	return Workflow{
		Type: kWorkflowTypeAsync,
		Process: func(ppap *PoseidonApp) error {
			TCPServerConfig, ok := ppap.Config.(TCPServerConfig)
			if !ok {
				log.Error("config is not TCPServerConfig")
				return errors.Errorf("config is not TCPServerConfig")
			}

			e := echo.New()
			e.Binder = &binder.PbRequestBinder{}

			e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{echo.GET, echo.POST},
			}))

			if err := handler(e, ppap); err != nil {
				return err
			}

			config := TCPServerConfig.GetTCPServerConfig()
			return e.Start(config.Host + ":" + strconv.Itoa(config.HttpPort))
		},
	}
}

func WithGrpcServer(handler func(server *grpc.Server, ppap *PoseidonApp) error) Workflow {
	return Workflow{
		Type: kWorkflowTypeAsync,
		Process: func(ppap *PoseidonApp) error {
			if ppap.ServiceRegistry == nil {
				err := ppap.startServiceRegistry()
				if err != nil {
					return err
				}
			}

			serverConfig, ok := ppap.Config.(TCPServerConfig)
			if !ok {
				log.Error("config is not TCPServerConfig")
				return errors.Errorf("config is not TCPServerConfig")
			}
			config := serverConfig.GetTCPServerConfig()
			address := fmt.Sprintf("%s:%d", config.Host, config.GrpcPort)
			listener, err := net.Listen("tcp", address)
			if err != nil {
				return err
			}
			server := ppap.ServiceRegistry.NewServer()
			if err := handler(server, ppap); err != nil {
				return err
			}
			return ppap.ServiceRegistry.RegisterServerAndRun(ppap.Name, server, listener)
		},
	}
}

func WithZipkinTracer() Workflow {

	getZipkinCreator := func(serviceName string) ConfigWatcher {
		return func(config interface{}) error {
			zipkinConfig := new(ZipkinConfig)
			if !viper.IsSet("ZipkinConfig") {
				log.Errorf("getZipkinCreator ZipkinConfig not set")
				return errors.Errorf("getZipkinCreator ZipkinConfig not set")
			}

			if err := viper.UnmarshalKey("ZipkinConfig", zipkinConfig); err != nil {
				log.Errorf("unable to get zipkin config")
				return errors.Wrapf(err, "getZipkinCreator unmarshal failed")
			}

			collector, err := zipkin.NewHTTPCollector(zipkinConfig.Url)
			if err != nil {
				log.Errorf("unable to create Zipkin HTTP collector: %+v\n", err)
				return err
			}

			recorder := zipkin.NewRecorder(collector, zipkinConfig.Debug, "0.0.0.0:0", serviceName)

			var sampler zipkin.Sampler
			switch zipkinConfig.Sampler {
			case "Mod":
				sampler = zipkin.ModuloSampler(uint64(zipkinConfig.Mod))
			case "Always":
				sampler = zipkin.NewBoundarySampler(1, 0)
			case "Never":
				sampler = zipkin.NewBoundarySampler(0, 0)
			default:
				sampler = zipkin.NewBoundarySampler(0, 0)
			}

			tracer, err := zipkin.NewTracer(
				recorder,
				zipkin.WithSampler(sampler),
				zipkin.ClientServerSameSpan(zipkinConfig.ClientServerSameSample),
				zipkin.TraceID128Bit(true),
				zipkin.DebugMode(zipkinConfig.Debug),
			)
			if err != nil {
				log.Errorf("unable to create Zipkin tracer: %+v\n", err)
				return err
			}

			opentracing.InitGlobalTracer(tracer)

			return nil
		}
	}

	return Workflow{
		Type: kWorkflowTypeSync,
		Process: func(ppap *PoseidonApp) error {
			if err := getZipkinCreator(ppap.Name)(ppap.Config); err != nil {
				return err
			}

			ppap.AddConfigWatcher(getZipkinCreator(ppap.Name))
			return nil
		},
	}
}

func WithCrontab(handler func(ppap *PoseidonApp) error) Workflow {
	return Workflow{
		Type: kWorkflowTypeSync,
		Process: func(ppap *PoseidonApp) error {
			if err := handler(ppap); err != nil {
				return err
			}
			return nil
		},
	}
}
