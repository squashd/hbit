package auth

import "github.com/SQUASHD/hbit/auth/authdb"

type (
	// AuthDTO has the access token, refresh token, and username
	// Currently the http handler only returns the username as the response
	AuthDTO struct {
		Username     string `json:"username"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

func toDTO(model authdb.Auth, accessToken, refreshToken string) AuthDTO {
	return AuthDTO{
		Username:     model.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
