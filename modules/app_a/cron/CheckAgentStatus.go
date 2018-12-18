package cron

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	log "github.com/sirupsen/logrus"
	"mcloud/public.v2/manager"
	"mcloud/ecs.v2/model"
	ecsManager "mcloud/ecs.v2/manager"
)

func CheckAgentStatusHandler(manager manager.MCloudManager) func() {
	return func() {
		db := manager.(*ecsManager.EcsManager).Db

		ecsAgents := make([]*model.EcsAgent, 0)

		if err := db.Find(&ecsAgents).Error; err != nil {
			log.Errorf("Failed to get ecs agents from db, error=%v", err)
			return
		}

		for _, agent := range ecsAgents {
			log.WithFields(log.Fields{"agent_url": agent.CheckUrl, "__sign": agent.CheckUrlSign,}).Debug("Start to check agent status.")
			status, err := CheckAgentStatus(agent.CheckUrl, agent.CheckUrlSign)
			if err != nil {
				log.Errorf("Get agent status failed, error=%v", err)
				continue
			}

			ecsAgent := &model.EcsAgent{}
			if err := db.Where(&agent).First(&ecsAgent).Error; err != nil {
				log.Errorf("Get agent from db failed, error=%v", err)
				continue
			}

			ecsAgent.AgentStatus = status
			ecsAgent.UpdatedTime = time.Now()
			if err := db.Save(ecsAgent).Error; err != nil {
				log.Errorf("Update agent status failed, error=%v", err)
			}
		}

	}
}

func CheckAgentStatus(url, sign string) (status int, err error) {
	if url == "" {
		return 0, errors.New("url is empty")
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("__sign", sign)
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	ret := CheckAgent{}
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return 0, err
	}
	if ret.Code == 200 {
		return ret.Data.Status, nil
	} else {
		return 0, errors.New("check status failed")
	}
}

type AgentInfo struct {
	ApplicationName string `json:"applicationName"`
	Ip              string `json:"ip"`
	Port            int    `json:"port"`
	RemoteAddress   string `json:"remoteAddress"`
	Host            string `json:"host"`
	Status          int    `json:"status"`
	Description     string `json:"description"`
}

type CheckAgent struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data AgentInfo `json:"data"`
}
