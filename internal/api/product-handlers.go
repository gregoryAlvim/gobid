package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gregoryAlvim/gobid/internal/usecase/product"
	"github.com/gregoryAlvim/gobid/internal/utils"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[product.CreateProductReq](r)
	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), "AuthenticateUserId").(uuid.UUID)
	if !ok {
		_ = utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected error, try again later"})
		return
	}

	id, err := api.ProductService.CreateProduct(r.Context(), userID, data.ProductName, data.Description, data.BasePrice, data.AuctionEnd)
	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "failed to create product auction, try again later"})
		return
	}

	utils.EncodeJson(w, r, http.StatusCreated, map[string]any{"product_id": id.String(), "message": "product auction created successfully"})
}
