package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gregoryAlvim/gobid/internal/services"
	"github.com/gregoryAlvim/gobid/internal/usecase/product"
	"github.com/gregoryAlvim/gobid/internal/utils"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[product.CreateProductReq](r)
	if err != nil {
		utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), "AuthenticateUserId").(uuid.UUID)
	if !ok {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected error, try again later"})
		return
	}

	productId, err := api.ProductService.CreateProduct(r.Context(), userID, data.ProductName, data.Description, data.BasePrice, data.AuctionEnd)
	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "failed to create product auction, try again later"})
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)
	auctionRoom := services.NewAuctionRoom(ctx, productId, api.BidsService)

	go auctionRoom.Run()
	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[productId] = auctionRoom
	api.AuctionLobby.Unlock()

	utils.EncodeJson(w, r, http.StatusCreated, map[string]any{"message": "Auction has started with success", "product_id": productId.String()})
}
