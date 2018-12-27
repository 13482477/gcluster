package app

import (
	"sync"
	"os"
	"time"
	"net/http"
	"fmt"
	"github.com/urfave/cli"
	"github.com/opentracing/opentracing-go"
	"github.com/hashicorp/consul/api"
	"github.com/go-kit/kit/sd/consul"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron"
	"github.com/rifflock/lfshook"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
	"gcluster/essential/manager"
	"gcluster/essential/config"
	"gcluster/essential/metric"
	"gcluster/essential/rpc"
	"gcluster/essential/registry"
	gHttp "gcluster/essential/http"
	gCron "gcluster/essential/cron"
)

var gClusterApp *GClusterApp
var gClusterAppOnce sync.Once

type RunType int

const (
	RunTypeSync  RunType = 1
	RunTypeAsync RunType = 2
)

type RunOption struct {
	Type    RunType
	Process func(mcApp *GClusterApp) error
}

type GClusterApp struct {
	Name       string
	Usage      string
	Version    string
	Config     config.GClusterConfig
	Metric     *metric.GClusterMetric
	Manager    manager.GClusterManager
	Client     consul.Client
	Registry   *registry.GClusterServiceRegistry
	Tracer     opentracing.Tracer
	RpcManager *rpc.GClusterRpcManager
	HttpServer *gHttp.GClusterHttpServer
	RunOptions []*RunOption
}

func GetGClusterApp() *GClusterApp {
	gClusterAppOnce.Do(func() {
		gClusterApp = &GClusterApp{}
	})
	return gClusterApp
}

func printLogo() {
	log.Infof("                                                                                        ")
	log.Infof("                                                                                        ")
	log.Infof("                     &&         &&    &&&   &      &&   &   &  &&&                      ")
	log.Infof("                    && &     & &&    &      &     &  &  &   &  &  &                     ")
	log.Infof("                   &&  &   &  &&     &      &     &  &  &   &  &  &                     ")
	log.Infof("                  &&   & &   &&      &      &     &  &  &   &  &  &                     ")
	log.Infof("                 &&    &    &&        &&&   &&&&   &&    &&&   &&&                      ")
	log.Infof("                                                                                        ")
	log.Infof("                                                                                        ")
}

func (gApp *GClusterApp) Run(runOptions ...*RunOption) error {
	gApp.RunOptions = runOptions

	app := &cli.App{
		Name:    gApp.Name,
		Usage:   gApp.Usage,
		Version: gApp.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config, c",
				Value: "config",
				Usage: "config file path",
			},
		},
		Action: func(ctx *cli.Context) error {

			filePath := fmt.Sprintf("./%s.log", gApp.Name)
			hook := lfshook.NewHook(filePath, nil)
			log.AddHook(hook)

			log.SetLevel(log.DebugLevel)

			printLogo()

			log.Infof("========================================================================================")
			log.Infof("======================================System start======================================")
			log.Infof("========================================================================================")
			log.WithField("SystemName", gApp.Name).Info()
			log.WithField("Version", gApp.Version).Info()

			configLoader := &config.GClusterConfigLoader{
				Name:     ctx.String("config"),
				FilePath: ".",
				Config:   gApp.Config,
			}

			if err := configLoader.Load(); err != nil {
				log.Panicf("load config file failed, error=%v", err)
			}

			wg := sync.WaitGroup{}
			wg.Add(1)

			for _, runOption := range gApp.RunOptions {
				if runOption.Type == RunTypeSync {
					if err := runOption.Process(gApp); err != nil {
						log.Panic(err)
					}
				} else {
					localRunOption := runOption
					wg.Add(1)
					go func() {
						if err := localRunOption.Process(gApp); err != nil {
							log.Panic(err)
						}
						wg.Done()
					}()
				}
			}

			wg.Done()
			wg.Wait()

			return nil
		},
	}

	return app.Run(os.Args)
}

func WithLoggerOption() *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *GClusterApp) error {
			logConfig := app.Config.(config.LogConfiguration)
			level, err := log.ParseLevel(logConfig.GetLogLevelConfig())
			if err != nil {
				log.Panicf("Parse log level failed, error=%v", err)
			}
			log.SetLevel(level)
			if log.DebugLevel == level {
				//log.SetReportCaller(true)
			}

			return nil
		},
	}
}

func WithMetricOption() *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *GClusterApp) error {

			fieldKeys := []string{"method"}

			app.Metric = &metric.GClusterMetric{
				RequestCount: kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
					Namespace: "gcluster",
					Subsystem: app.Name,
					Name:      "request_count",
					Help:      "Number of requests received.",
				}, fieldKeys),

				RequestLatency: kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
					Namespace: "gcluster",
					Subsystem: app.Name,
					Name:      "request_latency_microseconds",
					Help:      "Total duration of requests in microseconds.",
				}, fieldKeys),
			}

			log.Printf("Start GCluster metric server successfully!")
			return nil
		},
	}
}

func WithRegistryOption() *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *GClusterApp) error {

			cfg, ok := app.Config.(config.ServiceRegistryConfiguration)
			if !ok {
				log.Panic("ServiceRegistryConfiguration is not be configured.")
			}

			consulConfig := api.DefaultConfig()
			consulConfig.Address = cfg.GetServiceRegistryConfig().Address
			consulClient, err := api.NewClient(consulConfig)
			if err != nil {
				log.Panic("Failed to init consul client.")
			}

			gClusterServiceRegistry := registry.GClusterServiceRegistry{
				Client:     consul.NewClient(consulClient),
				ServerName: app.Name,
				Config:     app.Config,
			}

			app.Client = gClusterServiceRegistry.Client

			gClusterServiceRegistry.Register()
			log.Printf("Start GCluster server gClusterServiceRegistry successfully!")
			return nil
		},
	}
}

func WithOpenTracingOption() *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *GClusterApp) error {
			tracer := opentracing.GlobalTracer()
			app.Tracer = tracer
			log.Printf("Start GCluster trace server successfully!")
			return nil
		},
	}
}

func WithManagerOption(handler func(db *gorm.DB) (manager.GClusterManager, error)) *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *GClusterApp) error {
			dbConfig := app.Config.(config.DatabaseConfiguration).GetDataBaseConfig()

			var db *gorm.DB

			if dbConfig.Address != "" {
				database, err := gorm.Open("mysql", dbConfig.Address)
				if err != nil {
					log.Panicf("Fatal error mysql connection failed: %v", err)
				}
				db = database
				db.SingularTable(true)
				db.LogMode(dbConfig.LogMode)
				db.DB().SetMaxIdleConns(dbConfig.MaxIdle)
				db.DB().SetMaxOpenConns(dbConfig.MaxConns)
				db.DB().SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetime) * time.Second)
			}

			if gClusterManager, err := handler(db); err != nil {
				return err
			} else {
				app.Manager = gClusterManager
				if err := gClusterManager.StartGClusterManager(); err != nil {
					log.Panicf("Start GCluster manager failed, error=%v", err)
				} else {
					log.Infof("Start GCluster manager successfully!")
				}
				return nil
			}
		},
	}
}

func WithRpcOption(handler func() []*rpc.GClusterRpcOption) *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *GClusterApp) error {
			rpcManager := rpc.GetGClusterRpcManager(app.Client, app.Tracer)
			app.RpcManager = rpcManager

			options := handler()

			for _, v := range options {
				rpcManager.Subscript(v)
			}

			log.Printf("Start GCluster rcp server successfully!")
			return nil
		},
	}
}

func WithHttpEndpointOption(handler func() []*gHttp.GClusterHttpEndpointOption) *RunOption {
	return &RunOption{
		Type: RunTypeAsync,
		Process: func(app *GClusterApp) error {
			app.HttpServer = gHttp.GetHttpServer()
			app.HttpServer.Tracer = app.Tracer
			app.HttpServer.Metric = app.Metric

			opts := handler()

			for _, v := range opts {
				app.HttpServer.Register(app.Manager, v)
			}

			port := app.Config.(config.ServerConfiguration).GetServerConfig().Port

			log.Printf("Start GCluster http server, successfully, work on port:%d", port)
			log.WithError(http.ListenAndServe(fmt.Sprintf(":%d", port), app.HttpServer.Router))
			return nil
		},
	}
}

func WithCronOption(handler func(mgr manager.GClusterManager) []*gCron.GClusterCronOption) *RunOption {
	return &RunOption{
		Type: RunTypeAsync,
		Process: func(app *GClusterApp) error {

			c := cron.New()
			options := handler(app.Manager)

			for _, option := range options {
				if err := c.AddFunc(option.Spec, option.Handler(app.Manager)); err != nil {
					return err
				} else {
					log.Infof("Start GCluster crontab successful, name=%s, spec=%s, usage=%s", option.Name, option.Spec, option.Usage)
				}
			}
			c.Start()
			return nil
		},
	}
}
