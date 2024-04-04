package updates

import "github.com/SQUASHD/hbit/rpg"

type taggedPayload struct {
	tag string
	any
}

func determineTag(payload any) string {
	switch payload.(type) {
	case rpg.TaskRewardPayload:
		return "reward"
	case rpg.CharacterLevelUpPayload:
		return "levelup"
	default:
		return "unknown"
	}
}

func tagPayload(payload any) taggedPayload {
	return taggedPayload{
		tag: determineTag(payload),
		any: payload,
	}
}
