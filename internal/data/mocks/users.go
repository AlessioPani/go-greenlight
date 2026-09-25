package mocks

import (
	"context"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/data"
)

var InactiveUser = data.User{
	ID:        1,
	CreatedAt: time.Now(),
	Name:      "John Doe",
	Email:     "j.doe@mail.com",
	Activated: false,
	Version:   1,
}

var ActiveUser = data.User{
	ID:        2,
	CreatedAt: time.Now(),
	Name:      "Ada Wong",
	Email:     "a.wong@mail.com",
	Activated: true,
	Version:   1,
}

var InvalidUser = data.User{
	ID:        3,
	CreatedAt: time.Now(),
	Name:      "Jackie Chan",
	Email:     "j.chan",
	Activated: false,
	Version:   1,
}

type UserModel struct{}

func (u *UserModel) Insert(_ context.Context, user *data.User) error {
	if user.Email == ActiveUser.Email {
		return data.ErrDuplicateEmail
	}

	user.ID = 1

	return nil
}

func (u *UserModel) Register(ctx context.Context, user *data.User, ttl time.Duration) (*data.Token, error) {
	if err := u.Insert(ctx, user); err != nil {
		return nil, err
	}
	return &data.Token{Plaintext: "activationToken", UserID: user.ID, Scope: data.ScopeActivation}, nil
}

func (u *UserModel) Update(_ context.Context, user *data.User) error {
	if user.ID == 1 || user.ID == 2 {
		return nil
	}

	return data.ErrEditConflict
}

func (u *UserModel) UpdateAndDeleteTokens(ctx context.Context, user *data.User, scope string) error {
	return u.Update(ctx, user)
}

func (u *UserModel) GetByEmail(_ context.Context, email string) (*data.User, error) {
	switch email {
	case ActiveUser.Email:
		ActiveUser.Password.Set("valid_password")
		return &ActiveUser, nil
	case InactiveUser.Email:
		InactiveUser.Password.Set("valid_password")
		return &InactiveUser, nil
	}

	return nil, data.ErrRecordNotFound
}

func (u *UserModel) GetForToken(_ context.Context, tokenScope string, tokenPlaintext string) (*data.User, error) {
	switch tokenPlaintext {
	case "expiredtoken123456789token":
		return nil, data.ErrRecordNotFound
	case "validtokeninvaliduser12345":
		return &InvalidUser, nil
	}

	switch tokenScope {
	case data.ScopeActivation:
		return &InactiveUser, nil
	case data.ScopeAuthentication:
		return &ActiveUser, nil
	case data.ScopePasswordReset:
		return &ActiveUser, nil
	}

	return nil, data.ErrRecordNotFound
}
