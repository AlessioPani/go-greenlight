package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// Custom error for a user with an already used email.
var ErrDuplicateEmail = errors.New("duplicate email")

// Anonymous user.
var AnonymouseUser = &User{}

// User is a struct that represents an individual user.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

// User model struct that wraps a db connection pool.
type UserModel struct {
	DB *sql.DB
}

// Check if a user instance is an anonymous user.
func (u *User) IsAnonymous() bool {
	return u == AnonymouseUser
}

// Insert is a method used to add a new user to the User table.
func (m UserModel) Insert(user *User) error {
	query := `INSERT INTO users (name, email, password_hash, activated)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil

}

// GetByEmail is a method used to retrieve a user by its email.
func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, name, email, password_hash, activated, version
			  FROM users
			  WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.Password.hash, &user.Activated, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// Update is a method used to update a user in the database.
func (m UserModel) Update(user *User) error {
	query := `UPDATE users
			  SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
			  WHERE id = $5 AND version = $6
			  RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated, user.ID, user.Version}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

// GetForToken is a method used to get a user from a token specified in input.
func (m UserModel) GetForToken(tokenScope string, tokenPlaintext string) (*User, error) {
	// Get the hashed token.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
			  FROM users
			  INNER JOIN tokens
			  ON users.id = tokens.user_id
			  WHERE tokens.hash = $1
			  AND tokens.scope = $2
			  AND tokens.expiry > $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, tokenHash[:], tokenScope, time.Now()).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &User{}, ErrRecordNotFound
		default:
			return &User{}, err
		}
	}

	return &user, nil
}

// password is a struct that contains plaintext password and the related
// hash, expressed in slice of bytes.
type password struct {
	plaintext *string
	hash      []byte
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plaintestPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintestPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// ValidateEmail method validates the email by checking if it's not empty or
// if it's not a valid email.
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// ValidatePasswordPlaintext method validates the password by checking it has been provided
// and its length.
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

// ValidateUser validates a user by checking its fields.
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)

	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	// If the password hash is nil, this will be due to a logic error in the
	// codebase, thus we panic.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
