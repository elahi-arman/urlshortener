package server

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	id     string
	start  time.Time
	end    time.Time
	path   string
	user   string
	method string
	status int
	code   string
	logMsg string
}

//WriteResponse adds log context onto each response, passing status to WriteHeader
func (lw *loggingResponseWriter) WriteResponse(status int, logCode string, logMsg string) {
	lw.status = status
	lw.code = logCode
	lw.logMsg = logMsg
	lw.ResponseWriter.WriteHeader(status)
}

func (s *server) requestLogger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// in case handler doesn't call WriteHeader
		lw := loggingResponseWriter{
			w,
			ksuid.New().String(),
			time.Now(),
			time.Now(),
			"",
			"",
			"",
			200,
			"",
			"",
		}

		h(&lw, r, ps)
		lw.path = r.URL.Path
		lw.user = ""
		lw.method = r.Method

		lw.end = time.Now()

		if lw.status >= 400 {
			s.logRequest(lw, true)
		} else {
			s.logRequest(lw, false)
		}

		return
	}
}

func (s *server) logRequest(lw loggingResponseWriter, isError bool) {
	if isError {
		s.accessLog.Error("",
			zap.String("id", lw.id),
			zap.String("path", lw.path),
			zap.Time("start", lw.start),
			zap.Time("end", lw.end),
			zap.String("user", lw.user),
			zap.String("method", lw.method),
			zap.Int("status", lw.status),
			zap.String("log_code", lw.code),
			zap.String("log_msg", lw.logMsg),
		)
	} else {
		s.accessLog.Info("",
			zap.String("id", lw.id),
			zap.String("path", lw.path),
			zap.Time("start", lw.start),
			zap.Time("end", lw.end),
			zap.String("user", lw.user),
			zap.String("method", lw.method),
			zap.Int("status", lw.status),
			zap.String("log_code", lw.code),
			zap.String("log_msg", lw.logMsg),
		)
	}
}
