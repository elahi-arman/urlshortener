package server

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/elahi-arman/urlshortener/config"
	"github.com/elahi-arman/urlshortener/model"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	linker    model.Linker
	router    *httprouter.Router
	cfg       *config.AppConfig
	accessLog *zap.Logger
	appLog    *zap.SugaredLogger
	home      string
}

//StartServer creates a server based on the config passed in
func StartServer(cfgPath string, home string) error {

	var err error
	var s = new(server)

	s.home = home
	s.cfg, err = config.NewConfig(cfgPath)
	if err != nil {
		return err
	}

	s.setupLoggers()

	s.appLog.Debugw("", "msg", "Starting router setup")
	s.setupRouter()
	s.appLog.Debugw("", "msg", "Finished router setup")

	s.appLog.Debugw("", "msg", "Starting Linker setup")
	linker, err := model.NewRedisLinker(s.cfg.Redis)
	if err != nil {
		return err
	}
	s.linker = linker
	s.appLog.Debugw("", "msg", "Finished Linker setup")

	s.appLog.Infow("", "msg", fmt.Sprintf("Configuration all good, serving at %s", s.cfg.Server.Address))
	http.ListenAndServe(s.cfg.Server.Address, s.router)

	return nil
}

//setupLoggers after this point, s.appLog is available for logging
func (s *server) setupLoggers() error {

	var (
		err    error
		appLog *zap.Logger
		cfg    zap.Config
	)

	for i, path := range s.cfg.Log.OutputPaths {
		s.cfg.Log.OutputPaths[i] = s.home + "/" + path
	}

	appLog, err = s.cfg.Log.Build()
	if err != nil {
		return err
	}
	defer appLog.Sync()
	s.appLog = appLog.Sugar()
	s.appLog.Debugw("", "msg", "Finished setting up app log")

	s.appLog.Debugw("", "msg", "Starting access log setup")
	cfg = zap.NewProductionConfig()
	cfg.OutputPaths = []string{s.home + "/" + s.cfg.Server.AccessLog}
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.Encoding = "console"

	s.accessLog, err = cfg.Build()
	if err != nil {
		return err
	}
	defer s.accessLog.Sync()
	s.appLog.Debugw("", "msg", "Finished access log setup")

	return nil

}
