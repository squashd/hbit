package achievement

import (
	"time"

	"github.com/SQUASHD/hbit/achievement/database"
)

type (
	Achievement  = database.Achievement
	Achievements = []Achievement

	AchievementDTO struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Requirement string `json:"requirement"`
	}

	UserAchievement  = database.UserAchievement
	UserAchievements = []UserAchievement

	UserAchievementDTO struct {
		UserID        string    `json:"user_id"`
		AchievementID string    `json:"achievement_id"`
		CreatedAt     time.Time `json:"created_at"`
	}

	Mapper interface {
		achToDTO(ach Achievement) AchievementDTO
		achToDTOs(achs Achievements) []AchievementDTO

		userAchToDTO(userAch UserAchievement) UserAchievementDTO
		userAchToDTOs(userAchs UserAchievements) []UserAchievementDTO
	}
	mapper struct{}
)

func (m mapper) userAchToDTO(userAch database.UserAchievement) UserAchievementDTO {
	return UserAchievementDTO{
		UserID:        userAch.UserID,
		AchievementID: userAch.AchievementID,
		CreatedAt:     userAch.CreatedAt,
	}
}

func (m mapper) userAchToDTOs(userAchs []database.UserAchievement) []UserAchievementDTO {
	if len(userAchs) == 0 {
		return []UserAchievementDTO{}
	}
	dtos := make([]UserAchievementDTO, 0, len(userAchs))
	for _, userAch := range userAchs {
		dtos = append(dtos, m.userAchToDTO(userAch))
	}
	return dtos
}

func (m mapper) achToDTO(ach database.Achievement) AchievementDTO {
	return AchievementDTO{
		ID:          ach.ID,
		Name:        ach.Name,
		Requirement: ach.Requirement,
	}
}

func (m mapper) achToDTOs(achs []database.Achievement) []AchievementDTO {
	if len(achs) == 0 {
		return []AchievementDTO{}
	}
	dtos := make([]AchievementDTO, 0, len(achs))
	for _, ach := range achs {
		dtos = append(dtos, m.achToDTO(ach))
	}
	return dtos
}

func newMapper() Mapper {
	return mapper{}
}
