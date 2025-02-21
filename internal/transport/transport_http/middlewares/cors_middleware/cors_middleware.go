package cors_middleware

import (
	"net/http"

	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/middlewares"
)

// var allowedOrigins = []string{"http://localhost:3000", "http://localhost:5173", "http://192.168.1.35:3000"}

func CorsMiddleware() middlewares.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
