package config

type JwtConfig struct {
	JwtSecret                       string
	AccessDuration, RefreshDuration int
	AccessIssuer, RefreshIssuer     string
}

func NewJwtConfig(
	secret, accessIssuer, refreshIssuer string,
	accessDuration, refreshDuration int,
) JwtConfig {
	return JwtConfig{
		JwtSecret:       secret,
		AccessDuration:  accessDuration,
		RefreshDuration: refreshDuration,
		AccessIssuer:    accessIssuer,
		RefreshIssuer:   refreshIssuer,
	}
}
