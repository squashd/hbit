package rpg

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
	"github.com/wagslane/go-rabbitmq"
)

type (
	Service interface {
		EventService
	}

	EventService interface {
		HandleTaskCompleted(userId string) error
	}
	rpgService struct {
		charSvc   character.Service
		questSvc  quest.Service
		publisher *rabbitmq.Publisher
	}
)

func NewService(
	charSvc character.Service,
	questSvc quest.Service,
	publisher *rabbitmq.Publisher,
) EventService {
	return &rpgService{
		charSvc:   charSvc,
		questSvc:  questSvc,
		publisher: publisher,
	}
}

func (s *rpgService) HandleTaskCompleted(userId string) error {
	fmt.Printf("rpg service received task complete event for user: %s\n", userId)

	levelUpData := CharacterLevelUpPayload{
		Level: 69,
	}
	payload, err := json.Marshal(levelUpData)
	if err != nil {
		return err
	}

	levelUpEvent := hbit.EventMessage{
		Type:    "level_up",
		UserID:  userId,
		Payload: payload,
	}

	msg, err := json.Marshal(levelUpEvent)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(
		msg,
		[]string{"rpg.character"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		return err
	}

	fmt.Println("Published level up event")

	awardData := TaskRewardPayload{
		Gold: 50,
		Exp:  100,
		Mana: 10,
	}

	payload, err = json.Marshal(awardData)
	if err != nil {
		return err
	}

	rewardEvent := hbit.EventMessage{
		Type:    hbit.TaskRewardEvent,
		UserID:  userId,
		Payload: payload,
	}

	msg, err = json.Marshal(rewardEvent)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(
		msg,
		[]string{"rpg.rewards"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		return err
	}

	fmt.Println("Published rewards event")

	return nil
}

func (s *rpgService) CleanUp() error {
	var errs []error
	if err := s.charSvc.Cleanup(); err != nil {
		errs = append(errs, err)
	}
	if err := s.questSvc.Cleanup(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
