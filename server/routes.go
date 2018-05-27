package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/segmentio/ksuid"
)

func (s *server) setupRouter() {

	s.router = httprouter.New()

	s.router.RedirectFixedPath = true
	s.router.RedirectTrailingSlash = true
	s.router.HandleMethodNotAllowed = true
	s.router.HandleOPTIONS = true

	// GET/HEAD requests for health of the service + dependencies
	s.router.GET("/health", s.health())
	s.router.HEAD("/health", s.health())

	// Handle POST link
	s.router.POST("/v1/link", s.postLink())

	// Handle different types of GET requests
	s.router.GET("/v1/link/:scope", s.link())
	s.router.GET("/v1/link/:scope/:user", s.link())
	s.router.GET("/v1/link/:scope/:user/:title", s.link())
}

func generateRequestID() string {
	id := ksuid.New()
	return id.String()
}

// RequestIDHandler updates request header with unique request id and response header
func RequestIDHandler(handle http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqid := generateRequestID()
		r.Header.Set("X-Shortly-Request-Id", reqid)
		handle.ServeHTTP(w, r)
	}
}
