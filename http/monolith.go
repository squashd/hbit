package http

import (
	"github.com/SQUASHD/hbit/achievement"
	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/character"
	"github.com/SQUASHD/hbit/config"
	"github.com/SQUASHD/hbit/quest"
	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/user"
)

type (
	ServerMonolith struct {
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

func NewServerMonolith(
	serverConf config.ServerConfig,
	jwtConf config.JwtConfig,
	authSvc auth.Service,
	charSvc character.Service,
	questSvc quest.Service,
	achSvc achievement.Service,
	taskSvc task.Service,
	userSvc user.Service,
) *ServerMonolith {
	return &ServerMonolith{
		serverConf: serverConf,
		jwtConf:    jwtConf,
		authSvc:    authSvc,
		charSvc:    charSvc,
		questSvc:   questSvc,
		achSvc:     achSvc,
		taskSvc:    taskSvc,
		userSvc:    userSvc,
	}
}
