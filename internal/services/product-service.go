package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gregoryAlvim/gobid/internal/store/pgstore"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(
	ctx context.Context,
	sellerId uuid.UUID,
	productName,
	description string,
	base_price float64,
	auction_end time.Time,
) (uuid.UUID, error) {
	args := pgstore.CreateProductParams{
		SellerID:    sellerId,
		ProductName: productName,
		Description: description,
		BasePrice:   base_price,
		AuctionEnd:  auction_end,
	}

	id, err := ps.queries.CreateProduct(ctx, args)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
