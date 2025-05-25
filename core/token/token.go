package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokenable interface {
	Encode(identifier int64, duration time.Duration, validAfter time.Duration) (string, error)
	Decode(str string) (*Claims, error)
}

type Claims struct {
	Identifier int64 `json:"identifier"`
	jwt.RegisteredClaims
}
