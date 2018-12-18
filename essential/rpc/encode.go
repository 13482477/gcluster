package rpc

import (
	"encoding/json"
	"io/ioutil"
	"context"
	"net/http"
	"bytes"
)

func EncodeJSONRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}