package actions

import (
	"context"

	healthpb "github.com/helthtech/core-health/pkg/proto/health"
	"github.com/helthtech/public-api/internal/middleware"
	"github.com/helthtech/public-api/internal/requests"
	"github.com/porebric/resty/responses"
)

type HealthController struct {
	healthClient healthpb.HealthServiceClient
}

func NewHealthController(healthClient healthpb.HealthServiceClient) *HealthController {
	return &HealthController{healthClient: healthClient}
}

func (c *HealthController) ListCriteria(ctx context.Context, _ requests.ListCriteriaRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	resp, err := c.healthClient.ListCriteria(ctx, &healthpb.ListCriteriaRequest{
		UserId: userID,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to list criteria"}, 500
	}
	return successData(resp.Criteria), 200
}

func (c *HealthController) SetUserCriterion(ctx context.Context, req requests.SetUserCriterionRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.healthClient.SetUserCriterion(ctx, &healthpb.SetUserCriterionRequest{
		UserId:      userID,
		CriterionId: req.CriterionID,
		Value:       req.Value,
		Source:      "web",
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to set criterion"}, 500
	}
	return successData(resp), 200
}

func (c *HealthController) GetUserCriteria(ctx context.Context, _ requests.GetUserCriteriaRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.healthClient.GetUserCriteria(ctx, &healthpb.GetUserCriteriaRequest{
		UserId: userID,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to get user criteria"}, 500
	}
	return successData(resp.Entries), 200
}

func (c *HealthController) ResetCriteria(ctx context.Context, _ requests.ResetCriteriaRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.healthClient.ResetCriteria(ctx, &healthpb.ResetCriteriaRequest{
		UserId: userID,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to reset criteria"}, 500
	}
	return successData(resp), 200
}

func (c *HealthController) GetProgress(ctx context.Context, _ requests.GetProgressRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.healthClient.GetProgress(ctx, &healthpb.GetProgressRequest{UserId: userID})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to get progress"}, 500
	}
	return successData(resp), 200
}

func (c *HealthController) GetRecommendations(ctx context.Context, _ requests.GetRecommendationsRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.healthClient.GetRecommendations(ctx, &healthpb.GetRecommendationsRequest{UserId: userID})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to get recommendations"}, 500
	}
	return successData(resp.Recommendations), 200
}
