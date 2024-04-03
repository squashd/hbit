package rpg

type (
	CharacterLevelUpPayload struct {
		Level int `json:"level"`
	}

	TaskRewardPayload struct {
		Gold int `json:"gold"`
		Exp  int `json:"exp"`
		Mana int `json:"mana"`
	}
)
