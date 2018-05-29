package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *server) linkByUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		user := ps.ByName("user")
		link := ps.ByName("link")
		w.Write([]byte(fmt.Sprintf("retrieving link (%s) by user (%s)", link, user)))
		// s.linker.GetLink("*", user, link)

		w.WriteHeader(200)
		return
	}

}

func (s *server) allLinksByUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := ps.ByName("user")
		w.Write([]byte(fmt.Sprintf("retrieving all links by user (%s)", user)))

		// s.linker.GetLinksByUser(user)

		w.WriteHeader(200)
		return
	}
}
