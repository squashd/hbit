package task

type TaskDifficulty string

const (
	EASY   TaskDifficulty = "easy"
	MEDIUM TaskDifficulty = "medium"
	HARD   TaskDifficulty = "hard"
	EPIC   TaskDifficulty = "epic"
)

type TaskDonePayload struct {
	TaskId     string         `json:"task_id"`
	Difficulty TaskDifficulty `json:"difficulty"`
}

type TaskUndonePayload struct {
	TaskId string `json:"task_id"`
}
