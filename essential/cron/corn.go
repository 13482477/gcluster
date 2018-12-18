package cron

import "mcloud/public.v2/manager"

type MCloudCronHandler func(manager manager.MCloudManager) func()


type MCloudCronOption struct {
	Name    string
	Usage   string
	Spec    string
	Handler MCloudCronHandler
}