package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	return &WebServer{
		Router:        r,
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Router.Method(method, path, handler)
}

// start the server
func (s *WebServer) Start() {
	http.ListenAndServe(s.WebServerPort, s.Router)
}
