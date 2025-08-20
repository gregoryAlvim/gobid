package api

import (
	"net/http"

	"github.com/gregoryAlvim/gobid/internal/jsonutils"
	"github.com/gregoryAlvim/gobid/internal/usecase/user"
)

func (api *Api) handleSignUpUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
	}
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {}
