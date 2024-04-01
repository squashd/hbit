package http

import (
	"fmt"
	"net/http"

	"github.com/SQUASHD/hbit/config"
)

func NewServer(serverConf config.ServerConfig, handler http.Handler) (*http.Server, error) {

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConf.Port),
		Handler:      handler,
		IdleTimeout:  serverConf.TimeoutIdle,
		ReadTimeout:  serverConf.TimeoutRead,
		WriteTimeout: serverConf.TimeoutWrite,
	}
	return server, nil
}

var ErrServerClosed = http.ErrServerClosed
