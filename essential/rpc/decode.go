package rpc

import (
	"fmt"
	"context"
	"net/http"
	"encoding/json"
	"github.com/kataras/iris/core/errors"
)

func MakeDecodeJsonResponse(resp interface{}) func(_ context.Context, httpResp *http.Response) (interface{}, error) {
	return func(_ context.Context, httpResp *http.Response) (interface{}, error) {
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			err = errors.New(fmt.Sprintf("decode rpc response failed, error=%v", err))
			return nil, err
		}
		return resp, nil
	}
}
