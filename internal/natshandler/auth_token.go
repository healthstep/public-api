package natshandler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/porebric/resty/ws"
)

type AuthTokenMessage struct {
	Key    string `json:"key"`
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

type AuthTokenHandler struct {
	hub *ws.Hub
}

func NewAuthTokenHandler(hub *ws.Hub) *AuthTokenHandler {
	return &AuthTokenHandler{hub: hub}
}

func (h *AuthTokenHandler) Subscribe(nc *nats.Conn) error {
	_, err := nc.Subscribe("auth.token.*", func(msg *nats.Msg) {
		var m AuthTokenMessage
		if err := json.Unmarshal(msg.Data, &m); err != nil {
			log.Printf("nats auth.token unmarshal error: %v", err)
			return
		}

		payload, _ := json.Marshal(map[string]string{
			"type":    "auth",
			"token":   m.Token,
			"user_id": m.UserID,
		})

		h.hub.SendToClient(context.Background(), m.Key, nil, payload)
	})
	return err
}
