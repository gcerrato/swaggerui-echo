package swaggerui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:generate go run generate.go

//go:embed embed
var swagfs embed.FS

func byteHandler(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Write(b)
	}
}

// Handler returns a handler that will serve a self-hosted Swagger UI with your spec embedded
func Handler(spec []byte) http.Handler {
	// render the index template with the proper spec name inserted
	static, _ := fs.Sub(swagfs, "embed")
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger_spec", byteHandler(spec))
	mux.Handle("/", http.FileServer(http.FS(static)))
	return mux
}

// EchoHandler returns an echo.HandlerFunc that will serve a self-hosted Swagger UI with your spec embedded
func EchoHandler(strip string, spec []byte) echo.HandlerFunc {
	handler := Handler(spec)
	return echo.WrapHandler(http.StripPrefix(strip, handler))
}
