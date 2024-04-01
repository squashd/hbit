package config

import (
	"os"
	"strconv"
	"time"

	"github.com/SQUASHD/hbit"
)

type ServerConfig struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
}

func NewServerConfigFromEnv() (ServerConfig, error) {
	var c ServerConfig

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_PORT is not set"}
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_PORT is not set"}
	}

	timeoutRead := os.Getenv("SERVER_TIMEOUT_READ")
	if timeoutRead == "" {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_TIMEOUT_READ is not set"}
	}
	tr, err := time.ParseDuration(timeoutRead)
	if err != nil {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_TIMEOUT_READ is not a valid duration"}
	}

	timeoutWrite := os.Getenv("SERVER_TIMEOUT_WRITE")
	if timeoutWrite == "" {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_TIMEOUT_WRITE is not set"}
	}
	tw, err := time.ParseDuration(timeoutWrite)
	if err != nil {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_TIMEOUT_WRITE is not a valid duration"}
	}

	timeoutIdle := os.Getenv("SERVER_TIMEOUT_IDLE")
	if timeoutIdle == "" {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_TIMEOUT_IDLE is not set"}
	}
	ti, err := time.ParseDuration(timeoutIdle)
	if err != nil {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_TIMEOUT_IDLE is not a valid duration"}

	}

	debugStr := os.Getenv("SERVER_DEBUG")
	if debugStr == "" {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_DEBUG is not set"}
	}
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		return ServerConfig{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "SERVER_DEBUG is not a valid boolean"}
	}

	c.Port = port
	c.TimeoutRead = tr
	c.TimeoutWrite = tw
	c.TimeoutIdle = ti
	c.Debug = debug

	return c, nil
}
