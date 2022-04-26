package token

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Payload struct {
	Username  string    `json:"username"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) *Payload {
	return &Payload{
		Username:  username,
		ExpiredAt: time.Now().Add(duration),
	}
}

func (p *Payload) Validate() error {
	if time.Now().After(p.ExpiredAt) {
		return status.Errorf(codes.Unauthenticated, "token has expired")
	}

	return nil
}
