package applog

import (
	"sync"
	log "github.com/sirupsen/logrus"
	goKitLog "github.com/go-kit/kit/log"
)

var consulLogger goKitLog.Logger
var consulLoggerOnce sync.Once

func GetConsulLogger() goKitLog.Logger {
	consulLoggerOnce.Do(func() {
		consulLogger = goKitLog.NewLogfmtLogger(log.StandardLogger().WriterLevel(log.DebugLevel))
		consulLogger = goKitLog.With(consulLogger, "usage", "consul")
		consulLogger = goKitLog.With(consulLogger, "caller", goKitLog.DefaultCaller)
	})

	return consulLogger
}

var openTracingLogger goKitLog.Logger
var openTracingLoggerOnce sync.Once

func GetOpenTracingLogger() goKitLog.Logger {
	openTracingLoggerOnce.Do(func() {
		openTracingLogger = goKitLog.NewLogfmtLogger(log.StandardLogger().WriterLevel(log.DebugLevel))
		openTracingLogger = goKitLog.With(openTracingLogger, "usage", "openTracing")
		openTracingLogger = goKitLog.With(openTracingLogger, "caller", goKitLog.DefaultCaller)
	})
	return openTracingLogger
}

var endpointLogger goKitLog.Logger
var endpointLoggerOnce sync.Once

func GetEndpointLogger() goKitLog.Logger {
	endpointLoggerOnce.Do(func() {
		endpointLogger = goKitLog.NewLogfmtLogger(log.StandardLogger().WriterLevel(log.DebugLevel))
		endpointLogger = goKitLog.With(consulLogger, "usage", "endpoint")
		endpointLogger = goKitLog.With(consulLogger, "caller", goKitLog.DefaultCaller)
	})
	return endpointLogger
}
