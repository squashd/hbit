package rpg

import (
	"context"
	"log"
	"math"
	"sync"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/quest"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
	"github.com/SQUASHD/hbit/task"
)

type (
	TaskRewardRequest struct {
		Difficulty task.TaskDifficulty `json:"difficulty"`
		UserId     string              `json:"userId"`
	}

	TaskRewardResponse struct {
		BaseReward  BaseTaskReward         `json:"base_reward"`
		QuestReward quest.QuestRewards     `json:"quest_reward"`
		Character   character.CharacterDTO `json:"character"`
		QuestDTO    quest.UserQuestDTO     `json:"quest_dto"`
		LevelUp     bool                   `json:"level_up"`
	}
)

// CalculateRewards calculates the rewards for a task based on the character's current state.
// It returns the rewards and an error if any.
func (s *rpgService) CalculateRewards(
	ctx context.Context,
	rewardPayload TaskRewardRequest,
) (TaskRewardResponse, error) {
	var payload TaskRewardResponse

	userId := rewardPayload.UserId
	difficulty := rewardPayload.Difficulty

	// Since there is no foreign key constraint we simply have to pass in the
	// userId twice to the query
	dbParams := rpgdb.GetUserRpgStateParams{
		UserID:   userId,
		UserID_2: userId,
	}

	// We do one query to get all the data we need to process a task done event
	rpgState, err := s.queries.GetUserRpgState(ctx, dbParams)
	if err != nil {
		return payload, err
	}

	// After getting a row we fracture it up to delegate down to the appropriate services
	charData, questData := convertRPGStateToCharAndQuest(rpgState)

	var baseReward BaseTaskReward
	var questState quest.QuestState
	var questStateErr error
	var charOutput character.CharacterOutput

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		baseReward = calculateBaseTaskReward(charData, difficulty)
	}()
	go func() {
		defer wg.Done()
		questState, questStateErr = s.questSvc.GetQuestState(questData)
	}()
	go func() {
		defer wg.Done()
		// This is a simple calculation of the character's output of damage
		// and the total drop amount
		charOutput = s.charSvc.CalculateBaseCharacterOuput(charData)
	}()
	wg.Wait()

	if questStateErr != nil {
		hbit.ReportError(ctx, questStateErr) // uh oh, spaghetti-o
		log.Printf("Error getting questState: %s", questStateErr.Error())
	}
	// What now??

	resolveInput := taskResolutionInput{
		CharacterOutput: charOutput,
		QuestState:      questState,
	}

	taskRes := processTaskCompletion(resolveInput, baseReward)

	var charDto character.CharacterDTO
	var questDto quest.UserQuestDTO
	var charErr, questErr error
	var levelup bool

	wg.Add(2)
	go func() {
		defer wg.Done()
		charDto, levelup, charErr = s.charSvc.HandleCharacterChange(ctx, charData, taskRes.CharProgress)
	}()
	go func() {
		defer wg.Done()
		questDto, questErr = s.questSvc.HandleQuestChange(questData, taskRes.QuestProgress)
	}()
	wg.Wait()

	if charErr != nil {
		return payload, charErr
	}

	if questErr != nil {
		return payload, questErr
	}

	questRewards := questState.QuestRewards

	finalRewards := TaskRewardResponse{
		BaseReward: baseReward,
		QuestReward: quest.QuestRewards{
			Gold: questRewards.Gold,
			Exp:  questRewards.Exp,
		},
		Character: charDto,
		LevelUp:   levelup,
		QuestDTO:  questDto,
	}

	return finalRewards, nil
}

type UnresolvedTaskPayload struct {
	BaseReward BaseTaskReward
}

// How's this going to work?
func (s *rpgService) UndoRewards(
	ctx context.Context,
	request TaskRewardRequest,
) (UnresolvedTaskPayload, error) {
	userId := request.UserId
	_, err := s.queries.ReadCharacter(context.Background(), userId)
	if err != nil {
		return UnresolvedTaskPayload{}, err
	}

	return UnresolvedTaskPayload{}, nil
}

type (
	taskResolutionInput struct {
		character.CharacterOutput
		quest.QuestState
	}

	TaskResolution struct {
		QuestProgress quest.QuestProgression
		CharProgress  character.CharacterProgression
	}
)

// We take the inputs of quest state and character output to calculate the rewards
// This separates the concerns of the quest and character service but use their data
func processTaskCompletion(input taskResolutionInput, baseReward BaseTaskReward) TaskResolution {
	var questXp int

	questProgress := calculateQuestProgress(input)

	if questProgress.Completed {
		questXp = input.QuestRewards.Exp
	}

	totalXp := calculateTotalXp(baseReward.Exp, questXp)

	charProgress := character.CharacterProgression{
		ExperienceChange: totalXp,
		GoldChange:       baseReward.Gold,
		ManaChange:       baseReward.Mana,
	}

	return TaskResolution{
		QuestProgress: questProgress,
		CharProgress:  charProgress,
	}
}

// We measure the character output against the quest state
func calculateQuestProgress(input taskResolutionInput) quest.QuestProgression {
	progress := quest.QuestProgression{
		Completed: questCompletionCheck(input),
	}

	// A player can't get more drops or do more damage than the remaining amount
	// Hencethe max in calcDrops and calcDamage
	if input.DropsNeeded != nil {
		progress.DropChange = calcDrops(input)
	}
	if input.DamageNeeded != nil {
		progress.BossHealthChange = calcDamage(input)
	}

	return progress
}

// If the quest isn't a drop quest, then the drops are 0
func calcDrops(input taskResolutionInput) int {
	if input.DropsNeeded == nil {
		return 0
	}
	return int(math.Max(float64(*input.DropsNeeded-input.Drops), 0))
}

// If the quest isn't a damage quest, then the damage done is 0
func calcDamage(input taskResolutionInput) int {
	if input.DamageNeeded == nil {
		return 0
	}
	return int(math.Max(float64(*input.DamageNeeded-input.DamageDone), 0))
}

// If a quest is inactive then it's not possible to complete it
func questCompletionCheck(input taskResolutionInput) bool {

	if !input.Active {
		return false
	}
	// Tasks are either drops or damage
	if input.DropsNeeded != nil && *input.DropsNeeded-input.Drops <= 0 {
		return true
	}
	if input.DamageNeeded != nil && *input.DamageNeeded-input.DamageDone <= 0 {
		return true
	}
	return false
}

// If there are different sources to calculate the base reward
// we can just use a variadic function to sum them up
// might be more complicated down the road?
func calculateTotalXp(exp ...int) int {
	var total int
	for _, e := range exp {
		total += e
	}

	return total
}

func processTaskUndo(_ taskResolutionInput) UnresolvedTaskPayload {
	return UnresolvedTaskPayload{}
}
