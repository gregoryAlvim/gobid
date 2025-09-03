package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gregoryAlvim/gobid/internal/services"
	"github.com/gregoryAlvim/gobid/internal/utils"
)

func (api *Api) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductID := chi.URLParam(r, "product_id")

	productId, err := uuid.Parse(rawProductID)
	if err != nil {
		utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{"message": "invalid product id, must be a valid uuid"})
		return
	}

	_, err = api.ProductService.GetProductById(r.Context(), productId)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			utils.EncodeJson(w, r, http.StatusNotFound, map[string]any{"message": "no product found with the given id"})
			return
		}

		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"message": "unexpected error, try again later"})
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticateUserId").(uuid.UUID)
	if !ok {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"message": "unexpected error, try again later"})
		return
	}

	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productId]
	api.AuctionLobby.Unlock()

	if !ok {
		utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{"message": "the auction for this product has ended or does not exist"})
		return
	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"message": "could not upgrade connection to websocket protocol"})
		return
	}

	client := services.NewClient(room, conn, userId)

	room.Register <- client
	go client.ReadEventLoop()
	go client.WriteEventLoop()
}
