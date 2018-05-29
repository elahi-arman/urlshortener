package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *server) postLink() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func (s *server) link() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		link := ps.ByName("link")
		w.Write([]byte(fmt.Sprintf("retrieving closest link (%s)", link)))

		w.WriteHeader(200)
		return
	}

}
