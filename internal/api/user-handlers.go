package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gregoryAlvim/gobid/internal/services"
	"github.com/gregoryAlvim/gobid/internal/usecase/user"
	"github.com/gregoryAlvim/gobid/internal/utils"
)

func (api *Api) handleSignUpUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[user.CreateUserReq](r)
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)

	fmt.Println("id: ", id)
	fmt.Println("err: ", err)

	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{"error": "email or password already exists"})
			return
		}
	}

	_ = utils.EncodeJson(w, r, http.StatusCreated, map[string]any{"user_id": id})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[user.LoginUserReq](r)
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			_ = utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{"error": "invalid email or password"})
			return
		}

		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected internal server error"})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected internal server error"})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticateUserId", id)
	utils.EncodeJson(w, r, http.StatusOK, map[string]any{"message": "logged in successfully"})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected internal server error"})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthenticateUserId")
	utils.EncodeJson(w, r, http.StatusOK, map[string]any{"message": "logged out successfully"})
}
