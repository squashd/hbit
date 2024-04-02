package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type ServerOptions struct {
	Port         int
	TimeoutRead  time.Duration
	TimeoutWrite time.Duration
	TimeoutIdle  time.Duration
	Debug        bool
}

func NewServer(options ...func(*ServerOptions) error) (*ServerOptions, error) {
	opts := getDefaultServerOptions()
	for _, option := range options {
		if err := option(&opts); err != nil {
			return nil, err
		}
	}
	return &opts, nil
}

func getDefaultServerOptions() ServerOptions {
	return ServerOptions{
		Port:         8080,
		TimeoutRead:  5 * time.Second,
		TimeoutWrite: 10 * time.Second,
		TimeoutIdle:  15 * time.Second,
		Debug:        false,
	}
}

func WithServerOptionsPort(port int) func(*ServerOptions) error {
	return func(options *ServerOptions) error {
		options.Port = port
		return nil
	}
}

func WithServerOptionsPortFromEnv(key string) func(*ServerOptions) error {
	return func(options *ServerOptions) error {
		if portStr, ok := os.LookupEnv(key); ok {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf("invalid port value: %v", err)
			}
			options.Port = port
		} else {
			return fmt.Errorf("missing port value")
		}
		return nil
	}
}

func WithServerOptionsTimeoutRead(timeout time.Duration) func(*ServerOptions) error {
	return func(options *ServerOptions) error {
		options.TimeoutRead = timeout
		return nil
	}
}

func WithServerOptionsTimeoutReadFromEnv(key string) func(*ServerOptions) error {
	return func(options *ServerOptions) error {
		if timeoutStr, ok := os.LookupEnv(key); ok {
			timeout, err := time.ParseDuration(timeoutStr)
			if err != nil {
				return fmt.Errorf("invalid timeout value: %v", err)
			}
			options.TimeoutRead = timeout
		} else {
			return fmt.Errorf("missing timeout value")
		}
		return nil
	}
}

func WithServerOptionsTimeoutWrite(timeout time.Duration) func(*ServerOptions) error {
	return func(options *ServerOptions) error {
		options.TimeoutWrite = timeout
		return nil
	}
}

func WithServerOptionsTimeoutWriteFromEnv(key string) func(*ServerOptions) error {
	return func(options *ServerOptions) error {
		if timeoutStr, ok := os.LookupEnv(key); ok {
			timeout, err := time.ParseDuration(timeoutStr)
			if err != nil {
				return fmt.Errorf("invalid timeout value: %v", err)
			}
			options.TimeoutWrite = timeout
		} else {
			return fmt.Errorf("missing timeout value")
		}
		return nil
	}
}

func WithServerOptionsTimeoutIdle(timeout time.Duration) func(*ServerOptions) error {
	return func(options *ServerOptions) error {
		options.TimeoutIdle = timeout
		return nil
	}
}

func WithServerOptionsTimeoutIdleFromEnv(key string) func(*ServerOptions) error {
	return func(options *ServerOptions) error {
		if timeoutStr, ok := os.LookupEnv(key); ok {
			timeout, err := time.ParseDuration(timeoutStr)
			if err != nil {
				return fmt.Errorf("invalid timeout value: %v", err)
			}
			options.TimeoutIdle = timeout
		} else {
			return fmt.Errorf("missing timeout value")
		}
		return nil
	}
}
