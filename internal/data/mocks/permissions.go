package mocks

import (
	"github.com/AlessioPani/go-greenlight/internal/data"
)

var movieReadPermission = data.Permissions{"movies:read"}
var movieWritePermission = data.Permissions{"movies:read", "movies:write"}

type PermissionModel struct{}

func (p PermissionModel) GetAllForUser(userID int64) (data.Permissions, error) {
	switch userID {
	case 1:
		return movieReadPermission, nil
	case 2:
		return movieWritePermission, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (p PermissionModel) AddForUser(userID int64, codes ...string) error {
	return nil
}
