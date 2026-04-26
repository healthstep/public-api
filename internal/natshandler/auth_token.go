package natshandler

import (
	"context"
	"encoding/json"

	"github.com/helthtech/public-api/internal/obs"
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
			obs.BG("nats").Error(err, "nats auth.token unmarshal", "subject", msg.Subject)
			return
		}

		payload, _ := json.Marshal(map[string]string{
			"type":    "auth",
			"token":   m.Token,
			"user_id": m.UserID,
		})

		h.hub.SendToClient(context.Background(), m.Key, nil, payload)
		obs.BG("nats").Info("nats auth.token sent to ws", "key", m.Key, "user_id", m.UserID)
	})
	return err
}
