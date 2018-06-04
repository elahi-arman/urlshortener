package server

import (
	"net/http"

	"github.com/elahi-arman/urlshortener/model"
	"github.com/segmentio/ksuid"

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
		s := model.ServerContext{
			s.linker,
			s.appLog,
			ksuid.New().String(),
		}
		link := ps.ByName("link")
		model.SearchForLink(s, "aelahi", link)
		w.WriteHeader(200)
		return
	}

}
