package cron

import "gcluster/essential/manager"

type GClusterCronHandler func(manager manager.GClusterManager) func()


type GClusterCronOption struct {
	Name    string
	Usage   string
	Spec    string
	Handler GClusterCronHandler
}