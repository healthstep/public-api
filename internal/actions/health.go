package actions

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	healthpb "github.com/helthtech/core-health/pkg/proto/health"
	userspb "github.com/helthtech/core-users/pkg/proto/users"
	"github.com/helthtech/public-api/internal/middleware"
	"github.com/helthtech/public-api/internal/requests"
	"github.com/porebric/resty/responses"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthController struct {
	healthClient  healthpb.HealthServiceClient
	usersClient   userspb.UserServiceClient
	jwtSecret     string
	maxImportDocs int
}

func NewHealthController(healthClient healthpb.HealthServiceClient, usersClient userspb.UserServiceClient, jwtSecret string, maxImportDocs int) *HealthController {
	if maxImportDocs <= 0 {
		maxImportDocs = 5
	}
	return &HealthController{healthClient: healthClient, usersClient: usersClient, jwtSecret: jwtSecret, maxImportDocs: maxImportDocs}
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
		UserId:        userID,
		CriterionId:   req.CriterionID,
		Value:         req.Value,
		Source:        "web",
		MeasuredAt:    req.MeasuredAt,
		UserSex:       req.UserSex,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: "failed to set criterion"}, 500
	}
	return successData(resp), 200
}

func (c *HealthController) GetUserCriteria(ctx context.Context, req requests.GetUserCriteriaRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	resp, err := c.healthClient.GetUserCriteria(ctx, &healthpb.GetUserCriteriaRequest{
		UserId:   userID,
		UserSex:  req.UserSex,
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

// ConfirmLabImport applies or discards a pending import from the lab document flow.
func (c *HealthController) ConfirmLabImport(ctx context.Context, req requests.ConfirmLabImportRequest) (responses.Response, int) {
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		return &responses.ErrorResponse{Message: "unauthorized"}, 401
	}
	if req.PendingID == "" {
		return &responses.ErrorResponse{Message: "pending_id required"}, 400
	}
	out, err := c.healthClient.ConfirmPendingImport(ctx, &healthpb.ConfirmPendingImportRequest{
		UserId:   userID,
		PendingId: req.PendingID,
		Accept:   req.Accept,
		UserSex:  req.UserSex,
	})
	if err != nil {
		return &responses.ErrorResponse{Message: err.Error()}, 500
	}
	if !out.GetSuccess() {
		if out.GetErrorMessage() == "forbidden" {
			return &responses.ErrorResponse{Message: "forbidden"}, 403
		}
		return &responses.ErrorResponse{Message: out.GetErrorMessage()}, 400
	}
	return successData(out), 200
}

// UploadLabDocuments is an HTTP handler (multipart) — not via resty.Endpoint.
func (c *HealthController) UploadLabDocuments(w http.ResponseWriter, r *http.Request) {
	uid, err := middleware.UserIDFromHTTP(r, c.jwtSecret)
	if err != nil {
		http.Error(w, `{"success":false,"message":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, `{"success":false,"message":"invalid multipart"}`, http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		files = r.MultipartForm.File["file"]
	}
	if len(files) == 0 {
		http.Error(w, `{"success":false,"message":"no files; use form field files[]"}`, http.StatusBadRequest)
		return
	}
	if len(files) > c.maxImportDocs {
		http.Error(w, `{"success":false,"message":"too many files"}`, http.StatusBadRequest)
		return
	}
	userSex := r.FormValue("user_sex")
	ctx := r.Context()
	stream, err := c.healthClient.ImportCriteriaFromPdf(ctx)
	if err != nil {
		http.Error(w, `{"success":false,"message":"import unavailable"}`, http.StatusBadGateway)
		return
	}
	if err := stream.Send(&healthpb.ImportCriteriaFromPdfRequest{UserId: uid, UserSex: userSex}); err != nil {
		http.Error(w, `{"success":false,"message":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	const chunkSize = 64 * 1024
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			http.Error(w, `{"success":false,"message":"file open"}`, http.StatusInternalServerError)
			return
		}
		if err := stream.Send(&healthpb.ImportCriteriaFromPdfRequest{Filename: fh.Filename}); err != nil {
			_ = f.Close()
			http.Error(w, `{"success":false,"message":"`+err.Error()+`"}`, http.StatusBadGateway)
			return
		}
		buf := make([]byte, chunkSize)
		for {
			n, rerr := f.Read(buf)
			if n > 0 {
				if serr := stream.Send(&healthpb.ImportCriteriaFromPdfRequest{Chunk: buf[:n]}); serr != nil {
					_ = f.Close()
					http.Error(w, `{"success":false,"message":"`+serr.Error()+`"}`, http.StatusBadGateway)
					return
				}
			}
			if rerr == io.EOF {
				break
			}
			if rerr != nil {
				_ = f.Close()
				http.Error(w, `{"success":false,"message":"read file"}`, http.StatusInternalServerError)
				return
			}
		}
		_ = f.Close()
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unavailable {
			http.Error(w, `{"success":false,"message":"lab import not configured (parser/redis)"}`, http.StatusServiceUnavailable)
			return
		}
		http.Error(w, `{"success":false,"message":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&JSONResponse{Success: true, Data: map[string]any{
		"user_criteria":     resp.GetUserCriteria(),
		"pending_import_id": resp.GetPendingImportId(),
		"model_note":        resp.GetModelNote(),
	}})
}

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
