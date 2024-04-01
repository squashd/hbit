package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SQUASHD/hbit"
)

var codes = map[hbit.AppError]int{
	hbit.ECONFLICT:       http.StatusConflict,
	hbit.EINVALID:        http.StatusBadRequest,
	hbit.ENOTFOUND:       http.StatusNotFound,
	hbit.ENOTIMPLEMENTED: http.StatusNotImplemented,
	hbit.EUNAUTHORIZED:   http.StatusUnauthorized,
	hbit.EFORBIDDEN:      http.StatusForbidden,
	hbit.EINTERNAL:       http.StatusInternalServerError,
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
				fmt.Printf("Error: %v\n", e)
			}
		}
		RespondWithJSON(w, ErrorStatusCode(hbit.EINVALID), ErrorResponse{Error: strings.Join(messages, ", ")})
		return

	default:
		code, message := hbit.ErrorCode(err), hbit.ErrorMessage(err)
		if code == hbit.EINTERNAL {
			fmt.Printf("Error: %v\n", err)
		}
		RespondWithJSON(w, ErrorStatusCode(code), ErrorResponse{Error: message})
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
