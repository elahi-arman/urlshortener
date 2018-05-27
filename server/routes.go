package server

import (
	"github.com/julienschmidt/httprouter"
)

func (s *server) setupRouter() {

	s.router = httprouter.New()

	s.router.RedirectFixedPath = true
	s.router.RedirectTrailingSlash = true
	s.router.HandleMethodNotAllowed = true
	s.router.HandleOPTIONS = true

	// GET/HEAD requests for health of the service + dependencies
	s.router.GET("/health", s.requestLogger(s.health()))
	s.router.HEAD("/health", s.requestLogger(s.health()))

	// Handle POST link
	s.router.POST("/v1/link", s.requestLogger(s.postLink()))

	// Handle different types of GET requests
	s.router.GET("/v1/link/:link", s.requestLogger(s.link()))
	s.router.GET("/v1/link/:link/:user", s.requestLogger(s.link()))
	s.router.GET("/v1/link/:link/:user/:scope", s.requestLogger(s.link()))
}
