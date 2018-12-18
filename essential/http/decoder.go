package http

import (
	"encoding/json"
	transport "github.com/go-kit/kit/transport/http"
	"context"
	"net/http"
	"github.com/kataras/iris/core/errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func MakeRequestDecoder(option *MCloudHttpEndpointOption) transport.DecodeRequestFunc {
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
