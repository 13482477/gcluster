package main

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"gcluster/essential/app"
	"gcluster/essential/manager"
	publicHttp "gcluster/essential/http"
	userConfig "gcluster/modules/user/config"
	userManager "gcluster/modules/user/manager"
	userHttp "gcluster/modules/user/http"
)

const (
	AppName    = "user"
	AppUsage   = "user application"
	AppVersion = "1.0.0"
)

func main() {
	userApp := app.GetGClusterApp()
	userApp.Name = AppName
	userApp.Usage = AppUsage
	userApp.Version = AppVersion
	userApp.Config = &userConfig.UserConfig{}

	userApp.Run(
		app.WithLoggerOption(),
		app.WithMetricOption(),
		app.WithOpenTracingOption(),
		app.WithManagerOption(func(db *gorm.DB) (manager.GClusterManager, error) {
			return userManager.GetUserManager()
		}),
		app.WithHttpEndpointOption(func() []*publicHttp.GClusterHttpEndpointOption {
			return []*publicHttp.GClusterHttpEndpointOption{
				{
					Path:       "/user/login",
					Method:     "Login",
					HttpMethod: http.MethodPost,
					CreateReq: func() interface{} {
						return &userHttp.LoginRequest{}
					},
				},
			}
		}),
	)

}
