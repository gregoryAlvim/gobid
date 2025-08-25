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
	fmt.Println("err: ", err)
	fmt.Println("problems: ", problems)
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)

	fmt.Println("id: ", id)
	fmt.Println("err: ", err)

	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrPassword) {
			_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{"error": "email or password already exists"})
			return
		}
	}

	_ = utils.EncodeJson(w, r, http.StatusCreated, map[string]any{"user_id": id})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {}
