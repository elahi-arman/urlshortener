package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *server) health() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
