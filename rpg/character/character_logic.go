package character

import (
	"context"
	"log"
	"sync"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

const (
	BASE_DMG_DONE = 10
)

type InternalUserCharacterUtils interface {
	CalculateBaseCharacterOuput(rpgdb.CharacterState) CharacterOutput
	HandleCharacterChange(context.Context, rpgdb.CharacterState, CharacterProgression) (CharacterDTO, bool, error)
}

type CharacterProgression struct {
	ExperienceChange int
	GoldChange       int
	ManaChange       int
}

func (s *service) HandleCharacterChange(
	ctx context.Context,
	charState rpgdb.CharacterState,
	progress CharacterProgression,
) (CharacterDTO, bool, error) {
	progressedChar, levelUp := progressCharacter(charState, progress)

	// Source eventId at orchestrator?
	eventId := hbit.NewEventIdWithTimestamp("rpg")

	updatedChar, err := s.queries.UpdateCharacter(ctx, rpgdb.UpdateCharacterParams{
		UserID:         progressedChar.UserID,
		ClassID:        progressedChar.ClassID,
		EventID:        string(eventId),
		CharacterLevel: progressedChar.CharacterLevel,
		Experience:     progressedChar.Experience,
		Health:         progressedChar.Health,
		Mana:           progressedChar.Mana,
		Strength:       progressedChar.Strength,
		Dexterity:      progressedChar.Dexterity,
		Intelligence:   progressedChar.Intelligence,
	})
	if err != nil {
		return CharacterDTO{}, false, err
	}

	if levelUp {
		levelUpPayload := CharacterLevelUpPayload{
			Level: updatedChar.CharacterLevel,
		}

		event, err := hbit.NewEventMessage(hbit.NewEventMessageParams{
			EventType: hbit.LEVELUP,
			UserId:    updatedChar.UserID,
			EventId:   eventId,
			Payload:   levelUpPayload,
		})
		if err != nil {
			log.Printf("Failed to create event message: %v", err)
		}
		if err := s.Publish(event, []string{"character.levelup"}); err != nil {
			log.Println("Failed to publish level up event", err)
		}
	}

	dto := characterToDto(updatedChar)

	return dto, levelUp, nil
}

func progressCharacter(charState rpgdb.CharacterState, progress CharacterProgression) (rpgdb.CharacterState, bool) {
	var levelUp = progress.ExperienceChange > calculateExpToNextLevel(charState)
	if levelUp {
		level := charState.CharacterLevel
		charState.CharacterLevel++
		charState.Mana += manaProgression(level, charState.ClassID)
		charState.Strength += strengthProgression(level, charState.ClassID)
		charState.Dexterity += dexProgression(level, charState.ClassID)
		charState.Intelligence += intProgression(level, charState.ClassID)

	}
	charState.Mana += int64(progress.ManaChange)
	charState.Experience += int64(progress.ExperienceChange)

	// Need to add a player inventory
	// charState.Gold += int64(progress.GoldChange)
	return charState, levelUp
}

// This will be based on the characters's class
// Rogues more dex, wizards more int, warriors more strength
func manaProgression(level int64, classId string) int64 {
	switch classId {
	default:
		return level * 10
	}
}

func strengthProgression(level int64, classId string) int64 {
	switch classId {
	default:
		return level * 10
	}
}

func dexProgression(level int64, classId string) int64 {
	switch classId {
	default:
		return level * 10
	}
}

func intProgression(level int64, classId string) int64 {
	switch classId {
	default:
		return level * 10
	}
}

func calculateExpToNextLevel(char rpgdb.CharacterState) int {
	return int((char.CharacterLevel+1)*1000 - char.Experience)
}

type CharacterOutput struct {
	DamageDone, Drops int
}

func (s *service) CalculateBaseCharacterOuput(char rpgdb.CharacterState) CharacterOutput {
	var damage, drops int
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		damage = calculateDamageDone(char)
	}()
	go func() {
		defer wg.Done()
		drops = calculateDrops(char)
	}()
	wg.Wait()

	return CharacterOutput{
		DamageDone: damage,
		Drops:      drops,
	}
}

// TODO: class-specific damage ouput
func calculateDamageDone(char rpgdb.CharacterState) int {
	dexModifier := 1 + (float64(char.Dexterity) / 100)
	strModifier := 1 + (float64(char.Strength) / 100)
	levelModifier := 1 + (float64(char.CharacterLevel) / 100)

	return int(float64(BASE_DMG_DONE) * dexModifier * strModifier * levelModifier)
}

func calculateDrops(_ rpgdb.CharacterState) int {
	return 1
}
