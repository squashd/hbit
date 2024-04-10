package feat

import (
	"time"

	"github.com/SQUASHD/hbit/feat/featdb"
)

type (
	UserFeatsDTO struct {
		Name        string     `json:"name"`
		Requirement string     `json:"requirement"`
		UserID      *string    `json:"user_id"`
		AchievedAt  *time.Time `json:"achieved_at"`
	}
)

func toDTO(row featdb.ListUserFeatsRow) UserFeatsDTO {
	return UserFeatsDTO{
		Name:        row.Name,
		Requirement: row.Requirement,
		UserID:      row.UserID,
		AchievedAt:  row.AchievedAt,
	}
}

func toDTOs(rows []featdb.ListUserFeatsRow) []UserFeatsDTO {
	items := make([]UserFeatsDTO, len(rows))
	for i, r := range rows {
		items[i] = toDTO(r)
	}
	return items
}
