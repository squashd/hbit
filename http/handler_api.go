package http

import (
	"net/http"

	"github.com/SQUASHD/hbit"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	Error(w, r, &hbit.Error{Code: hbit.ENOTFOUND, Message: "Could not find the resource you were looking for"})
}
