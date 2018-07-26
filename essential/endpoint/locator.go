package endpoint

import (
	"fmt"
	"strings"
)

const (
	LBStrategyRR = "rr"
	LBStrategyHash = "hash"
)

type Locator struct {
	Cluster string
	Version string
	Service string
	Method  string
	Lb      string			//lb 策略
}

func ParseLocator(loc string) Locator {
	loc = strings.Trim(loc, "/")
	list := strings.Split(loc, "/")
	if len(list) >= 4 {
		return Locator{
			Cluster: list[0],
			Version: list[1],
			Service: list[2],
			Method:  strings.Join(list[3:], "/"),
			Lb: LBStrategyRR,
		}
	}
	return Locator{}
}

func (loc Locator) String() string {
	return fmt.Sprintf("/%s/%s/%s/%s", loc.Cluster, loc.Version, loc.Service, loc.Method)
}
