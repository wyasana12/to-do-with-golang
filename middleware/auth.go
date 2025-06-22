package middleware

import (
	"context"
	"net/http"
	"strings"
	"to-do-list-go/helper"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			helper.Response(w, 401, "unauthorized: Authorization header missing", nil)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			helper.Response(w, 401, "unauthorized: Invalid Authorization header format. Expected 'Bearer <token>'", nil)
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := helper.ValidateToken(accessToken) // Sekarang hanya token mentah yang masuk ke ValidateToken
		if err != nil {
			helper.Response(w, 401, err.Error(), nil)
			return
		}

		ctx := context.WithValue(r.Context(), "User Info", user)
		next.ServeHTTP(w, r.WithContext(ctx))
		// accessToken := r.Header.Get("Authorization")

		// if accessToken == "" {
		// 	helper.Response(w, 401, "unauthorized", nil)
		// 	return
		// }

		// user, err := helper.ValidateToken(accessToken)
		// if err != nil {
		// 	helper.Response(w, 401, err.Error(), nil)
		// 	return
		// }

		// ctx := context.WithValue(r.Context(), "User Info", user)
		// next.ServeHTTP(w, r.WithContext(ctx))
	})
}
