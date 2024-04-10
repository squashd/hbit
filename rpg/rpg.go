package rpg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	"github.com/wagslane/go-rabbitmq"
)

type (
	Service interface {
		CalculateRewards(context.Context, TaskRewardRequest) (TaskRewardResponse, error)
		UndoRewards(context.Context, TaskRewardRequest) (UnresolvedTaskPayload, error)
		hbit.Publisher
		hbit.UserDataHandler
		CleanUp() error
	}
	rpgService struct {
		charSvc   character.Service
		questSvc  quest.Service
		publisher *rabbitmq.Publisher
		queries   *rpgdb.Queries
		db        *sql.DB
	}
)

type NewServiceParams struct {
	CharacterSvc character.Service
	QuestSvc     quest.Service
	Publisher    *rabbitmq.Publisher
	Queries      *rpgdb.Queries
	Db           *sql.DB
}

func NewService(
	params NewServiceParams,
) Service {
	return &rpgService{
		charSvc:   params.CharacterSvc,
		questSvc:  params.QuestSvc,
		publisher: params.Publisher,
		queries:   params.Queries,
		db:        params.Db,
	}
}

func (s *rpgService) DeleteData(ctx context.Context, userId string) error {
	tx, er := s.db.Begin()
	if er != nil {
		return &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to start transaction"}
	}
	defer tx.Rollback()
	qtx := s.queries.WithTx(tx)
	err := qtx.DeleteUserQuestData(ctx, userId)
	if err != nil {
		return err
	}
	err = qtx.DeleteUserCharacterData(ctx, userId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to commit transaction"}
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

func (s *rpgService) CleanUp() error {
	var errs []error
	s.publisher.Close()
	if err := s.charSvc.CleanUp(); err != nil {
		errs = append(errs, err)
	}
	if err := s.db.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.questSvc.CleanUp(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
