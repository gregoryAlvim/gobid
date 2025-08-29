package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	// Requests
	PlaceBid MessageKind = iota

	// Success
	SuccessfyllyPlacedBid

	// Errors
	FailedToPlaceBid

	// Infos
	NewBidPlaced
	AuctionFinished
)

type Message struct {
	Message string
	Amount  float64
	Kind    MessageKind
	UserID  uuid.UUID
}

type Auctionlobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id          uuid.UUID
	Context     context.Context
	Broadcast   chan Message
	Register    chan *Client
	Unregister  chan *Client
	Clients     map[uuid.UUID]*Client
	BidsService BidsService
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, bidsService BidsService) *AuctionRoom {
	return &AuctionRoom{
		Id:          id,
		Context:     ctx,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[uuid.UUID]*Client),
		BidsService: bidsService,
	}
}

func (ar *AuctionRoom) registerClient(c *Client) {
	slog.Info("new user connected", "client", c)
	ar.Clients[c.UserId] = c
}

func (ar *AuctionRoom) unregisterClient(c *Client) {
	slog.Info("user desconected", "client", c)
	delete(ar.Clients, c.UserId)
}

func (ar *AuctionRoom) broadcastMessage(m Message) {
	slog.Info("new message recieved", "room_id", ar.Id, "message", m, "user_id", m.UserID)
	switch m.Kind {
	case PlaceBid:
		bid, err := ar.BidsService.PlaceBid(ar.Context, ar.Id, m.UserID, m.Amount)
		client, ok := ar.Clients[m.UserID]

		if err != nil {
			if errors.Is(err, ErrBidTooLow) {
				if ok {
					client.Send <- Message{
						Kind:    FailedToPlaceBid,
						Message: ErrBidTooLow.Error(),
					}
				}
			}

			return
		}

		client.Send <- Message{
			Kind:    SuccessfyllyPlacedBid,
			Message: "your bid was placed with success ",
		}

		for id, client := range ar.Clients {
			newBidMessage := Message{
				Kind:    NewBidPlaced,
				Message: "A new bid was placed", Amount: bid.BidAmount,
			}

			if id == m.UserID {
				continue
			}

			client.Send <- newBidMessage
		}
	}
}

func (ar *AuctionRoom) Run() {
	defer func() {
		close(ar.Broadcast)
		close(ar.Register)
		close(ar.Unregister)
	}()

	for {
		select {
		case client := <-ar.Register:
			ar.registerClient(client)
		case client := <-ar.Unregister:
			ar.unregisterClient(client)
		case message := <-ar.Broadcast:
			ar.broadcastMessage(message)
		case <-ar.Context.Done():
			slog.Info("Auction has ended.", "auction_id", ar.Id)

			for _, client := range ar.Clients {
				client.Send <- Message{
					Kind:    AuctionFinished,
					Message: "The auction has ended. Thank you for participating!",
				}
			}

			return
		}
	}
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserId uuid.UUID
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userId uuid.UUID) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserId: userId,
	}
}
