package data

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/lib/pq"
)

// Permissions is a slice of string used to contain the permission for
// a single user.
type Permissions []string

// Include is a method used to check if a specific permission is included
// in the Permission slice.
func (p Permissions) Include(code string) bool {
	return slices.Contains(p, code)
}

// Permission model struct that wraps a db connection pool.
type PermissionModel struct {
	DB *sql.DB
}

// GetAllForUser is a method that retrieves all the permission for a specific
// user from the DB.
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	// SQL query to retreive all movie records.
	query := `SELECT permissions.code
			  FROM permissions
			  INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
			  INNER JOIN users ON users.id = users_permissions.user_id
			  WHERE users.id = $1`

	// Creates a context with a 3 seconds timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Executes the query.
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	// Scan rows to fill the result.
	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	// Checks again for errors.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// AddForUser is a method used to add permissions to a user.
func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	// SQL query to retreive all movie records.
	query := `INSERT INTO users_permissions
			  SELECT $1, permissions.id FROM permissions
			  WHERE permissions.code = ANY($2)`

	// Creates a context with a 3 seconds timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Executes the query.
	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))

	return err
}
