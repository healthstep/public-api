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

// --- Analysis ---

type ListAnalysisRequest struct {
	AuthenticatedRequest
}

func (ListAnalysisRequest) Validate() (bool, string, string) { return true, "", "" }
func (ListAnalysisRequest) Methods() []string                { return []string{"GET"} }
func (ListAnalysisRequest) Path() (string, bool)             { return "/api/v1/health/analysis", false }
func (ListAnalysisRequest) String() string                   { return "list-analysis" }

func NewListAnalysisRequest(ctx context.Context, r *http.Request) (context.Context, ListAnalysisRequest, error) {
	return ctx, ListAnalysisRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

// --- Reset analysis criteria ---

type ResetAnalysisCriteriaRequest struct {
	AuthenticatedRequest
	AnalysisID string `json:"analysis_id"`
}

func (ResetAnalysisCriteriaRequest) Validate() (bool, string, string) { return true, "", "" }
func (ResetAnalysisCriteriaRequest) Methods() []string                { return []string{"DELETE"} }
func (ResetAnalysisCriteriaRequest) Path() (string, bool) {
	return "/api/v1/health/user-criteria/reset", false
}
func (ResetAnalysisCriteriaRequest) String() string { return "reset-analysis-criteria" }

func NewResetAnalysisCriteriaRequest(ctx context.Context, r *http.Request) (context.Context, ResetAnalysisCriteriaRequest, error) {
	var req ResetAnalysisCriteriaRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

// --- Criteria ---

type ListCriteriaRequest struct {
	AuthenticatedRequest
	AnalysisID string
}

func (ListCriteriaRequest) Validate() (bool, string, string) { return true, "", "" }
func (ListCriteriaRequest) Methods() []string                { return []string{"GET"} }
func (ListCriteriaRequest) Path() (string, bool)             { return "/api/v1/health/criteria", false }
func (ListCriteriaRequest) String() string                   { return "list-criteria" }

func NewListCriteriaRequest(ctx context.Context, r *http.Request) (context.Context, ListCriteriaRequest, error) {
	req := ListCriteriaRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}
	req.AnalysisID = r.URL.Query().Get("analysis_id")
	return ctx, req, nil
}

// --- User Criteria ---

type SetUserCriterionRequest struct {
	AuthenticatedRequest
	CriterionID string `json:"criterion_id"`
	Value       string `json:"value"`
}

func (SetUserCriterionRequest) Validate() (bool, string, string) { return true, "", "" }
func (SetUserCriterionRequest) Methods() []string                { return []string{"POST"} }
func (SetUserCriterionRequest) Path() (string, bool)             { return "/api/v1/health/user-criteria", false }
func (SetUserCriterionRequest) String() string                   { return "set-user-criterion" }

func NewSetUserCriterionRequest(ctx context.Context, r *http.Request) (context.Context, SetUserCriterionRequest, error) {
	var req SetUserCriterionRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

type GetUserCriteriaRequest struct {
	AuthenticatedRequest
}

func (GetUserCriteriaRequest) Validate() (bool, string, string) { return true, "", "" }
func (GetUserCriteriaRequest) Methods() []string                { return []string{"GET"} }
func (GetUserCriteriaRequest) Path() (string, bool)             { return "/api/v1/health/user-criteria", false }
func (GetUserCriteriaRequest) String() string                   { return "get-user-criteria" }

func NewGetUserCriteriaRequest(ctx context.Context, r *http.Request) (context.Context, GetUserCriteriaRequest, error) {
	return ctx, GetUserCriteriaRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

// --- Progress & Recommendations ---

type GetProgressRequest struct {
	AuthenticatedRequest
}

func (GetProgressRequest) Validate() (bool, string, string) { return true, "", "" }
func (GetProgressRequest) Methods() []string                { return []string{"GET"} }
func (GetProgressRequest) Path() (string, bool)             { return "/api/v1/health/progress", false }
func (GetProgressRequest) String() string                   { return "get-progress" }

func NewGetProgressRequest(ctx context.Context, r *http.Request) (context.Context, GetProgressRequest, error) {
	return ctx, GetProgressRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

type GetRecommendationsRequest struct {
	AuthenticatedRequest
}

func (GetRecommendationsRequest) Validate() (bool, string, string) { return true, "", "" }
func (GetRecommendationsRequest) Methods() []string                { return []string{"GET"} }
func (GetRecommendationsRequest) Path() (string, bool) {
	return "/api/v1/health/recommendations", false
}
func (GetRecommendationsRequest) String() string { return "get-recommendations" }

func NewGetRecommendationsRequest(ctx context.Context, r *http.Request) (context.Context, GetRecommendationsRequest, error) {
	return ctx, GetRecommendationsRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}
