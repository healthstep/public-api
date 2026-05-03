package requests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

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

type CheckAuthKeyRequest struct {
	Key string
}

func (CheckAuthKeyRequest) Validate() (bool, string, string) { return true, "", "" }
func (CheckAuthKeyRequest) Methods() []string                { return []string{"GET"} }
func (CheckAuthKeyRequest) Path() (string, bool)             { return "/api/v1/auth/check", false }
func (CheckAuthKeyRequest) String() string                   { return "check-auth-key" }

func NewCheckAuthKeyRequest(ctx context.Context, r *http.Request) (context.Context, CheckAuthKeyRequest, error) {
	return ctx, CheckAuthKeyRequest{Key: r.URL.Query().Get("key")}, nil
}

type LoginWithPasswordRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (LoginWithPasswordRequest) Validate() (bool, string, string) { return true, "", "" }
func (LoginWithPasswordRequest) Methods() []string                { return []string{"POST"} }
func (LoginWithPasswordRequest) Path() (string, bool)             { return "/api/v1/auth/login", false }
func (LoginWithPasswordRequest) String() string                   { return "login-with-password" }

func NewLoginWithPasswordRequest(ctx context.Context, r *http.Request) (context.Context, LoginWithPasswordRequest, error) {
	var req LoginWithPasswordRequest
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
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

// --- Criteria ---

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

// --- Reset all criteria ---

type ResetCriteriaRequest struct {
	AuthenticatedRequest
}

func (ResetCriteriaRequest) Validate() (bool, string, string) { return true, "", "" }
func (ResetCriteriaRequest) Methods() []string                { return []string{"DELETE"} }
func (ResetCriteriaRequest) Path() (string, bool) {
	return "/api/v1/health/user-criteria", false
}
func (ResetCriteriaRequest) String() string { return "reset-criteria" }

func NewResetCriteriaRequest(ctx context.Context, r *http.Request) (context.Context, ResetCriteriaRequest, error) {
	return ctx, ResetCriteriaRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

// --- User Criteria ---

type SetUserCriterionRequest struct {
	AuthenticatedRequest
	CriterionID string `json:"criterion_id"`
	Value         string `json:"value"`
	MeasuredAt    string `json:"measured_at"` // optional ISO date "2006-01-02"
	UserSex       string `json:"user_sex"`  // optional: "male" | "female" | ""; used for entry in response (same as GetUserCriteria)
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
	UserSex string // query: user_sex (optional)
}

func (GetUserCriteriaRequest) Validate() (bool, string, string) { return true, "", "" }
func (GetUserCriteriaRequest) Methods() []string                { return []string{"GET"} }
func (GetUserCriteriaRequest) Path() (string, bool)             { return "/api/v1/health/user-criteria", false }
func (GetUserCriteriaRequest) String() string                   { return "get-user-criteria" }

func NewGetUserCriteriaRequest(ctx context.Context, r *http.Request) (context.Context, GetUserCriteriaRequest, error) {
	return ctx, GetUserCriteriaRequest{
		AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)},
		UserSex:                r.URL.Query().Get("user_sex"),
	}, nil
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

// --- Weekly Recommendations ---

type GetWeeklyRecommendationsRequest struct {
	AuthenticatedRequest
}

func (GetWeeklyRecommendationsRequest) Validate() (bool, string, string) { return true, "", "" }
func (GetWeeklyRecommendationsRequest) Methods() []string                { return []string{"GET"} }
func (GetWeeklyRecommendationsRequest) Path() (string, bool) {
	return "/api/v1/health/weekly-recommendations", false
}
func (GetWeeklyRecommendationsRequest) String() string { return "get-weekly-recommendations" }

func NewGetWeeklyRecommendationsRequest(ctx context.Context, r *http.Request) (context.Context, GetWeeklyRecommendationsRequest, error) {
	return ctx, GetWeeklyRecommendationsRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

// --- Groups ---

type ListGroupsRequest struct {
	AuthenticatedRequest
}

func (ListGroupsRequest) Validate() (bool, string, string) { return true, "", "" }
func (ListGroupsRequest) Methods() []string                { return []string{"GET"} }
func (ListGroupsRequest) Path() (string, bool)             { return "/api/v1/health/groups", false }
func (ListGroupsRequest) String() string                   { return "list-groups" }

func NewListGroupsRequest(ctx context.Context, r *http.Request) (context.Context, ListGroupsRequest, error) {
	return ctx, ListGroupsRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}, nil
}

// --- Admin Recommendations ---

type AdminListRecommendationsRequest struct {
	AuthenticatedRequest
	CriterionID string
}

func (AdminListRecommendationsRequest) Validate() (bool, string, string) { return true, "", "" }
func (AdminListRecommendationsRequest) Methods() []string                { return []string{"GET"} }
func (AdminListRecommendationsRequest) Path() (string, bool) {
	return "/api/v1/admin/recommendations", false
}
func (AdminListRecommendationsRequest) String() string { return "admin-list-recommendations" }

func NewAdminListRecommendationsRequest(ctx context.Context, r *http.Request) (context.Context, AdminListRecommendationsRequest, error) {
	req := AdminListRecommendationsRequest{
		AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)},
		CriterionID:          r.URL.Query().Get("criterion_id"),
	}
	return ctx, req, nil
}

type AdminUpsertRecommendationRequest struct {
	AuthenticatedRequest
	ID          string   `json:"id"`
	CriterionID string   `json:"criterion_id"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Texts       []string `json:"texts"`
	BaseWeight  int32    `json:"base_weight"`
}

func (AdminUpsertRecommendationRequest) Validate() (bool, string, string) { return true, "", "" }
func (AdminUpsertRecommendationRequest) Methods() []string                { return []string{"POST"} }
func (AdminUpsertRecommendationRequest) Path() (string, bool) {
	return "/api/v1/admin/recommendations", false
}
func (AdminUpsertRecommendationRequest) String() string { return "admin-upsert-recommendation" }

func NewAdminUpsertRecommendationRequest(ctx context.Context, r *http.Request) (context.Context, AdminUpsertRecommendationRequest, error) {
	var req AdminUpsertRecommendationRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

type AdminDeleteRecommendationRequest struct {
	AuthenticatedRequest
	ID string
}

func (AdminDeleteRecommendationRequest) Validate() (bool, string, string) { return true, "", "" }
func (AdminDeleteRecommendationRequest) Methods() []string                { return []string{"DELETE"} }
func (AdminDeleteRecommendationRequest) Path() (string, bool) {
	return "/api/v1/admin/recommendations/{id}", true
}
func (AdminDeleteRecommendationRequest) String() string { return "admin-delete-recommendation" }

func NewAdminDeleteRecommendationRequest(ctx context.Context, r *http.Request) (context.Context, AdminDeleteRecommendationRequest, error) {
	req := AdminDeleteRecommendationRequest{
		AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)},
		ID:                   r.PathValue("id"),
	}
	return ctx, req, nil
}

// --- Admin Criteria ---

type AdminUpsertCriterionRequest struct {
	AuthenticatedRequest
	ID         string   `json:"id"`
	GroupID    string   `json:"group_id"`
	AnalysisID *int64   `json:"analysis_id,omitempty"`
	Name       string   `json:"name"`
	Level      int32    `json:"level"`
	Sex        string   `json:"sex"`
	InputType  string   `json:"input_type"`
	Lifetime   int32    `json:"lifetime"`
	SortOrder  int32    `json:"sort_order"`
	MinValue   *float64 `json:"min_value,omitempty"`
	MaxValue   *float64 `json:"max_value,omitempty"`
	Delta      *float64 `json:"delta,omitempty"`
}

func (AdminUpsertCriterionRequest) Validate() (bool, string, string) { return true, "", "" }
func (AdminUpsertCriterionRequest) Methods() []string                { return []string{"POST"} }
func (AdminUpsertCriterionRequest) Path() (string, bool)             { return "/api/v1/admin/criteria", false }
func (AdminUpsertCriterionRequest) String() string                   { return "admin-upsert-criterion" }

func NewAdminUpsertCriterionRequest(ctx context.Context, r *http.Request) (context.Context, AdminUpsertCriterionRequest, error) {
	var req AdminUpsertCriterionRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

// --- Admin analyses ---

type AdminListAnalysesRequest struct {
	AuthenticatedRequest
}

func (AdminListAnalysesRequest) Validate() (bool, string, string) { return true, "", "" }
func (AdminListAnalysesRequest) Methods() []string                { return []string{"GET"} }
func (AdminListAnalysesRequest) Path() (string, bool)             { return "/api/v1/admin/analyses", false }
func (AdminListAnalysesRequest) String() string                   { return "admin-list-analyses" }

func NewAdminListAnalysesRequest(ctx context.Context, r *http.Request) (context.Context, AdminListAnalysesRequest, error) {
	req := AdminListAnalysesRequest{AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)}}
	return ctx, req, nil
}

type AdminUpsertAnalysisRequest struct {
	AuthenticatedRequest
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Instruction string `json:"instruction"`
}

func (AdminUpsertAnalysisRequest) Validate() (bool, string, string) { return true, "", "" }
func (AdminUpsertAnalysisRequest) Methods() []string                { return []string{"POST"} }
func (AdminUpsertAnalysisRequest) Path() (string, bool)             { return "/api/v1/admin/analyses", false }
func (AdminUpsertAnalysisRequest) String() string                   { return "admin-upsert-analysis" }

func NewAdminUpsertAnalysisRequest(ctx context.Context, r *http.Request) (context.Context, AdminUpsertAnalysisRequest, error) {
	var req AdminUpsertAnalysisRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}

type AdminDeleteAnalysisRequest struct {
	AuthenticatedRequest
	ID int64
}

func (AdminDeleteAnalysisRequest) Validate() (bool, string, string) { return true, "", "" }
func (AdminDeleteAnalysisRequest) Methods() []string                { return []string{"DELETE"} }
func (AdminDeleteAnalysisRequest) Path() (string, bool)             { return "/api/v1/admin/analyses/{id}", true }
func (AdminDeleteAnalysisRequest) String() string                   { return "admin-delete-analysis" }

func NewAdminDeleteAnalysisRequest(ctx context.Context, r *http.Request) (context.Context, AdminDeleteAnalysisRequest, error) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return ctx, AdminDeleteAnalysisRequest{
		AuthenticatedRequest: AuthenticatedRequest{Token: middleware.ExtractBearerToken(r)},
		ID:                   id,
	}, nil
}

// --- Lab import (confirm AI-extracted values, staged in Redis) ---

type ConfirmLabImportRequest struct {
	AuthenticatedRequest
	PendingID string `json:"pending_id"`
	Accept    bool   `json:"accept"`
	UserSex   string `json:"user_sex"`
}

func (ConfirmLabImportRequest) Validate() (bool, string, string) { return true, "", "" }
func (ConfirmLabImportRequest) Methods() []string                { return []string{"POST"} }
func (ConfirmLabImportRequest) Path() (string, bool)             { return "/api/v1/health/lab-import/confirm", false }
func (ConfirmLabImportRequest) String() string                   { return "confirm-lab-import" }

func NewConfirmLabImportRequest(ctx context.Context, r *http.Request) (context.Context, ConfirmLabImportRequest, error) {
	var req ConfirmLabImportRequest
	req.Token = middleware.ExtractBearerToken(r)
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &req)
	return ctx, req, nil
}
