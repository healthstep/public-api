package actions

import (
	"context"

	userspb "github.com/helthtech/core-users/pkg/proto/users"
	"github.com/helthtech/public-api/internal/middleware"
	"github.com/helthtech/public-api/internal/requests"
	"github.com/porebric/resty/responses"
)

type UserController struct {
	usersClient userspb.UserServiceClient
}

func NewUserController(usersClient userspb.UserServiceClient) *UserController {
	return &UserController{usersClient: usersClient}
}

func (c *UserController) GetMe(ctx context.Context, _ requests.GetMeRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.usersClient.GetUser(ctx, &userspb.GetUserRequest{UserId: userID})
	if err != nil {
		return &responses.ErrorResponse{Message: "user not found"}, 404
	}
	return successData(resp), 200
}

func (c *UserController) UpdateMe(ctx context.Context, req requests.UpdateMeRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	updateReq := &userspb.UpdateUserRequest{UserId: userID}
	if req.DisplayName != nil {
		updateReq.DisplayName = req.DisplayName
	}
	if req.Locale != nil {
		updateReq.Locale = req.Locale
	}
	if req.Timezone != nil {
		updateReq.Timezone = req.Timezone
	}
	if req.BirthDate != nil {
		updateReq.BirthDate = req.BirthDate
	}
	if req.Sex != nil {
		updateReq.Sex = req.Sex
	}
	resp, err := c.usersClient.UpdateUser(ctx, updateReq)
	if err != nil {
		return &responses.ErrorResponse{Message: "update failed"}, 500
	}
	return successData(resp), 200
}
