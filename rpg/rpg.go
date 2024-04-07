package rpg

import "github.com/SQUASHD/hbit/task"

type (
	CaclulateRewardPayload struct {
		Difficulty task.TaskDifficulty `json:"difficulty"`
		UserId     string              `json:"userId"`
	}
)
