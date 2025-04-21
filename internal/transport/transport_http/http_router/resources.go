package http_router

import (
	"net/http"

	"github.com/yazmeyaa/ton_game/internal/transport/transport_http/http_controllers"
)

func setupResourcesRouter(rc http_controllers.ResourcesController) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /all", http.HandlerFunc(rc.All))
	mux.Handle("GET /items", http.HandlerFunc(rc.Items))
	mux.Handle("GET /crafts", http.HandlerFunc(rc.Crafts))
	mux.Handle("GET /adventures", http.HandlerFunc(rc.Adventures))

	return mux
}
