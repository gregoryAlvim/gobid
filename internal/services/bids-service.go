package services

import (
	"context"
	"errors"
	"log/slog"

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
	tx, err := bs.pool.Begin(ctx)
	if err != nil {
		return pgstore.Bid{}, err
	}

	defer tx.Rollback(ctx)

	qtx := bs.queries.WithTx(tx)

	product, err := qtx.GetProductById(ctx, product_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, errors.New("product not found")
		}

		return pgstore.Bid{}, err
	}

	highestBid, err := qtx.GetHighestBidByProductId(ctx, product_id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return pgstore.Bid{}, err
	}

	isFirstBid := errors.Is(err, pgx.ErrNoRows)
	if isFirstBid {
		if amount <= product.BasePrice {
			slog.Info("BID REJECTED: Amount is less than or equal to base price.")
			return pgstore.Bid{}, ErrBidTooLow
		}
	} else {
		if amount <= highestBid.BidAmount {
			slog.Info("BID REJECTED: Amount is less than or equal to highest bid.")
			return pgstore.Bid{}, ErrBidTooLow
		}
	}

	args := pgstore.CreateBidParams{
		ProductID: product_id,
		BidderID:  bidder_id,
		BidAmount: amount,
	}

	newBid, err := qtx.CreateBid(ctx, args)
	if err != nil {
		return pgstore.Bid{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return pgstore.Bid{}, err
	}

	return newBid, nil
}
