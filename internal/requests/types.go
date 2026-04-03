package requests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/helthtech/public-api/internal/middleware"
)

type BrowserChallengeRequest struct{}

func (BrowserChallengeRequest) Validate() (bool, string, string) { return true, "", "" }
func (BrowserChallengeRequest) Methods() []string                { return []string{"POST"} }
func (BrowserChallengeRequest) Path() (string, bool)             { return "/api/v1/auth/browser-challenge", false }
func (BrowserChallengeRequest) String() string                   { return "browser-challenge" }

func NewBrowserChallengeRequest(ctx context.Context, _ *http.Request) (context.Context, BrowserChallengeRequest, error) {
	return ctx, BrowserChallengeRequest{}, nil
}

type AuthenticatedRequest struct {
	Token  string
	UserID string
}

func (r AuthenticatedRequest) GetAuthToken() string { return r.Token }

type GetMeRequest struct {
	AuthenticatedRequest
}

func (GetMeRequest) Validate() (bool, string, string) { return true, "", "" }
func (GetMeRequest) Methods() []string                { return []string{"GET"} }
func (GetMeRequest) Path() (string, bool)             { return "/api/v1/users/me", false }
func (GetMeRequest) String() string                   { return "get-me" }

func NewGetMeRequest(ctx context.Context, r *http.Request) (context.Context, GetMeRequest, error) {
	return ctx, GetMeRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

type UpdateMeRequest struct {
	AuthenticatedRequest
	DisplayName *string `json:"display_name,omitempty"`
	Locale      *string `json:"locale,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	BirthDate   *string `json:"birth_date,omitempty"`
	Sex         *string `json:"sex,omitempty"`
}

func (UpdateMeRequest) Validate() (bool, string, string) { return true, "", "" }
func (UpdateMeRequest) Methods() []string                { return []string{"PATCH"} }
func (UpdateMeRequest) Path() (string, bool)             { return "/api/v1/users/me", false }
func (UpdateMeRequest) String() string                   { return "update-me" }

func NewUpdateMeRequest(ctx context.Context, r *http.Request) (context.Context, UpdateMeRequest, error) {
	var req UpdateMeRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

type ListCriteriaRequest struct {
	AuthenticatedRequest
}

func (ListCriteriaRequest) Validate() (bool, string, string) { return true, "", "" }
func (ListCriteriaRequest) Methods() []string                { return []string{"GET"} }
func (ListCriteriaRequest) Path() (string, bool)             { return "/api/v1/health/criteria", false }
func (ListCriteriaRequest) String() string                   { return "list-criteria" }

func NewListCriteriaRequest(ctx context.Context, r *http.Request) (context.Context, ListCriteriaRequest, error) {
	return ctx, ListCriteriaRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

type ListLabTestsRequest struct {
	AuthenticatedRequest
}

func (ListLabTestsRequest) Validate() (bool, string, string) { return true, "", "" }
func (ListLabTestsRequest) Methods() []string                { return []string{"GET"} }
func (ListLabTestsRequest) Path() (string, bool)             { return "/api/v1/health/lab-tests", false }
func (ListLabTestsRequest) String() string                   { return "list-lab-tests" }

func NewListLabTestsRequest(ctx context.Context, r *http.Request) (context.Context, ListLabTestsRequest, error) {
	return ctx, ListLabTestsRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

type CreateNumericEventRequest struct {
	AuthenticatedRequest
	HealthCriterionID string  `json:"health_criterion_id"`
	LabTestID         string  `json:"lab_test_id,omitempty"`
	NumericValue      float64 `json:"numeric_value"`
	OccurredAt        string  `json:"occurred_at,omitempty"`
	Note              string  `json:"note,omitempty"`
}

func (CreateNumericEventRequest) Validate() (bool, string, string) { return true, "", "" }
func (CreateNumericEventRequest) Methods() []string                { return []string{"POST"} }
func (CreateNumericEventRequest) Path() (string, bool) {
	return "/api/v1/health/events/numeric", false
}
func (CreateNumericEventRequest) String() string { return "create-numeric-event" }

func NewCreateNumericEventRequest(ctx context.Context, r *http.Request) (context.Context, CreateNumericEventRequest, error) {
	var req CreateNumericEventRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

type CreateBooleanEventRequest struct {
	AuthenticatedRequest
	HealthCriterionID string `json:"health_criterion_id"`
	BooleanValue      string `json:"boolean_value"`
	OccurredAt        string `json:"occurred_at,omitempty"`
	Note              string `json:"note,omitempty"`
}

func (CreateBooleanEventRequest) Validate() (bool, string, string) { return true, "", "" }
func (CreateBooleanEventRequest) Methods() []string                { return []string{"POST"} }
func (CreateBooleanEventRequest) Path() (string, bool) {
	return "/api/v1/health/events/boolean", false
}
func (CreateBooleanEventRequest) String() string { return "create-boolean-event" }

func NewCreateBooleanEventRequest(ctx context.Context, r *http.Request) (context.Context, CreateBooleanEventRequest, error) {
	var req CreateBooleanEventRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

type CreateMarkDoneEventRequest struct {
	AuthenticatedRequest
	HealthCriterionID string `json:"health_criterion_id"`
	OccurredAt        string `json:"occurred_at,omitempty"`
	Note              string `json:"note,omitempty"`
}

func (CreateMarkDoneEventRequest) Validate() (bool, string, string) { return true, "", "" }
func (CreateMarkDoneEventRequest) Methods() []string                { return []string{"POST"} }
func (CreateMarkDoneEventRequest) Path() (string, bool) {
	return "/api/v1/health/events/mark-done", false
}
func (CreateMarkDoneEventRequest) String() string { return "create-mark-done-event" }

func NewCreateMarkDoneEventRequest(ctx context.Context, r *http.Request) (context.Context, CreateMarkDoneEventRequest, error) {
	var req CreateMarkDoneEventRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

type GetDashboardRequest struct {
	AuthenticatedRequest
}

func (GetDashboardRequest) Validate() (bool, string, string) { return true, "", "" }
func (GetDashboardRequest) Methods() []string                { return []string{"GET"} }
func (GetDashboardRequest) Path() (string, bool)             { return "/api/v1/health/dashboard", false }
func (GetDashboardRequest) String() string                   { return "get-dashboard" }

func NewGetDashboardRequest(ctx context.Context, r *http.Request) (context.Context, GetDashboardRequest, error) {
	return ctx, GetDashboardRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}
