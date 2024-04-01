package achievement

import (
	"context"
	"encoding/json"
	"net/http"
)

type (
	HttpService interface {
		GetUserAchievements(ctx context.Context, userID string) (UserAchievements, error)
	}

	MessageService interface {
		TaskComplet(msg json.RawMessage) (string, error)
		UserDeleted(msg json.RawMessage) (string, error)
	}

	Service interface {
		HttpService
		MessageService
	}

	Handler interface {
		GetAchievements(w http.ResponseWriter, r *http.Request, userId string)
	}

	handler struct {
		service Service
	}
)

func NewHandler(s Service) Handler {
	return &handler{
		service: s,
	}
}

func (h *handler) GetAchievements(w http.ResponseWriter, r *http.Request, userId string) {
}
