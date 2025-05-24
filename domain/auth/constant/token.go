package constant

import "time"

const (
	TokenValidImmediately         = 0 * time.Minute
	TokenValidAfterFifteenMinutes = 15 * time.Minute
	TokenValidFifteenMinutes      = 15 * time.Minute
	TokenValidSevenDays           = 7 * 24 * time.Hour
)
