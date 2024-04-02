package config

import (
	"os"
	"time"
)

type JwtOptions struct {
	JwtSecret                       string
	AccessDuration, RefreshDuration time.Duration
	AccessIssuer, RefreshIssuer     string
}

func NewJwtConfig(options ...func(*JwtOptions),
) JwtOptions {
	opts := getDefaultJwtOptions()
	for _, option := range options {
		option(&opts)
	}
	return opts
}

func getDefaultJwtOptions() JwtOptions {
	return JwtOptions{
		JwtSecret:       "",
		AccessDuration:  time.Minute * 15,
		RefreshDuration: time.Minute * 24,
		AccessIssuer:    "access",
		RefreshIssuer:   "refresh",
	}
}

func WithJwtOptionsSecret(secret string) func(*JwtOptions) {
	return func(options *JwtOptions) {
		options.JwtSecret = secret
	}
}

func WithJwtOptionsAccessDuration(duration int) func(*JwtOptions) {
	return func(options *JwtOptions) {
		options.AccessDuration = time.Minute * time.Duration(duration)
	}
}

func WithJwtOptionsRefreshDuration(duration int) func(*JwtOptions) {
	return func(options *JwtOptions) {
		options.RefreshDuration = time.Minute * time.Duration(duration)
	}
}

func WithJwtOptionsAccessIssuer(issuer string) func(*JwtOptions) {
	return func(options *JwtOptions) {
		options.AccessIssuer = issuer
	}
}

func WithJwtOptionsRefreshIssuer(issuer string) func(*JwtOptions) {
	return func(options *JwtOptions) {
		options.RefreshIssuer = issuer
	}
}

func WithJwtOptionsSecretFromEnv(env string) func(*JwtOptions) {
	return func(options *JwtOptions) {
		secret := os.Getenv(env)
		options.JwtSecret = secret
	}
}
