package character

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	"github.com/wagslane/go-rabbitmq"
)

type (
	// Aggregate service for instantiation top level, should not be used directly
	Service interface {
		UserCharacterService
		AdminCharacterService
		InternalUserCharacterUtils
		hbit.Publisher
		CleanUp() error
	}

	UserCharacterService interface {
		GetCharacter(ctx context.Context, userId string) (CharacterDTO, error)
		CreateCharacter(ctx context.Context, form CreateCharacterForm) (CharacterDTO, error)
		UpdateCharacter(ctx context.Context, form UpdateCharacterForm) (CharacterDTO, error)
	}

	AdminCharacterService interface {
		List(ctx context.Context) ([]CharacterDTO, error)
	}

	service struct {
		db        *sql.DB
		queries   *rpgdb.Queries
		publisher *rabbitmq.Publisher
	}
)

func NewService(
	db *sql.DB,
	queries *rpgdb.Queries,
	publisher *rabbitmq.Publisher,
) Service {
	return &service{
		db:        db,
		queries:   queries,
		publisher: publisher,
	}
}

func (s *service) Publish(event hbit.EventMessage, routingKeys []string) error {
	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = s.publisher.Publish(
		msg,
		routingKeys,
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CleanUp() error {
	s.publisher.Close()
	return nil
}
