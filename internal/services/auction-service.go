package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	// Requests
	PlaceBid MessageKind = iota

	// Success
	SuccessfullyPlacedBid

	// Errors
	FailedToPlaceBid
	InvalidJson

	// Infos
	NewBidPlaced
	AuctionFinished
)

type Message struct {
	Message string      `json:"message,omitempty"`
	Amount  float64     `json:"amount,omitempty"`
	Kind    MessageKind `json:"kind"`
	UserID  uuid.UUID   `json:"user_id,omitempty"`
}

type AuctionLobby struct {
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
	slog.Info("user disconnected", "client", c)
	delete(ar.Clients, c.UserId)
}

func (ar *AuctionRoom) broadcastMessage(m Message) {
	slog.Info("new message received", "room_id", ar.Id, "message", m, "user_id", m.UserID)
	switch m.Kind {
	case PlaceBid:
		bid, err := ar.BidsService.PlaceBid(ar.Context, ar.Id, m.UserID, m.Amount)
		if err != nil {
			if errors.Is(err, ErrBidTooLow) {
				if client, ok := ar.Clients[m.UserID]; ok {
					client.Send <- Message{Kind: FailedToPlaceBid, Message: ErrBidTooLow.Error(), UserID: m.UserID}
				}
				return
			}
		}

		if client, ok := ar.Clients[m.UserID]; ok {
			client.Send <- Message{Kind: SuccessfullyPlacedBid, Message: "Your bid was Successfully placed.", UserID: m.UserID}
		}

		for id, client := range ar.Clients {
			newBidMessage := Message{
				Kind:    NewBidPlaced,
				Message: "A new bid was placed", Amount: bid.BidAmount,
				UserID: m.UserID,
			}

			if id == m.UserID {
				continue
			}

			client.Send <- newBidMessage
		}

	case InvalidJson:
		client, ok := ar.Clients[m.UserID]
		if !ok {
			slog.Info("client not found", "user_id", m.UserID)
		}

		client.Send <- m
	}
}

func (ar *AuctionRoom) Run() {
	slog.Info("Auction has begun.", "auction_id", ar.Id)

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

const (
	maxMessageSize = 512
	readDeadline   = 60 * time.Second
	writeWait      = 10 * time.Second
	pingPeriod     = (readDeadline * 9) / 10
)

func (c *Client) ReadEventLoop() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(readDeadline))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(readDeadline))
		return nil
	})

	for {
		var m Message
		err := c.Conn.ReadJSON(&m)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("unexpected close error", "error", err)
				return
			}

			slog.Warn("invalid json received from client", "user_id", c.UserId, "error", err)
			continue
		}

		m.UserID = c.UserId
		c.Room.Broadcast <- m
	}
}

func (c *Client) WriteEventLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteJSON(Message{
					Kind:    websocket.CloseMessage,
					Message: "closing websocket connection",
				})
				return
			}

			if message.Kind == AuctionFinished {
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteJSON(message)
			if err != nil {
				c.Room.Unregister <- c
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("unexpected write error", "error", err)
				return
			}
		}
	}
}
