package token_test

import (
	"testing"
	"time"

	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/stretchr/testify/assert"
)

const (
	tokenValidImmediately         = 0 * time.Minute
	tokenValidAfterFifteenMinutes = 15 * time.Minute
	tokenValidFifteenMinutes      = 15 * time.Minute
	tokenValidSevenDays           = 7 * 24 * time.Hour
)

func TestJWTEncode(t *testing.T) {
	tests := []struct {
		identifier    int64
		duration      time.Duration
		validAfter    time.Duration
		wantError     bool
		expectedError error
	}{
		{
			identifier:    1,
			duration:      tokenValidFifteenMinutes,
			validAfter:    tokenValidImmediately,
			wantError:     false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			tokenable := token.NewJWT("secret")

			tokenStr, err := tokenable.Encode(tt.identifier, tt.duration, tt.validAfter)

			if tt.wantError {
				assert.Empty(t, tokenStr)
				assert.NotNil(t, err)

				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NotEmpty(t, tokenStr)
				assert.Nil(t, err)
			}
		})
	}
}

func TestJWTDecode(t *testing.T) {
	tests := []struct {
		name           string
		identifier     int64
		duration       time.Duration
		validAfter     time.Duration
		wantError      bool
		expectedError  error
		expectedResult *token.Claims
	}{
		{
			name:          "success",
			identifier:    1,
			duration:      tokenValidFifteenMinutes,
			validAfter:    tokenValidImmediately,
			wantError:     false,
			expectedError: nil,
			expectedResult: &token.Claims{
				Identifier: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenable := token.NewJWT("secret")

			tokenStr, _ := tokenable.Encode(tt.identifier, tt.duration, tt.validAfter)

			claims, err := tokenable.Decode(tokenStr)

			if tt.wantError {
				assert.Nil(t, claims)
				assert.NotNil(t, err)

				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NotNil(t, claims)
				assert.Nil(t, err)

				assert.Equal(t, claims.Identifier, tt.identifier)
			}
		})
	}
}
