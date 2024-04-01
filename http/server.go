package http

import (
	"fmt"
	"net/http"

	"github.com/SQUASHD/hbit/achievement"
	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/character"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/quest"
	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/user"
)

type (
	Server struct {
		serverConf config.ServerConfig
		jwtConf    config.JwtConfig
		authSvc    auth.Service
		charSvc    character.Service
		questSvc   quest.Service
		achSvc     achievement.Service
		taskSvc    task.Service
		userSvc    user.Service
	}
)

func NewServer(
	serverConf config.ServerConfig,
	jwtConf config.JwtConfig,
	authSvc auth.Service,
	charSvc character.Service,
	questSvc quest.Service,
	achievementSvc achievement.Service,
	taskSvc task.Service,
) (*http.Server, error) {

	NewServer := &Server{
		serverConf: serverConf,
		jwtConf:    jwtConf,
		authSvc:    authSvc,
		charSvc:    charSvc,
		questSvc:   questSvc,
		achSvc:     achievementSvc,
		taskSvc:    taskSvc,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConf.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  serverConf.TimeoutIdle,
		ReadTimeout:  serverConf.TimeoutRead,
		WriteTimeout: serverConf.TimeoutWrite,
	}
	return server, nil
}

func (s *Server) registerUserRoutes(mainRouter *http.ServeMux) http.Handler {
	return mainRouter
}

var ErrServerClosed = http.ErrServerClosed
