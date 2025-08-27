-- name: CreateBid :one
INSERT INTO bids ("product_id", "bidder_id", "bid_amount")
VALUES ($1, $2, $3) 
RETURNING *;

-- name: GetBidsByProductId :many
SELECT id, product_id, bidder_id, bid_amount, created_at FROM bids WHERE product_id = $1 ORDER BY bid_amount DESC;

-- name: GetHighestBidByProductId :one
SELECT id, product_id, bidder_id, bid_amount, created_at FROM bids WHERE product_id = $1 ORDER BY bid_amount DESC LIMIT 1;