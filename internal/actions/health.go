package actions

import (
	"context"

	healthpb "github.com/helthtech/core-health/pkg/proto/health"
	"github.com/helthtech/public-api/internal/middleware"
	"github.com/helthtech/public-api/internal/requests"
	"github.com/porebric/resty/responses"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HealthController struct {
	healthClient healthpb.HealthServiceClient
}

func NewHealthController(healthClient healthpb.HealthServiceClient) *HealthController {
	return &HealthController{healthClient: healthClient}
}

func (c *HealthController) ListCriteria(ctx context.Context, _ requests.ListCriteriaRequest) (responses.Response, int) {
	resp, err := c.healthClient.ListCriteria(ctx, &healthpb.ListCriteriaRequest{})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to list criteria"}, 500
	}
	return successData(resp.Criteria), 200
}

func (c *HealthController) ListLabTests(ctx context.Context, _ requests.ListLabTestsRequest) (responses.Response, int) {
	resp, err := c.healthClient.ListLabTests(ctx, &healthpb.ListLabTestsRequest{})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to list lab tests"}, 500
	}
	return successData(resp.LabTests), 200
}

func (c *HealthController) CreateNumericEvent(ctx context.Context, req requests.CreateNumericEventRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	pbReq := &healthpb.CreateNumericEventRequest{
		UserId:            userID,
		HealthCriterionId: req.HealthCriterionID,
		LabTestId:         req.LabTestID,
		NumericValue:      req.NumericValue,
		OccurredAt:        timestamppb.Now(),
		Source:            "web",
		Note:              req.Note,
	}
	resp, err := c.healthClient.CreateNumericEvent(ctx, pbReq)
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to create event"}, 500
	}
	return successData(resp), 200
}

func (c *HealthController) CreateBooleanEvent(ctx context.Context, req requests.CreateBooleanEventRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	pbReq := &healthpb.CreateBooleanEventRequest{
		UserId:            userID,
		HealthCriterionId: req.HealthCriterionID,
		BooleanValue:      req.BooleanValue,
		OccurredAt:        timestamppb.Now(),
		Source:            "web",
		Note:              req.Note,
	}
	resp, err := c.healthClient.CreateBooleanEvent(ctx, pbReq)
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to create event"}, 500
	}
	return successData(resp), 200
}

func (c *HealthController) CreateMarkDoneEvent(ctx context.Context, req requests.CreateMarkDoneEventRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	pbReq := &healthpb.CreateMarkDoneEventRequest{
		UserId:            userID,
		HealthCriterionId: req.HealthCriterionID,
		OccurredAt:        timestamppb.Now(),
		Source:            "web",
		Note:              req.Note,
	}
	resp, err := c.healthClient.CreateMarkDoneEvent(ctx, pbReq)
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to create event"}, 500
	}
	return successData(resp), 200
}

func (c *HealthController) GetDashboard(ctx context.Context, _ requests.GetDashboardRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.healthClient.GetDashboard(ctx, &healthpb.GetDashboardRequest{UserId: userID})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to get dashboard"}, 500
	}
	return successData(resp), 200
}
