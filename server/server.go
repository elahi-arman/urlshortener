package server

import (
	"net/http"

	"github.com/elahi-arman/urlshortener/config"
	"github.com/elahi-arman/urlshortener/model"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	linker model.Linker
	router *httprouter.Router
	cfg    *config.AppConfig
}

//StartServer creates a server based on the config passed in
func StartServer(cfgPath string) error {

	var err error
	var s = new(server)
	s.setupRouter()
	s.cfg, err = config.NewConfig(cfgPath)

	if err != nil {
		return err
	}

	linker, err := model.NewRedisLinker(s.cfg.Redis)
	if err != nil {
		return err
	}

	s.linker = linker
	http.ListenAndServe(s.cfg.Server.Address, s.router)

	return nil
}
