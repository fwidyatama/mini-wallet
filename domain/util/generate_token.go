package util

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	CustomerXid string `json:"customer_xid"`
}

func GenerateToken(customerXid string) (string, error) {
	var JWTSecret = viper.Get("JWT_SECRET").(string)

	expirationTime := time.Now().Add(1 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		CustomerXid: customerXid,
	})

	jwtSecretKey := []byte(JWTSecret)

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return tokenString, nil

}
