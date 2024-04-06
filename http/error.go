package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/SQUASHD/hbit"
)

// This centralization of error handling is lifted directly from wtf:
// https://github.com/benbjohnson/wtf
// I really like it
var codes = map[hbit.AppError]int{
	hbit.ECONFLICT:       http.StatusConflict,
	hbit.EINVALID:        http.StatusBadRequest,
	hbit.ENOTFOUND:       http.StatusNotFound,
	hbit.ENOTIMPLEMENTED: http.StatusNotImplemented,
	hbit.EUNAUTHORIZED:   http.StatusUnauthorized,
	hbit.EFORBIDDEN:      http.StatusForbidden,
	hbit.EINTERNAL:       http.StatusInternalServerError,
	hbit.EASYNC:          http.StatusAccepted,
}

func ErrorStatusCode(code hbit.AppError) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

func FromErrorStatusCode(code int) hbit.AppError {
	for k, v := range codes {
		if v == code {
			return k
		}
	}
	return hbit.EINTERNAL
}

// Error responds with an error message.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	switch err := err.(type) {

	case *hbit.MultiError:
		var messages []string
		for _, e := range err.Errors {
			if e.Code != hbit.EINTERNAL {
				messages = append(messages, e.Message)
			} else {
				LogError(r, e)
			}
		}
		respondWithJSON(w, ErrorStatusCode(hbit.EINVALID), ErrorResponse{Error: strings.Join(messages, ", ")})
		return

	default:
		code, message := hbit.ErrorCode(err), hbit.ErrorMessage(err)
		// TODO: Remove this once the app has fewer unknown errors.
		LogError(r, err)
		if code == hbit.EINTERNAL {
			hbit.ReportError(r.Context(), err, r)

		}
		respondWithJSON(w, ErrorStatusCode(code), ErrorResponse{Error: message})
	}
}

func LogError(r *http.Request, err error) {
	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var ErrServerClosed = http.ErrServerClosed
