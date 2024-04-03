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
		Publish(event hbit.EventMessage, routingKeys []string) error
	}
	rpgService struct {
		charSvc   character.CharacterService
		questSvc  quest.QuestService
		publisher *rabbitmq.Publisher
	}
)

func NewService(
	charSvc character.CharacterService,
	questSvc quest.QuestService,
	publisher *rabbitmq.Publisher,
) EventService {
	return &rpgService{
		charSvc:   charSvc,
		questSvc:  questSvc,
		publisher: publisher,
	}
}

func (s *rpgService) HandleTaskCompleted(userId string) error {

	levelUpData := CharacterLevelUpPayload{
		Level: 69,
	}
	levelMsg, err := hbit.NewEventMessage(
		hbit.CharacterLevelUpEvent,
		userId,
		hbit.NewUUID(),
		levelUpData,
	)
	err = s.Publish(levelMsg, []string{"rpg.levelup"})
	if err != nil {
		return err
	}

	awardData := TaskRewardPayload{
		Gold: 50,
		Exp:  100,
		Mana: 10,
	}
	rewardMsg, err := hbit.NewEventMessage(
		hbit.CharacterLevelUpEvent,
		userId,
		hbit.NewUUID(),
		awardData,
	)
	err = s.Publish(rewardMsg, []string{"rpg.rewards"})
	if err != nil {
		return err
	}

	fmt.Println("Published rewards event")

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
	if err := s.charSvc.CleanUp(); err != nil {
		errs = append(errs, err)
	}
	if err := s.questSvc.CleanUp(); err != nil {
		errs = append(errs, err)
	}
	s.publisher.Close()
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
