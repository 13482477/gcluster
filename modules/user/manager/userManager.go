package manager

import (
	"sync"
	"context"
	userHttp "gcluster/modules/user/http"
	"errors"
	"github.com/zheng-ji/goSnowFlake"
	"fmt"
	"strconv"
)

var userManager *UserManager
var userManagerOnce sync.Once
var userManagerError error

var idWorker *goSnowFlake.IdWorker
var tokenMap = new(sync.Map)

type UserManager struct {
}

func GetUserManager() (*UserManager, error) {
	userManagerOnce.Do(func() {
		userManager = &UserManager{}
		if iw, err := goSnowFlake.NewIdWorker(1); err != nil {
			userManagerError = err
		} else {
			idWorker = iw
		}
	})
	return userManager, userManagerError
}

func (manager *UserManager) Login(ctx context.Context, req *userHttp.LoginRequest) (*userHttp.LoginResponse, error) {
	if req.Username != "lizhiqiang" && req.Password != "password" {
		return nil, errors.New("用户名或密码错误")
	}
	if tokenId, err := idWorker.NextId(); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to generate token, error=%v", err))
	} else {
		token := strconv.FormatInt(tokenId, 10)
		tokenMap.Store(token, "")
		return &userHttp.LoginResponse{
			Token: token,
		}, nil
	}
}
func (manager *UserManager) StartGClusterManager() error {
	return nil
}
