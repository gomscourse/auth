package auth

import (
	"context"
	"github.com/gomscourse/auth/internal/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"strings"
)

const authPrefix = "Bearer "

func (s *serv) Check(ctx context.Context, endpointAddress string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "checking endpoint address")
	span.SetTag("endpoint", endpointAddress)
	defer span.Finish()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(s.jwtConfig.AccessTokenSecret()))
	if err != nil {
		return errors.New("access token is invalid")
	}

	accessRules, _, err := s.accessRepo.GetRuleByEndpoint(ctx, endpointAddress)
	if err != nil {
		return errors.New("failed to get access rules")
	}

	if len(accessRules) == 0 {
		return nil
	}

	accessGranted := false

	for _, rule := range accessRules {
		if rule.Role == claims.Role {
			accessGranted = true
			break
		}
	}

	if !accessGranted {
		return errors.New("access denied")
	}

	return nil
}
