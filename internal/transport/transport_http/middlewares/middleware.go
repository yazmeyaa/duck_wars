package middlewares

import (
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

func ApplyMiddleware(to http.Handler, middlewares ...Middleware) http.Handler {
	for _, mw := range middlewares {
		to = mw(to)
	}

	return to
}
