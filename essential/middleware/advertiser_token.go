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

const (
	FakeTypeAgency = 1
	FakeTypeAE     = 2
)

type AdvertiserTokenAllClaims struct {
	AdvertiserTokenClaims
	jwt.StandardClaims
}

type AdvertiserInfo struct {
	AdvertiserId int32  `json:"advertiser_id"`
	AgencyId     int32  `json:"agency_id"`
	Type         int32  `json:"type"`
	CompanyName  string `json:"company_name"`
}

type AccountExecutive struct {
}

type FakeInfo struct {
	Type             int32             `json:"type"`
	Agency           *AdvertiserInfo   `json:"agency"`
	AccountExecutive *AccountExecutive `json:"account_executive"`
}

type AdvertiserTokenClaims struct {
	Advertiser *AdvertiserInfo
	Agency     *AdvertiserInfo
	Fake       *FakeInfo
}

func GenerateAdvertiserToken(data AdvertiserTokenClaims, secretKey string, expireTime int64) (string, error) {
	expireAt := time.Now().Unix() + expireTime
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &AdvertiserTokenAllClaims{
		AdvertiserTokenClaims: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt,
		}})

	return token.SignedString([]byte(secretKey))
}

func DecodeAdvertiserToken(tokenString string, secretKey string) (*AdvertiserTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdvertiserTokenAllClaims{}, func(token *jwt.Token) (interface{}, error) {
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
		return &AdvertiserTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*AdvertiserTokenAllClaims); ok && token.Valid {
		return &AdvertiserTokenClaims{
			Advertiser: claims.Advertiser,
			Agency:     claims.Agency,
			Fake:       claims.Fake,
		}, nil
	} else {
		return nil, TokenInvalid
	}
}

func GetAdvertiserInfoFromContext(ctx context.Context) (AdvertiserTokenClaims, error) {
	var advertiserInfo AdvertiserTokenClaims

	advertiserMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return AdvertiserTokenClaims{}, fmt.Errorf("advertiser token not found")
	} else if _, ok := advertiserMetadata[TokenUserMetadataKey]; !ok {
		return AdvertiserTokenClaims{}, fmt.Errorf("advertiser token not found")
	}

	decodeErr := json.Unmarshal([]byte(strings.Join(advertiserMetadata[TokenUserMetadataKey], "")), &advertiserInfo)
	if decodeErr != nil {
		return AdvertiserTokenClaims{}, fmt.Errorf("advertiser token decode fail")
	}

	return advertiserInfo, nil

}
