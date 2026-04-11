package actions

import (
	"context"

	userspb "github.com/helthtech/core-users/pkg/proto/users"
	"github.com/helthtech/public-api/internal/requests"
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

type CheckAuthKeyResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func (c *AuthController) CheckAuthKey(ctx context.Context, req requests.CheckAuthKeyRequest) (responses.Response, int) {
	if req.Key == "" {
		return &responses.ErrorResponse{Message: "key required"}, 400
	}
	resp, err := c.usersClient.CheckAuthToken(ctx, &userspb.CheckAuthTokenRequest{Key: req.Key})
	if err != nil {
		return &responses.ErrorResponse{Message: "check failed"}, 500
	}
	if resp.Token == "" {
		return successData(CheckAuthKeyResponse{}), 200
	}
	return successData(CheckAuthKeyResponse{Token: resp.Token, UserID: resp.UserId}), 200
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func (c *AuthController) LoginWithPassword(ctx context.Context, req requests.LoginWithPasswordRequest) (responses.Response, int) {
	if req.Phone == "" || req.Password == "" {
		return &responses.ErrorResponse{Message: "phone and password required"}, 400
	}
	resp, err := c.usersClient.LoginWithPassword(ctx, &userspb.LoginWithPasswordRequest{
		PhoneE164: req.Phone,
		Password:  req.Password,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "invalid credentials"}, 401
	}
	return successData(LoginResponse{Token: resp.Token, UserID: resp.UserId}), 200
}
