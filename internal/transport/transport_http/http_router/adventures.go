package http_router

import "net/http"

func setupAdventuresRouter() http.Handler {
	mux := http.NewServeMux()

	return mux
}
