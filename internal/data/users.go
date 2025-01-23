package data

import (
	"errors"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// User is a struct that represents an individual user.
type User struct {
	ID          int64     `json:"id"`
	CreatedDate time.Time `json:"created_date"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    password  `json:"-"`
	Activated   bool      `json:"activated"`
	Version     int       `json:"version"`
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
