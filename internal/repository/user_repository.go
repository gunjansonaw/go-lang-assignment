package repository

import (
	"database/sql"
	"time"

	"github.com/gunjan/user-api/internal/models"
	"go.uber.org/zap"
)

type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// CreateUser creates a new user and returns the ID
func (r *UserRepository) CreateUser(name string, dob time.Time) (int64, error) {
	result, err := r.db.Exec("INSERT INTO users (name, dob) VALUES (?, ?)", name, dob.Format("2006-01-02"))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	var dobStr string

	err := r.db.QueryRow(
		"SELECT id, name, dob, created_at, updated_at FROM users WHERE id = ? LIMIT 1",
		id,
	).Scan(&user.ID, &user.Name, &dobStr, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	user.DOB, err = time.Parse("2006-01-02", dobStr)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users with pagination
func (r *UserRepository) GetAllUsers(offset, limit int) ([]models.User, int64, error) {
	// Get total count
	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	rows, err := r.db.Query(
		"SELECT id, name, dob, created_at, updated_at FROM users ORDER BY id LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		var dobStr string

		err := rows.Scan(&user.ID, &user.Name, &dobStr, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		user.DOB, err = time.Parse("2006-01-02", dobStr)
		if err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	return users, total, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(id int64, name string, dob time.Time) error {
	result, err := r.db.Exec(
		"UPDATE users SET name = ?, dob = ? WHERE id = ?",
		name, dob.Format("2006-01-02"), id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id int64) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

