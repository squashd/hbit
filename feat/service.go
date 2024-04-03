package feat

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/feat/featdb"
	"github.com/wagslane/go-rabbitmq"
)

type (
	service struct {
		db        *sql.DB
		queries   *featdb.Queries
		publisher *rabbitmq.Publisher
	}
)

func NewService(
	db *sql.DB,
	queries *featdb.Queries,
	publisher *rabbitmq.Publisher,
) Service {
	return &service{
		db:        db,
		queries:   queries,
		publisher: publisher,
	}
}

func (s *service) GetUserFeats(ctx context.Context, userID string) ([]UserFeatsDTO, error) {
	feats, err := s.queries.ListUserFeats(ctx, userID)
	if err != nil {
		return []UserFeatsDTO{}, err
	}

	dtos := toDTOs(feats)
	return dtos, nil
}

func (s *service) Publish(event hbit.EventMessage) error {

	msg, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	return s.publisher.Publish(
		msg,
		[]string{"feat"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
}

func (s *service) Cleanup() error {
	return s.db.Close()
}
