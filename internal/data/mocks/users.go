package mocks

import (
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

func (u *UserModel) Insert(user *data.User) error {
	if user.Email == ActiveUser.Email {
		return data.ErrDuplicateEmail
	}

	user.ID = 1

	return nil
}

func (u *UserModel) Update(user *data.User) error {
	if user.ID == 1 || user.ID == 2 {
		return nil
	}

	return data.ErrEditConflict
}

func (u *UserModel) GetByEmail(email string) (*data.User, error) {
	if email == ActiveUser.Email {
		return &ActiveUser, nil
	}

	return nil, data.ErrRecordNotFound
}

func (u *UserModel) GetForToken(tokenScope string, tokenPlaintext string) (*data.User, error) {
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
