package server

import (
	"io/fs"
	v1 "kratos-demo/api/helloworld/v1"
	"kratos-demo/internal/conf"
	"kratos-demo/internal/service"
	"kratos-demo/third_party"
	"mime"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func getOpenAPIHandler() stdhttp.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
	if err != nil {
		log.Fatalf("couldn't create sub filesystem: %v", err)
	}
	return stdhttp.FileServer(stdhttp.FS(subFS))
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			// tracing.Server(),
			// logging.Server(logger),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	srv.HandlePrefix("/swagger-ui", getOpenAPIHandler())
	srv.HandlePrefix("/redoc", getOpenAPIHandler())
	return srv
}
