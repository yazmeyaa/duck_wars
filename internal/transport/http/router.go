package http_server

import "net/http"

func New() *http.Server {
	router := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &server
}
