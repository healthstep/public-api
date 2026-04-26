package boot

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	healthpb "github.com/helthtech/core-health/pkg/proto/health"
	userspb "github.com/helthtech/core-users/pkg/proto/users"
	"github.com/helthtech/public-api/internal/actions"
	"github.com/helthtech/public-api/internal/middleware"
	"github.com/helthtech/public-api/internal/natshandler"
	"github.com/helthtech/public-api/internal/obs"
	"github.com/helthtech/public-api/internal/requests"
	"github.com/nats-io/nats.go"
	"github.com/porebric/configs"
	"github.com/porebric/logger"
	"github.com/porebric/resty"
	restyerrors "github.com/porebric/resty/errors"
	"github.com/porebric/resty/ws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(ctx context.Context) error {
	usersConn, err := grpc.NewClient(
		configs.Value(ctx, "grpc_core_users").String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("grpc core-users: %w", err)
	}

	healthConn, err := grpc.NewClient(
		configs.Value(ctx, "grpc_core_health").String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("grpc core-health: %w", err)
	}

	usersClient := userspb.NewUserServiceClient(usersConn)
	healthClient := healthpb.NewHealthServiceClient(healthConn)

	nc, err := nats.Connect(configs.Value(ctx, "nats_url").String())
	if err != nil {
		return fmt.Errorf("nats: %w", err)
	}

	wsHub := ws.NewHub(
		func(_ context.Context, msg *ws.LoginMessage) (context.Context, ws.Error) {
			return nil, ws.Error{}
		},
		func(r *http.Request) string {
			return r.URL.Query().Get("key")
		},
	)

	authHandler := natshandler.NewAuthTokenHandler(wsHub)
	if err = authHandler.Subscribe(nc); err != nil {
		return fmt.Errorf("nats subscribe: %w", err)
	}

	jwtSecret := configs.Value(ctx, "jwt_secret").String()
	jwtMW := middleware.NewJWTAuth(jwtSecret)

	restyerrors.Init(nil)

	router := resty.NewRouter(func() *logger.Logger { return obs.L }, wsHub)
	router.MuxRouter().Use(middleware.AccessLog())

	origins := strings.Split(configs.Value(ctx, "cors_origins").String(), ",")
	methods := strings.Split(configs.Value(ctx, "cors_methods").String(), ",")
	headers := strings.Split(configs.Value(ctx, "cors_headers").String(), ",")
	router.SetCors(origins, methods, headers)

	tgBotURL := configs.Value(ctx, "tg_bot_url").String()
	maxBotURL := configs.Value(ctx, "max_bot_url").String()

	authCtrl := actions.NewAuthController(usersClient, tgBotURL, maxBotURL)
	userCtrl := actions.NewUserController(usersClient)
	maxLab := int(configs.Value(ctx, "max_lab_import_files").Int())
	if maxLab <= 0 {
		maxLab = 5
	}
	healthCtrl := actions.NewHealthController(healthClient, usersClient, jwtSecret, maxLab)

	resty.Endpoint(router, requests.NewBrowserChallengeRequest, authCtrl.BrowserChallenge)
	resty.Endpoint(router, requests.NewCheckAuthKeyRequest, authCtrl.CheckAuthKey)
	resty.Endpoint(router, requests.NewLoginWithPasswordRequest, authCtrl.LoginWithPassword)
	resty.Endpoint(router, requests.NewGetMeRequest, userCtrl.GetMe, jwtMW)
	resty.Endpoint(router, requests.NewUpdateMeRequest, userCtrl.UpdateMe, jwtMW)
	resty.Endpoint(router, requests.NewListGroupsRequest, healthCtrl.ListGroups, jwtMW)
	resty.Endpoint(router, requests.NewListCriteriaRequest, healthCtrl.ListCriteria, jwtMW)
	resty.Endpoint(router, requests.NewSetUserCriterionRequest, healthCtrl.SetUserCriterion, jwtMW)
	resty.Endpoint(router, requests.NewGetUserCriteriaRequest, healthCtrl.GetUserCriteria, jwtMW)
	resty.Endpoint(router, requests.NewGetProgressRequest, healthCtrl.GetProgress, jwtMW)
	resty.Endpoint(router, requests.NewGetRecommendationsRequest, healthCtrl.GetRecommendations, jwtMW)
	resty.Endpoint(router, requests.NewGetWeeklyRecommendationsRequest, healthCtrl.GetWeeklyRecommendations, jwtMW)
	resty.Endpoint(router, requests.NewResetCriteriaRequest, healthCtrl.ResetCriteria, jwtMW)
	resty.Endpoint(router, requests.NewConfirmLabImportRequest, healthCtrl.ConfirmLabImport, jwtMW)
	// Admin endpoints
	resty.Endpoint(router, requests.NewAdminListRecommendationsRequest, healthCtrl.AdminListRecommendations, jwtMW)
	resty.Endpoint(router, requests.NewAdminUpsertRecommendationRequest, healthCtrl.AdminUpsertRecommendation, jwtMW)
	resty.Endpoint(router, requests.NewAdminDeleteRecommendationRequest, healthCtrl.AdminDeleteRecommendation, jwtMW)
	resty.Endpoint(router, requests.NewAdminUpsertCriterionRequest, healthCtrl.AdminUpsertCriterion, jwtMW)

	router.MuxRouter().HandleFunc("/api/v1/health/lab-import", healthCtrl.UploadLabDocuments).Methods(http.MethodPost)
	router.RegisterAppDoc("/api/v1/health/lab-import", []string{http.MethodPost},
		"Import lab PDFs", "Multipart form field `files` (up to N PDFs), optional `user_sex` form field.",
		"", map[int]string{200: "pending_import_id, user_criteria, model_note"}, nil)

	obs.L.Info("public-api starting")
	resty.RunServer(ctx, router, func(ctx context.Context) error {
		usersConn.Close()
		healthConn.Close()
		nc.Close()
		return nil
	})

	return nil
}
