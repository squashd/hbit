package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/SQUASHD/hbit"
	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/rpg/character"
	"github.com/SQUASHD/hbit/rpg/rpgdb"
)

type (
	RegistrationOrchestrator interface {
		OrchestrateRegistration(w http.ResponseWriter, r *http.Request)
	}

	orchestratorReg struct {
		authSvc   auth.Service
		rpgSvcUrl string
		client    *http.Client
	}
)

func NewRegistrationOrchestrator(
	authSvc auth.Service,
	rpgSvcUrl string,
	client *http.Client,
) RegistrationOrchestrator {
	return &orchestratorReg{
		authSvc:   authSvc,
		rpgSvcUrl: rpgSvcUrl,
		client:    client,
	}
}
func (o *orchestratorReg) OrchestrateRegistration(w http.ResponseWriter, r *http.Request) {
	var form auth.CreateUserForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		Error(w, r, &hbit.Error{Code: hbit.EINVALID, Message: "Invalid JSON Body"})
		return
	}

	// All new users will be assigned the new_hero class
	// And are assigned a UUID as their user ID
	userId := hbit.NewUUID()
	classId := "new_hero"

	registrationForm := struct {
		auth.CreateUserForm
		character.CreateCharacterForm
	}{
		CreateUserForm: auth.CreateUserForm{
			UserID:          userId, // Previously we used SQLi to generate a UUID
			Username:        form.Username,
			Password:        form.Password,
			ConfirmPassword: form.ConfirmPassword,
		},
		CreateCharacterForm: character.CreateCharacterForm{
			CreateCharacterParams: rpgdb.CreateCharacterParams{
				UserID:  userId,
				ClassID: classId,
			},
			RequestedById: userId,
		},
	}

	var rpgRes *http.Response
	var authDto auth.AuthDTO
	var authErr, rpgErr error

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		authDto, authErr = o.callCreateAuth(r.Context(), registrationForm.CreateUserForm)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rpgRes, rpgErr = o.callCreateCharacter(registrationForm.CreateCharacterForm)
	}()

	wg.Wait()

	if authErr != nil && rpgErr == nil {
		LogError(r, authErr)
		go func() {
			// since authErr is not nil, we can assume that the user was not created
			// and we should not attempt to delete the user
			// instead we should publish an event to delete the user
			var pubErr error
			event, err := hbit.NewEventMessage(
				hbit.AUTHDELETE,
				hbit.UserId(userId),
				hbit.NewEventIdWithTimestamp("auth"),
				nil,
			)
			if err != nil {
				LogError(r, err)
				return
			}

			pubErr = o.authSvc.Publish(event, []string{"auth.delete"})
			if pubErr != nil {
				LogError(r, pubErr)
				return
			}
		}()
		Error(w, r, authErr)
		return
	}

	if rpgErr != nil && authErr == nil {
		LogError(r, rpgErr)
		go func() {
			// since rpgErr is not nil, we can assume that the character was not created
			// and thus the application is in an inconsistent state
			// we should delete the user from the auth service
			// the auth service does not subscribe to its own events so we call the
			// delete method directly
			delErr := o.callDeleteAuth(r.Context(), userId)
			if delErr != nil {
				LogError(r, delErr)
				return
			}
		}()
		Error(w, r, rpgErr)
		return
	}

	var charDTO character.DTO
	if err := json.NewDecoder(rpgRes.Body).Decode(&charDTO); err != nil {
		Error(w, r, err)
		return
	}

	registrationRes := struct {
		auth.AuthDTO
		character.DTO
	}{
		AuthDTO: authDto,
		DTO:     charDTO,
	}

	respondWithJSON(w, http.StatusCreated, registrationRes)
}

// A wrapper around the auth service's Register method in case we want to offload
// authentication to a separate service
func (o *orchestratorReg) callCreateAuth(ctx context.Context, form auth.CreateUserForm) (auth.AuthDTO, error) {
	userDto, err := o.authSvc.Register(ctx, form)
	if err != nil {
		return auth.AuthDTO{}, err
	}
	return userDto, nil
}

// A wrapper around the auth service's DeleteUser method in case we want to offload
// authentication to a separate service
func (o *orchestratorReg) callDeleteAuth(ctx context.Context, userId string) error {
	err := o.authSvc.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

// Since services delete user data based on event messages from the message broker
// we rely on the auth service publishing rather than a direct call to reverse the
// registration process
func (o *orchestratorReg) callCreateCharacter(form character.CreateCharacterForm) (*http.Response, error) {
	url := fmt.Sprintf("%s/characters", o.rpgSvcUrl)

	jsonData, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// This is done since we need to act as the user upon registration
	setUserIdInRequestHeader(req, form.RequestedById)

	res, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
