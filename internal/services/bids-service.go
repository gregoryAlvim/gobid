package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gregoryAlvim/gobid/internal/store/pgstore"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BidsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidsService(pool *pgxpool.Pool) BidsService {
	return BidsService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

var ErrBidTooLow = errors.New("bid amount is too low")

func (bs *BidsService) PlaceBid(ctx context.Context, product_id, bidder_id uuid.UUID, amount float64) (pgstore.Bid, error) {
	product, err := bs.queries.GetProductById(ctx, product_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, errors.New("product not found")
		}

		return pgstore.Bid{}, err
	}

	highestBid, err := bs.queries.GetHighestBidByProductId(ctx, product_id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	if product.BasePrice >= amount || highestBid.BidAmount >= amount {
		return pgstore.Bid{}, ErrBidTooLow
	}

	args := pgstore.CreateBidParams{
		ProductID: product_id,
		BidderID:  bidder_id,
		BidAmount: amount,
	}

	highestBid, err = bs.queries.CreateBid(ctx, args)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	return highestBid, nil
}
