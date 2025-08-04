// Copyright The Linux Foundation and each contributor to LFX.
// SPDX-License-Identifier: MIT

package mock

import (
	"context"
	"log/slog"
	"os"

	"github.com/linuxfoundation/lfx-v2-committee-service/internal/domain/port"
	"github.com/linuxfoundation/lfx-v2-committee-service/pkg/errors"
)

// MockAuthService provides a mock implementation of the authentication service
type MockAuthService struct{}

// ParsePrincipal parses and validates a JWT token, returning a mock principal
func (m *MockAuthService) ParsePrincipal(ctx context.Context, token string, logger *slog.Logger) (string, error) {

	principal := os.Getenv("JWT_AUTH_DISABLED_MOCK_LOCAL_PRINCIPAL")

	if principal == "" {
		return "", errors.NewValidation("empty token provided")
	}

	logger.DebugContext(ctx, "parsed principal",
		"user_id", principal,
	)

	return principal, nil
}

// NewMockAuthService creates a new mock authentication service
func NewMockAuthService() port.Authenticator {
	return &MockAuthService{}
}
