package main

import (
	"embed"
	"net/http"

	"github.com/easyp-tech/grpc-cource-2/third_party/swagger"
)

func serveSwagger(mux *http.ServeMux, swaggerSpecs embed.FS) {
	swaggerStaticsHandler := http.StripPrefix("/swagger", http.FileServer(http.FS(swagger.Content)))
	mux.Handle("GET /swagger/", swaggerStaticsHandler)

	swaggerSpecsHandler := http.StripPrefix("/swagger/specs", http.FileServer(http.FS(swaggerSpecs)))
	mux.Handle("GET /swagger/specs/", swaggerSpecsHandler)
}
