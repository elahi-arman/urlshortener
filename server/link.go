package server

import (
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

		scope := ps.ByName("scope")
		user := ps.ByName("user")
		title := ps.ByName("title")

		if user != "" {
			if title != "" {
				s.linker.GetLink(scope, user, title)
			} else {
				s.linker.GetLinksByUser(user)
			}
		} else {
			s.linker.GetLinksInScope(scope)
		}
		return
	}

}
