package quest

import (
	"context"
)

// User facing service
type UserQuestService interface {
	ListQuests(ctx context.Context, userId string) ([]QuestDTO, error)
	// StartQuest(ctx context.Context, userId, questId string) error
	// StopQuest(ctx context.Context, userId, questId string) error
}

func (s *service) ListQuests(ctx context.Context, userId string) ([]QuestDTO, error) {
	quests, err := s.queries.ListQuests(ctx)
	if err != nil {
		return nil, err
	}

	return questsToDTOs(quests), nil
}

//func (s *service) ReadQuest(ctx context.Context, id string) (QuestDTO, error) {
//	return QuestDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
//}
//
//func (s *service) CreateQuest(ctx context.Context, data any) (QuestDTO, error) {
//	return QuestDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
//}
//
//func (s *service) UpdateQuest(ctx context.Context, data any) (QuestDTO, error) {
//	return QuestDTO{}, &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
//
//}
//
//func (s *service) DeleteQuest(ctx context.Context, id string) error {
//	return &hbit.Error{Code: hbit.EINTERNAL, Message: "Not implemented"}
//
//}
