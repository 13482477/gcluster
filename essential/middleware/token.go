package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/armon/go-radix"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"google.golang.org/grpc/metadata"
	"net/http"
	"poseidon/essential/app"
	"strings"
	"time"
)

type DefaultResponse struct {
	Code         int
	ErrorMessage string
	Data         interface{}
}

type UserTokenAllClaims struct {
	UserTokenClaims
	jwt.StandardClaims
}

const (
	TokenUserMetadataKey = "custom_user_info"
)

type UserTokenClaims struct {
	UserId int32 `json:"user_id"`
}

var (
	TokenExpired error = errors.New("Token is expired")
	TokenInvalid error = errors.New("Token is invalid")
)

func GenerateToken(data UserTokenClaims, secretKey string, expireTime int64) (string, error) {
	expireAt := time.Now().Unix() + expireTime
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &UserTokenAllClaims{
		UserTokenClaims: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt,
		}})

	return token.SignedString([]byte(secretKey))
}

func TokenMiddleware(conf app.TokenConfig) echo.MiddlewareFunc {
	r := radix.New()
	for _, path := range conf.HcCheckPath {
		r.Insert(path, "HC")
	}
	for _, path := range conf.HcIgnorePath {
		r.Insert(path, "HC_IGNORE")
	}
	for _, path := range conf.JwCheckPath {
		r.Insert(path, "JW")
	}
	for _, path := range conf.JwIgnorePath {
		r.Insert(path, "JW_IGNORE")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.LastIndex(c.Request().RequestURI, "/inner/") != -1 {
				return c.JSON(http.StatusForbidden, DefaultResponse{
					Code:         403,
					ErrorMessage: "forbidden",
				})
			}
			if _, val, ok := r.LongestPrefix(c.Request().RequestURI); ok {
				tokenString := ""
				auth := c.Request().Header.Get("Authorization")
				scheme := "Bearer"
				l := len(scheme)
				if len(auth) > l+1 && auth[:l] == scheme {
					tokenString = auth[l+1:]
				}
				if val.(string) == "HC" {
					tokenData, err := DecodeToken(tokenString, conf.Secret)
					if err != nil {
						return c.JSON(200, DefaultResponse{
							Code:         401,
							ErrorMessage: err.Error(),
						})
					}
					c.Set("user", tokenData)
				} else if val.(string) == "JW" {
					tokenData, err := DecodeAdvertiserToken(tokenString, conf.Secret)
					if err != nil {
						return c.JSON(200, DefaultResponse{
							Code:         401,
							ErrorMessage: err.Error(),
						})
					}
					c.Set("user", tokenData)
				}
			}
			return next(c)
		}
	}
}

func DecodeToken(tokenString string, secretKey string) (*UserTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserTokenAllClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else {
				return nil, TokenInvalid
			}
		}
		return &UserTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*UserTokenAllClaims); ok && token.Valid {
		return &UserTokenClaims{
			UserId: claims.UserId,
		}, nil
	} else {
		return nil, TokenInvalid
	}
}

func GetUserInfoFromContext(ctx context.Context) (UserTokenClaims, error) {
	var userInfo UserTokenClaims

	userInfoMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return UserTokenClaims{}, fmt.Errorf("user token not found")
	} else if _, ok := userInfoMetadata[TokenUserMetadataKey]; !ok {
		return UserTokenClaims{}, fmt.Errorf("user token not found")
	}

	decodeErr := json.Unmarshal([]byte(strings.Join(userInfoMetadata[TokenUserMetadataKey], "")), &userInfo)
	if decodeErr != nil {
		return UserTokenClaims{}, fmt.Errorf("user token decode fail")
	}

	return userInfo, nil

}
