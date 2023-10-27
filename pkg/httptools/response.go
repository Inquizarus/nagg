package httptools

import (
	"net/http"

	"github.com/inquizarus/nagg/pkg/logging"
)

type ResponseWriterWrapper interface {
	http.ResponseWriter
	Status() int
}

func NewResponseWriterWrapper(w http.ResponseWriter, statusCode int, log logging.Logger) responseWriterWrapper {
	return responseWriterWrapper{
		w,
		statusCode,
		log,
	}
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	log        logging.Logger
}

func (rww *responseWriterWrapper) Status() int {
	return rww.statusCode
}

func (rww *responseWriterWrapper) WriteHeader(code int) {
	rww.log.Debugf("changing response status code from %d to %d", rww.statusCode, code)
	rww.statusCode = code
	rww.ResponseWriter.WriteHeader(code)
}
