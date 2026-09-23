package data

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/lib/pq"
)

// Permissions contains the permission codes assigned to
// a single user.
type Permissions []string

// Include reports whether the permission code is present in
// the Permissions slice.
func (p Permissions) Include(code string) bool {
	return slices.Contains(p, code)
}

// PermissionModelInterface defines permission data operations.
type PermissionModelInterface interface {
	GetAllForUser(userID int64) (Permissions, error)
	AddForUser(userID int64, codes ...string) error
}

// PermissionModel provides access to user permissions.
type PermissionModel struct {
	DB *sql.DB
}

// GetAllForUser retrieves all permission codes assigned to a
// user.
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	// SQL query to retrieve the user permissions.
	query := `SELECT permissions.code
			  FROM permissions
			  INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
			  INNER JOIN users ON users.id = users_permissions.user_id
			  WHERE users.id = $1`

	// Create a context with a three-second timeout.
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

	// Check for errors encountered while iterating over the rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// AddForUser is a method used to add permissions to a user.
func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	// SQL query to add the requested permissions to the user.
	query := `INSERT INTO users_permissions
			  SELECT $1, permissions.id FROM permissions
			  WHERE permissions.code = ANY($2)`

	// Create a context with a three-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Executes the query.
	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))

	return err
}
