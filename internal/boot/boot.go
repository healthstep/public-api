package boot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	healthpb "github.com/helthtech/core-health/pkg/proto/health"
	userspb "github.com/helthtech/core-users/pkg/proto/users"
	"github.com/helthtech/public-api/internal/actions"
	"github.com/helthtech/public-api/internal/middleware"
	"github.com/helthtech/public-api/internal/natshandler"
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

	l := logger.New(logger.InfoLevel)
	router := resty.NewRouter(func() *logger.Logger { return l }, wsHub)

	origins := strings.Split(configs.Value(ctx, "cors_origins").String(), ",")
	methods := strings.Split(configs.Value(ctx, "cors_methods").String(), ",")
	headers := strings.Split(configs.Value(ctx, "cors_headers").String(), ",")
	router.SetCors(origins, methods, headers)

	tgBotURL := configs.Value(ctx, "tg_bot_url").String()
	maxBotURL := configs.Value(ctx, "max_bot_url").String()

	authCtrl := actions.NewAuthController(usersClient, tgBotURL, maxBotURL)
	userCtrl := actions.NewUserController(usersClient)
	healthCtrl := actions.NewHealthController(healthClient)

	resty.Endpoint(router, requests.NewBrowserChallengeRequest, authCtrl.BrowserChallenge)
	resty.Endpoint(router, requests.NewGetMeRequest, userCtrl.GetMe, jwtMW)
	resty.Endpoint(router, requests.NewUpdateMeRequest, userCtrl.UpdateMe, jwtMW)
	resty.Endpoint(router, requests.NewListCriteriaRequest, healthCtrl.ListCriteria, jwtMW)
	resty.Endpoint(router, requests.NewSetUserCriterionRequest, healthCtrl.SetUserCriterion, jwtMW)
	resty.Endpoint(router, requests.NewGetUserCriteriaRequest, healthCtrl.GetUserCriteria, jwtMW)
	resty.Endpoint(router, requests.NewGetProgressRequest, healthCtrl.GetProgress, jwtMW)
	resty.Endpoint(router, requests.NewGetRecommendationsRequest, healthCtrl.GetRecommendations, jwtMW)
	resty.Endpoint(router, requests.NewResetCriteriaRequest, healthCtrl.ResetCriteria, jwtMW)

	log.Println("public-api starting")
	resty.RunServer(ctx, router, func(ctx context.Context) error {
		usersConn.Close()
		healthConn.Close()
		nc.Close()
		return nil
	})

	return nil
}
