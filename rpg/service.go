package rpg

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	"github.com/SQUASHD/hbit/task"
	"github.com/wagslane/go-rabbitmq"
)

type (
	EventService interface {
		HandleTaskCompleted(userId hbit.UserId, difficulty task.TaskDifficulty) error
		hbit.Publisher
		hbit.UserDataHandler
	}
	rpgService struct {
		publisher *rabbitmq.Publisher
		queries   *rpgdb.Queries
		db        *sql.DB
	}
)

// DeleteData implements EventService.
func (s *rpgService) DeleteData(userId string) error {
	tx, er := s.db.Begin()
	if er != nil {
		return &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to start transaction"}
	}
	defer tx.Rollback()
	qtx := s.queries.WithTx(tx)
	err := qtx.DeleteUserData(context.Background(), userId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to commit transaction"}
	}
	return nil

}

func NewService(
	publisher *rabbitmq.Publisher,
	queries *rpgdb.Queries,
	db *sql.DB,
) EventService {
	return &rpgService{
		publisher: publisher,
		queries:   queries,
		db:        db,
	}
}

func (s *rpgService) HandleTaskCompleted(
	userId hbit.UserId,
	difficulty task.TaskDifficulty,
) error {
	char, err := s.queries.ReadCharacter(context.Background(), userId.String())
	if err != nil {
		return err
	}

	// TODO: check character's current quest and state
	reward := determineReward(char, difficulty)
	msg, err := hbit.NewEventMessage(
		hbit.RPGREWARD,
		userId,
		hbit.NewEventIdWithTimestamp("rpg"),
		reward,
	)
	if err != nil {
		return err
	}

	err = s.Publish(msg, []string{"rpg.reward"})
	if err != nil {
		return err
	}
	return nil
}

func (s *rpgService) Publish(event hbit.EventMessage, routingKeys []string) error {
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

func (s *rpgService) CleanUp() {
	s.publisher.Close()
}
