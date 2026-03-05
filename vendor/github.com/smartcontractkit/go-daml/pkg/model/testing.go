package model

import (
	"time"
)

type GetTimeRequest struct{}

type GetTimeResponse struct {
	CurrentTime time.Time
}

type SetTimeRequest struct {
	CurrentTime time.Time
	NewTime     time.Time
}

type SetTimeResponse struct{}
