package token

import (
	"github.com/garcia-paulo/upvotes-grpc/infra/config"
	"github.com/o1egl/paseto"
)

type TokenMaker struct {
	token  *paseto.V2
	config *config.Config
}

func NewTokenMaker(config *config.Config) *TokenMaker {
	maker := &TokenMaker{
		token:  paseto.NewV2(),
		config: config,
	}
	return maker
}

func (m *TokenMaker) CreateToken(username string) (string, error) {
	payload, err := NewPayload(username, m.config.TokenDuration)
	if err != nil {
		return "", err
	}

	return m.token.Encrypt([]byte(m.config.TokenKey), payload, nil)
}

func (m *TokenMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := m.token.Decrypt(token, []byte(m.config.TokenKey), payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Validate()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
