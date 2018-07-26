package utils

import (
	"github.com/robfig/cron"
	"github.com/zieckey/etcdsync"
	"time"
	log "github.com/sirupsen/logrus"
)

type CronRunner struct {
	Key        string
	RunFunc    func()
	Spec       string
	Ttl        int
	EtcdServer []string
}

func NewCronRunner() *CronRunner {
	return &CronRunner{}
}

func (r *CronRunner) Start() {
	go func() {
		for {
			if r.Ttl < 30 {
				r.Ttl = 30
			}
			m, err := etcdsync.New(r.Key, r.Ttl, r.EtcdServer)
			if m == nil || err != nil {
				time.Sleep(time.Second * time.Duration(r.Ttl/3))
				continue
			}
			err = m.Lock()
			if err != nil {
				time.Sleep(time.Second * time.Duration(r.Ttl/3))
				continue
			} else {
				c := cron.New()
				c.AddFunc(r.Spec, r.RunFunc)
				c.Start()
				for {
					time.Sleep(time.Second * time.Duration(r.Ttl/3))
					m.RefreshLockTTL(time.Duration(r.Ttl))
				}
				m.Unlock()
			}
		}
	}()
}

func (r *CronRunner) StartSimple() {
	log.Debug("-------- start simple -----------")
	c := cron.New()
	c.AddFunc(r.Spec, r.RunFunc)
	c.Start()
}
