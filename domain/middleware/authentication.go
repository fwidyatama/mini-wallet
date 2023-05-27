package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"mini-wallet/domain/util"
	"net/http"
	"strings"
)

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.Background()

		authHeader := r.Header.Get("Authorization")
		splitHeader := strings.Split(authHeader, " ")

		if strings.ToLower(splitHeader[0]) != "bearer" || authHeader == "" {
			util.FailedResponseWriter(w, "unauthorized", http.StatusUnauthorized)
			return
		} else {
			tokenString := splitHeader[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				var JWTSecret = viper.Get("JWT_SECRET").(string)
				jwtSecretKey := []byte(JWTSecret)

				return jwtSecretKey, nil
			})

			if err != nil {
				util.FailedResponseWriter(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx = context.WithValue(ctx, "ownerId", claims["customer_xid"])
				r = r.WithContext(ctx)
				next(w, r, ps)
			} else {
				util.FailedResponseWriter(w, "unauthorized", http.StatusUnauthorized)
				return
			}

		}

	}
}
