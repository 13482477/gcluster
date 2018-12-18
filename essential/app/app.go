package app

import (
	"gcluster/essential/rpc"
	"gcluster/essential/registry"
	"sync"
	"github.com/urfave/cli"
	"os"
	log "github.com/sirupsen/logrus"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/hashicorp/consul/api"
	"github.com/go-kit/kit/sd/consul"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/opentracing/opentracing-go"
	mcloudHttp "gcluster/essential/http"
	"net/http"
	"fmt"
	"gcluster/essential/manager"
	"gcluster/essential/config"
	"gcluster/essential/metric"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron"
	ecsCron "gcluster/essential/cron"
	"github.com/rifflock/lfshook"
)

var mcloudApp *McloudApp
var mcloudAppOnce sync.Once

type RunType int

const (
	RunTypeSync  RunType = 1
	RunTypeAsync RunType = 2
)

type RunOption struct {
	Type    RunType
	Process func(mcApp *McloudApp) error
}

type McloudApp struct {
	Name       string
	Usage      string
	Version    string
	Config     config.McloudConfig
	Metric     *metric.MCloudMetric
	Manager    manager.MCloudManager
	Client     consul.Client
	Registry   *registry.McloudServiceRegistry
	Tracer     opentracing.Tracer
	RpcManager *rpc.MCloudRpcManager
	HttpServer *mcloudHttp.MCloudHttpServer
	RunOptions []*RunOption
}

func GetMcloudApp() *McloudApp {
	mcloudAppOnce.Do(func() {
		mcloudApp = &McloudApp{}
	})
	return mcloudApp
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

func (mcApp *McloudApp) Run(runOptions ...*RunOption) error {
	mcApp.RunOptions = runOptions

	app := &cli.App{
		Name:    mcApp.Name,
		Usage:   mcApp.Usage,
		Version: mcApp.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config, c",
				Value: "config",
				Usage: "config file path",
			},
		},
		Action: func(ctx *cli.Context) error {

			filePath := fmt.Sprintf("./%s.log", mcApp.Name)
			hook := lfshook.NewHook(filePath, nil)
			log.AddHook(hook)

			log.SetLevel(log.DebugLevel)

			printLogo()

			log.Infof("========================================================================================")
			log.Infof("======================================System start======================================")
			log.Infof("========================================================================================")
			log.WithField("SystemName", mcApp.Name).Info()
			log.WithField("Version", mcApp.Version).Info()

			configLoader := &config.MCloudConfigLoader{
				Name:     ctx.String("config"),
				FilePath: ".",
				Config:   mcApp.Config,
			}

			if err := configLoader.Load(); err != nil {
				log.Panicf("load config file failed, error=%v", err)
			}

			wg := sync.WaitGroup{}
			wg.Add(1)

			for _, runOption := range mcApp.RunOptions {
				if runOption.Type == RunTypeSync {
					if err := runOption.Process(mcApp); err != nil {
						log.Panic(err)
					}
				} else {
					localRunOption := runOption
					wg.Add(1)
					go func() {
						if err := localRunOption.Process(mcApp); err != nil {
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
		Process: func(app *McloudApp) error {
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
		Process: func(app *McloudApp) error {

			fieldKeys := []string{"method"}

			app.Metric = &metric.MCloudMetric{
				RequestCount: kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
					Namespace: "mcloud",
					Subsystem: app.Name,
					Name:      "request_count",
					Help:      "Number of requests received.",
				}, fieldKeys),

				RequestLatency: kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
					Namespace: "mcloud",
					Subsystem: app.Name,
					Name:      "request_latency_microseconds",
					Help:      "Total duration of requests in microseconds.",
				}, fieldKeys),
			}

			log.Printf("Start MCloud metric server successfully!")
			return nil
		},
	}
}

func WithRegistryOption() *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *McloudApp) error {

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

			mcloudServiceRegistry := registry.McloudServiceRegistry{
				Client:     consul.NewClient(consulClient),
				ServerName: app.Name,
				Config:     app.Config,
			}

			app.Client = mcloudServiceRegistry.Client

			mcloudServiceRegistry.Register()
			log.Printf("Start MCloud server mcloudServiceRegistry successfully!")
			return nil
		},
	}
}

func WithOpenTracingOption() *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *McloudApp) error {
			tracer := opentracing.GlobalTracer()
			app.Tracer = tracer
			log.Printf("Start MCloud trace server successfully!")
			return nil
		},
	}
}

func WithManagerOption(handler func(db *gorm.DB) (manager.MCloudManager, error)) *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *McloudApp) error {
			dbConfig := app.Config.(config.DatabaseConfiguration).GetDataBaseConfiguration()
			db, err := gorm.Open("mysql", dbConfig.Address)
			if err != nil {
				log.Panicf("Fatal error mysql connection failed: %v", err)
			}
			db.SingularTable(true)
			db.LogMode(dbConfig.LogMode)
			db.DB().SetMaxIdleConns(dbConfig.MaxIdle)
			db.DB().SetMaxOpenConns(dbConfig.MaxConns)
			db.DB().SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetime) * time.Second)

			if mCloudManager, err := handler(db); err != nil {
				return err
			} else {
				app.Manager = mCloudManager
				if err := mCloudManager.StartMcloudManager(); err != nil {
					log.Panicf("Start MCloud manager failed, error=%v", err)
				} else {
					log.Infof("Start MCloud manager successfully!")
				}
				return nil
			}
		},
	}
}

func WithRpcOption(handler func() []*rpc.MCloudRpcOption) *RunOption {
	return &RunOption{
		Type: RunTypeSync,
		Process: func(app *McloudApp) error {
			rpcManager := rpc.GetMCloudRpcManager(app.Client, app.Tracer)
			app.RpcManager = rpcManager

			options := handler()

			for _, v := range options {
				rpcManager.Subscript(v)
			}

			log.Printf("Start MCloud rcp server successfully!")
			return nil
		},
	}
}

func WithHttpEndpointOption(handler func() []*mcloudHttp.MCloudHttpEndpointOption) *RunOption {
	return &RunOption{
		Type: RunTypeAsync,
		Process: func(app *McloudApp) error {
			app.HttpServer = mcloudHttp.GetHttpServer()
			app.HttpServer.Tracer = app.Tracer
			app.HttpServer.Metric = app.Metric

			opts := handler()

			for _, v := range opts {
				app.HttpServer.Register(app.Manager, v)
			}

			port := app.Config.(config.ServerConfiguration).GetServerConfig().Port

			log.Printf("Start MCloud http server, successfully, work on port:%d", port)
			log.WithError(http.ListenAndServe(fmt.Sprintf(":%d", port), app.HttpServer.Router))
			return nil
		},
	}
}

func WithCronOption(handler func(mgr manager.MCloudManager) []*ecsCron.MCloudCronOption) *RunOption {
	return &RunOption{
		Type: RunTypeAsync,
		Process: func(app *McloudApp) error {

			c := cron.New()
			options := handler(app.Manager)

			for _, option := range options {
				if err := c.AddFunc(option.Spec, option.Handler(app.Manager)); err != nil {
					return err
				} else {
					log.Infof("Start MCloud crontab successful, name=%s, spec=%s, usage=%s", option.Name, option.Spec, option.Usage)
				}
			}
			c.Start()
			return nil
		},
	}
}
