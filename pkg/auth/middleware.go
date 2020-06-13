package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/sunnywalden/sync-data/helper"
	"github.com/sunnywalden/sync-data/pkg/errors"
	"github.com/sunnywalden/sync-data/pkg/types"
)


//func middlewareHandler(next http.Handler) http.Handler{
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
//		// 执行handler之前的逻辑
//		next.ServeHTTP(w, r)
//		// 执行完毕handler后的逻辑
//	})
//}

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var res = types.Response{
				Code: http.StatusUnauthorized,
				Msg: errors.ErrNotAuthorized.Error(),
				Data: nil,
			}
			var status = http.StatusUnauthorized

		tokenStr := r.Header.Get("x-authorization-token")
		if tokenStr == "" {
			helper.ResponseWithJson(
				w,
				status,
				res,
			)
		} else {
			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					helper.ResponseWithJson(
						w,
						status,
						res,
						)
					return nil, errors.ErrNotAuthorized
				}
				return []byte("secret"), nil
			})
			if !token.Valid {
				helper.ResponseWithJson(
					w,
					status,
					res,
					)
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}
