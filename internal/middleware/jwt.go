package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	restyerrors "github.com/porebric/resty/errors"
	restymiddleware "github.com/porebric/resty/middleware"
	"github.com/porebric/resty/requests"
)

type contextKey string

const (
	UserIDCtxKey   contextKey = "user_id"
	UserTypeCtxKey contextKey = "user_type"
)

func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(UserIDCtxKey).(string)
	return v
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"uid"`
	UserType string `json:"utype"`
}

type JWTAuth struct {
	secret []byte
	next   restymiddleware.Middleware
}

func NewJWTAuth(secret string) func() restymiddleware.Middleware {
	return func() restymiddleware.Middleware {
		return &JWTAuth{secret: []byte(secret)}
	}
}

func (m *JWTAuth) Execute(ctx context.Context, req requests.Request) (context.Context, int32, string) {
	type authCarrier interface {
		GetAuthToken() string
	}

	if ac, ok := req.(authCarrier); ok {
		token := ac.GetAuthToken()
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimPrefix(token, "bearer ")

		if token == "" {
			return ctx, restyerrors.ErrorUserUnauthorized, "missing token"
		}

		claims, err := validateJWT(token, m.secret)
		if err != nil {
			return ctx, restyerrors.ErrorUserUnauthorized, "invalid token"
		}

		ctx = context.WithValue(ctx, UserIDCtxKey, claims.UserID)
		ctx = context.WithValue(ctx, UserTypeCtxKey, claims.UserType)
	}

	if m.next != nil {
		return m.next.Execute(ctx, req)
	}
	return ctx, 0, ""
}

func (m *JWTAuth) SetNext(next restymiddleware.Middleware) {
	m.next = next
}

func validateJWT(tokenStr string, secret []byte) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

func ExtractBearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	h = strings.TrimPrefix(h, "Bearer ")
	h = strings.TrimPrefix(h, "bearer ")
	return h
}
