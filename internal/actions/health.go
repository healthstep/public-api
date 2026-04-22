package actions

import (
	"context"

	healthpb "github.com/helthtech/core-health/pkg/proto/health"
	userspb "github.com/helthtech/core-users/pkg/proto/users"
	"github.com/helthtech/public-api/internal/middleware"
	"github.com/helthtech/public-api/internal/requests"
	"github.com/porebric/resty/responses"
)

type HealthController struct {
	healthClient healthpb.HealthServiceClient
	usersClient  userspb.UserServiceClient
}

func NewHealthController(healthClient healthpb.HealthServiceClient, usersClient userspb.UserServiceClient) *HealthController {
	return &HealthController{healthClient: healthClient, usersClient: usersClient}
}

func (c *HealthController) ListGroups(ctx context.Context, _ requests.ListGroupsRequest) (responses.Response, int) {
	resp, err := c.healthClient.ListGroups(ctx, &healthpb.ListGroupsRequest{})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to list groups"}, 500
	}
	return successData(resp.Groups), 200
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
		MeasuredAt:  req.MeasuredAt,
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

func (c *HealthController) GetWeeklyRecommendations(ctx context.Context, _ requests.GetWeeklyRecommendationsRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.healthClient.GetWeeklyRecommendations(ctx, &healthpb.GetWeeklyRecommendationsRequest{UserId: userID})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to get weekly recommendations"}, 500
	}
	return successData(resp), 200
}

// --- Admin ---

func (c *HealthController) requireAdmin(ctx context.Context) bool {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return false
	}
	resp, err := c.usersClient.GetUser(ctx, &userspb.GetUserRequest{UserId: userID})
	if err != nil {
		return false
	}
	return resp.GetIsAdmin()
}

func (c *HealthController) AdminListRecommendations(ctx context.Context, req requests.AdminListRecommendationsRequest) (responses.Response, int) {
	if !c.requireAdmin(ctx) {
		return &responses.ErrorResponse{Message: "forbidden"}, 403
	}
	resp, err := c.healthClient.AdminListRecommendations(ctx, &healthpb.AdminListRecommendationsRequest{
		CriterionId: req.CriterionID,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to list recommendations"}, 500
	}
	return successData(resp.Recommendations), 200
}

func (c *HealthController) AdminUpsertRecommendation(ctx context.Context, req requests.AdminUpsertRecommendationRequest) (responses.Response, int) {
	if !c.requireAdmin(ctx) {
		return &responses.ErrorResponse{Message: "forbidden"}, 403
	}
	pr := &healthpb.AdminRecommendation{
		Id:          req.ID,
		CriterionId: req.CriterionID,
		Type:        req.Type,
		Title:       req.Title,
		Texts:       req.Texts,
		BaseWeight:  req.BaseWeight,
	}
	resp, err := c.healthClient.AdminUpsertRecommendation(ctx, &healthpb.AdminUpsertRecommendationRequest{
		Recommendation: pr,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to upsert recommendation"}, 500
	}
	return successData(resp.Recommendation), 200
}

func (c *HealthController) AdminDeleteRecommendation(ctx context.Context, req requests.AdminDeleteRecommendationRequest) (responses.Response, int) {
	if !c.requireAdmin(ctx) {
		return &responses.ErrorResponse{Message: "forbidden"}, 403
	}
	resp, err := c.healthClient.AdminDeleteRecommendation(ctx, &healthpb.AdminDeleteRecommendationRequest{
		Id: req.ID,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to delete recommendation"}, 500
	}
	return successData(resp), 200
}

func (c *HealthController) AdminUpsertCriterion(ctx context.Context, req requests.AdminUpsertCriterionRequest) (responses.Response, int) {
	if !c.requireAdmin(ctx) {
		return &responses.ErrorResponse{Message: "forbidden"}, 403
	}
	pc := &healthpb.Criterion{
		Id:          req.ID,
		GroupId:     req.GroupID,
		Name:        req.Name,
		Level:       req.Level,
		Sex:         req.Sex,
		InputType:   req.InputType,
		Lifetime:    req.Lifetime,
		SortOrder:   req.SortOrder,
		MinValue:    req.MinValue,
		MaxValue:    req.MaxValue,
		Delta:       req.Delta,
		Instruction: req.Instruction,
	}
	resp, err := c.healthClient.AdminUpsertCriterion(ctx, &healthpb.AdminUpsertCriterionRequest{
		Criterion: pc,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to upsert criterion"}, 500
	}
	return successData(resp.Criterion), 200
}
