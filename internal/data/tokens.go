package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/validator"
)

// ScopeActivation is the scope related to the activation of a user.
const ScopeActivation = "activation"

// ScopeAuthentication is the scope related to the authentication of a user.
const ScopeAuthentication = "authentication"

// ScopePasswordReset is the scope related to the password reset of a user.
const ScopePasswordReset = "password-reset"

// Token is a struct used to hold the data for an individual token. This includes the
// plaintext and hashed versions of the token, associated user ID, expiry time and
// scope.
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// generateToken is a function used to generate a Token instance.
func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Create a Token instance containing the user ID, expiry, and scope information.
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Create a slice of 16 bytes.
	randomBytes := make([]byte, 16)

	// Fill randomBytes with random bytes from the OS's CSPRNG (cryptographically secure random
	// number generator).
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Generate the plaintext token.
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate the hashed version of the token (SHA-256).
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

// Check that the plaintext token has been provided and is exactly 26 bytes long.
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// Interface for the token model.
type TokenModelInterface interface {
	New(userID int64, ttl time.Duration, scope string) (*Token, error)
	Insert(token *Token) error
	DeleteAllForUser(scope string, userID int64) error
}

// Define the TokenModel type.
type TokenModel struct {
	DB *sql.DB
}

// The New() method is a shortcut which creates a new Token struct and then inserts the
// data in the tokens table.
func (m *TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

// Insert() adds the data for a specific token to the tokens table.
func (m *TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope)
		      VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)

	return err
}

// DeleteAllForUser() deletes all tokens for a specific user and scope.
func (m *TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `DELETE FROM tokens
		      WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)

	return err
}
