package handler

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"poseidon/essential/endpoint"

	"encoding/json"
	"google.golang.org/grpc/metadata"
	"net/http"
	"regexp"
	"strings"
)

type HttpApiHandler struct {
	Endpoint *endpoint.EndPoint
}

type DefaultResponse struct {
	Code         int         `json:"code"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

type RedirectResponse struct {
	Code        int    `json:"code"`
	RedirectUrl string `json:"redirect_url"`
}

func (handler *HttpApiHandler) HandleMethod(ctx echo.Context) error {
	log.Debugf("Handle %s %s", ctx.Request().Method, ctx.Request().URL.Path)
	span, requestContext := opentracing.StartSpanFromContext(ctx.Request().Context(), ctx.Request().URL.Path)
	defer span.Finish()

	loc := endpoint.ParseLocator(ctx.Request().URL.Path)
	method, ok := handler.Endpoint.GetMethodDetail(loc)
	if !ok {
		return ctx.JSON(http.StatusNotFound, DefaultResponse{Code: 404, ErrorMessage: fmt.Sprintf("method not found [%s]", loc)})
	}

	req := method.CreateRequest()

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusOK, DefaultResponse{Code: 500, ErrorMessage: fmt.Sprintf("cannot parse request body, error:[%s]", err)})
	}
	tokenUser := ctx.Get("user")

	if tokenUser != nil {
		tokenUserJson, _ := json.Marshal(tokenUser)
		requestContext = metadata.AppendToOutgoingContext(requestContext, "custom_user_info", string(tokenUserJson))
	}

	resp, err := handler.Endpoint.Call(requestContext, loc, req)
	if err != nil {
		re := regexp.MustCompile(`(?U)\{\{([^\{]*[^\}])\}\}`)
		match := re.FindAllStringSubmatch(err.Error(), -1)
		var err string
		if len(match) > 0 {
			err = match[len(match)-1][1]
		}
		return ctx.JSON(http.StatusOK, DefaultResponse{
			Code:         500,
			ErrorMessage: err,
		})
	}

	marshaler := jsonpb.Marshaler{EmitDefaults: true}
	responseBody, err := marshaler.MarshalToString(resp.(proto.Message))
	if err != nil {
		return ctx.JSON(http.StatusOK, DefaultResponse{
			Code:         500,
			ErrorMessage: fmt.Sprintf("encode error: %v", err),
		})
	}

	return ctx.Blob(http.StatusOK, echo.MIMEApplicationJSONCharsetUTF8, []byte(responseBody))
}

func (handler *HttpApiHandler) HandleRedirect(ctx echo.Context) error {
	log.Debugf("Handle %s %s", ctx.Request().Method, ctx.Request().URL.Path)
	span, requestContext := opentracing.StartSpanFromContext(ctx.Request().Context(), ctx.Request().URL.Path)
	defer span.Finish()

	loc := endpoint.ParseLocator(ctx.Request().URL.Path)
	method, ok := handler.Endpoint.GetMethodDetail(loc)
	if !ok {
		return ctx.JSON(http.StatusNotFound, DefaultResponse{Code: 404, ErrorMessage: fmt.Sprintf("method not found [%s]", loc)})
	}

	req := method.CreateRequest()

	if err := ctx.Bind(req); err != nil {
		return err
	}
	resp, err := handler.Endpoint.Call(requestContext, loc, req)
	if err != nil {
		errSlice := strings.Split(err.Error(), ":")

		return ctx.JSON(http.StatusOK, DefaultResponse{
			Code:         500,
			ErrorMessage: errSlice[len(errSlice)-1],
		})
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		ctx.JSON(http.StatusOK, DefaultResponse{
			Code:         500,
			ErrorMessage: fmt.Sprintf("encode error: %v", err),
		})
	}

	redirectData := RedirectResponse{}

	if err := json.Unmarshal(jsonResp, &redirectData); err != nil {
		return ctx.JSON(http.StatusOK, DefaultResponse{
			Code:         500,
			ErrorMessage: fmt.Sprintf("redirect info error: %v", err),
		})
	}

	return ctx.Redirect(redirectData.Code, redirectData.RedirectUrl)
}

func (handler *HttpApiHandler) HandleDefault(context echo.Context) error {
	return context.JSON(http.StatusOK, DefaultResponse{
		Code:         0,
		ErrorMessage: "success",
	})
}
