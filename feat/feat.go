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
	Service interface {
		UserFeatService
		hbit.Publisher
		hbit.UserDataHandler
		CleanUp() error
	}
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

func (s *service) DeleteData(ctx context.Context, userId string) error {
	tx, er := s.db.Begin()
	if er != nil {
		return &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to start transaction"}
	}
	defer tx.Rollback()
	qtx := s.queries.WithTx(tx)
	err := qtx.DeleteUserQuestLogs(ctx, userId)
	if err != nil {
		return err
	}
	err = qtx.DeleteUserFeats(ctx, userId)
	if err != nil {
		return err
	}
	err = qtx.DeleteUserTaskLogs(ctx, userId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return &hbit.Error{Code: hbit.EINTERNAL, Message: "failed to commit transaction"}
	}
	return nil

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
	return s.db.Close()
}
