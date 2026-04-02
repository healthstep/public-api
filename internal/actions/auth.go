package actions

import (
	"context"

	userspb "github.com/helthtech/core-users/pkg/proto/users"
	"github.com/helthtech/public-api/inte
	"github.com/porebric/resty/responses"
)

type AuthController struct {
	usersClient userspb.UserServiceClient
	tgBotURL    string
	maxBotURL   string
}

func NewAuthController(usersClient userspb.UserServiceClient, tgBotURL, maxBotURL string) *AuthController {
	return &AuthController{usersClient: usersClient, tgBotURL: tgBotURL, maxBotURL: maxBotURL}
}

type BrowserChallengeResponse struct {
	Key       string `json:"key"`
	TgBotURL  string `json:"tg_bot_url"`
	MaxBotURL string `json:"max_bot_url"`
}

func (c *AuthController) BrowserChallenge(ctx context.Context, _ requests.BrowserChallengeRequest) (responses.Response, int) {
	resp, err := c.usersClient.PrepareAuth(ctx, &userspb.PrepareAuthRequest{})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to prepare auth"}, 500
	}
	return successData(BrowserChallengeResponse{
		Key:       resp.Key,
		TgBotURL:  c.tgBotURL + "?start=" + resp.Key,
		MaxBotURL: c.maxBotURL + "?start=" + resp.Key,
	}), 200
}
