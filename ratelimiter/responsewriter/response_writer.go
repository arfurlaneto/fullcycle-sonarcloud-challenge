package responsewriter

import "net/http"

type RateLimiterResponseWriter interface {
	WriteResponse(w *http.ResponseWriter) error
	WriteError(w *http.ResponseWriter, err error) error
}
