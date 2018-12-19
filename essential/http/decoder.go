package http

import (
	"fmt"
	"context"
	"net/http"
	"encoding/json"
	"github.com/kataras/iris/core/errors"
	transport "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
)

func MakeRequestDecoder(option *GClusterHttpEndpointOption) transport.DecodeRequestFunc {
	if option.ReqDecoder != nil {
		return option.ReqDecoder
	} else {
		return func(ctx context.Context, r *http.Request) (interface{}, error) {
			req := option.CreateReq()
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				err = errors.New(fmt.Sprintf("http request decode failed, error=%v", err))
				log.WithError(err).Error()
				return nil, err
			} else {
				return req, nil
			}
		}
	}
}
