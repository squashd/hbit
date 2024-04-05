package updates

import "github.com/SQUASHD/hbit/rpg"

type Tag string

type taggedPayload struct {
	tag Tag
	any
}

const (
	TASKREWARD  Tag = "taskReward"
	QUESTREWARD Tag = "questReward"
	LEVELUP     Tag = "levelUp"
	UNKNOWN     Tag = "unknown"
)

func determineTag(payload any) Tag {
	switch payload.(type) {
	case rpg.TaskRewardPayload:
		return TASKREWARD
	case rpg.CharacterLevelUpPayload:
		return LEVELUP
	default:
		return UNKNOWN
	}
}

func tagPayload(payload any) taggedPayload {
	return taggedPayload{
		tag: determineTag(payload),
		any: payload,
	}
}
