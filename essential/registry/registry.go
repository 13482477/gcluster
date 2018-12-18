package registry

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"gcluster/essential/config"
	"fmt"
	applog "gcluster/essential/log"
	"net"
)

type McloudServiceRegistry struct {
	ServerName string
	Config     config.McloudConfig
	Client     consul.Client
}

func (registry *McloudServiceRegistry) Register() {
	name := registry.ServerName
	localAddress := registry.Config.(config.ServerConfiguration).GetServerConfig().Address
	if localAddress == "" {
		localAddress = resolveLocalIp()
	}

	port := registry.Config.(config.ServerConfiguration).GetServerConfig().Port

	asg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s:%s:%d", name, localAddress, port),
		Name:    name,
		Port:    port,
		Address: localAddress,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "5m",
			Interval:                       "10s",
			TCP:                            fmt.Sprintf("%s:%d", localAddress, port),
		},
	}

	register := consul.NewRegistrar(registry.Client, asg, applog.GetConsulLogger())
	register.Register()
}

func resolveLocalIp() string {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}

		}
	}

	return ""
}
