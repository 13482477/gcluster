package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

type AccountExecutiveTokenAllClaims struct {
	AccountExecutiveTokenClaims
	jwt.StandardClaims
}

type AccountExecutiveTokenClaims struct {
	AccountExecutive *AccountExecutiveInfo
}

type AccountExecutiveInfo struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	Type int32  `json:"type"`
}

func GenerateAccountExecutiveToken(data AccountExecutiveTokenClaims, secretKey string, expireTime int64) (string, error) {
	expireAt := time.Now().Unix() + expireTime
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &AccountExecutiveTokenAllClaims{
		AccountExecutiveTokenClaims: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt,
		}})

	return token.SignedString([]byte(secretKey))
}

func DecodeAccountExecutiveToken(tokenString string, secretKey string) (*AccountExecutiveTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccountExecutiveTokenAllClaims{}, func(token *jwt.Token) (interface{}, error) {
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
		return &AccountExecutiveTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*AccountExecutiveTokenAllClaims); ok && token.Valid {
		return &AccountExecutiveTokenClaims{
			AccountExecutive: claims.AccountExecutive,
		}, nil
	} else {
		return nil, TokenInvalid
	}
}

func GetAccountExecutiveInfoFromContext(ctx context.Context) (AccountExecutiveTokenClaims, error) {
	var accountExecutiveInfo AccountExecutiveTokenClaims

	accountExecutiveInfoMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return AccountExecutiveTokenClaims{}, fmt.Errorf("account executive token not found")
	} else if _, ok := accountExecutiveInfoMetadata[TokenUserMetadataKey]; !ok {
		return AccountExecutiveTokenClaims{}, fmt.Errorf("account executive token not found")
	}

	decodeErr := json.Unmarshal([]byte(strings.Join(accountExecutiveInfoMetadata[TokenUserMetadataKey], "")), &accountExecutiveInfo)
	if decodeErr != nil {
		return AccountExecutiveTokenClaims{}, fmt.Errorf("account executive  token decode fail")
	}

	return accountExecutiveInfo, nil

}
