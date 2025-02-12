package mocks

import (
	"errors"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/data"
)

var activationToken = data.Token{
	Plaintext: "activationToken",
	Hash:      []byte{},
	UserID:    1,
	Expiry:    time.Now(),
	Scope:     data.ScopeActivation,
}

var authToken = data.Token{
	Plaintext: "authenticationToken",
	Hash:      []byte{},
	UserID:    1,
	Expiry:    time.Now(),
	Scope:     data.ScopeAuthentication,
}

var passToken = data.Token{
	Plaintext: "passwordResetToken",
	Hash:      []byte{},
	UserID:    1,
	Expiry:    time.Now(),
	Scope:     data.ScopePasswordReset,
}

type TokenModel struct{}

func (t *TokenModel) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	if userID == 2 {
		return nil, errors.New("while generating a token")
	}

	switch scope {
	case data.ScopeActivation:
		return &activationToken, nil
	case data.ScopeAuthentication:
		return &authToken, nil
	case data.ScopePasswordReset:
		return &passToken, nil
	default:
		return nil, nil
	}
}

func (t *TokenModel) Insert(token *data.Token) error {
	return nil
}

func (t *TokenModel) DeleteAllForUser(scope string, userID int64) error {
	if scope == data.ScopeActivation && userID == ActiveUser.ID {
		return data.ErrRecordNotFound
	}

	return nil
}
