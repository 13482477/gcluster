package http

import (
	"encoding/json"
	"context"
	"net/http"
	transport "github.com/go-kit/kit/transport/http"
)

func MakeResponseEncoder(option *GClusterHttpEndpointOption) transport.EncodeResponseFunc {
	if option.RespEncoder != nil {
		return option.RespEncoder
	} else {
		return EncodeJSONResponseWithBaseResponse
	}
}

func EncodeJSONResponseWithBaseResponse(_ context.Context, w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if header, ok := data.(transport.Headerer); ok {
		for k, values := range header.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := data.(transport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	resp := &BaseResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}

	return json.NewEncoder(w).Encode(resp)
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(&BaseResponse{
		Code:    500,
		Message: err.Error(),
	})
}
