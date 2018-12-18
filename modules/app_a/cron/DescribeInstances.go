package cron

import (
	log "github.com/sirupsen/logrus"
	"mcloud/public.v2/manager"
	ecsManager "mcloud/ecs.v2/manager"
)

func DescribeInstancesHandler(manager manager.MCloudManager) func() {
	return func() {
		log.Infof("======start to sync ecs instance======")
		if err := manager.(*ecsManager.EcsManager).DescribeInstances(); err != nil {
			log.Errorf("sync instance status failed, error=%v", err)
		}
		log.Infof("======end to sync ecs instance======")
	}
}
