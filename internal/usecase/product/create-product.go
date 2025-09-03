package product

import (
	"context"
	"time"

	"github.com/gregoryAlvim/gobid/internal/validator"
)

type CreateProductReq struct {
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "this field cannot be blank")

	eval.CheckField(validator.NotBlank(req.Description), "description", "this field cannot be blank")
	eval.CheckField((validator.MinChars(req.Description, 10) && validator.MaxChars(req.Description, 255)), "description", "this field must have a length between 10 and 255 characters")

	eval.CheckField(req.BasePrice > 0, "base_price", "this field cannot be zero")

	eval.CheckField(req.AuctionEnd.Sub(time.Now().UTC()) >= minAuctionDuration, "auction_end", "must be at least two hours duration")

	return eval
}
