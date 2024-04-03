package rpg

import (
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
)

type Repository interface {
	character.Repository
	quest.Repository
}
