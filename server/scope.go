package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *server) linkInScope() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		scope := ps.ByName("scope")
		link := ps.ByName("link")
		w.Write([]byte(fmt.Sprintf("retrieving link(%s) in scope (%s)", link, scope)))
		// s.linker.GetLink(scope, "*", link)

		w.WriteHeader(200)
		return
	}

}

func (s *server) allLinksInScope() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		scope := ps.ByName("scope")
		w.Write([]byte(fmt.Sprintf("retrieving all links in scope (%s)", scope)))
		// s.linker.GetLinksInScope(scope)

		w.WriteHeader(200)
		return
	}
}
