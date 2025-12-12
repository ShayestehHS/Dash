package postgresql

import (
	"database/sql"
	"fmt"

	"Dash/core/user"
	"Dash/pkg/logger"

	"github.com/Masterminds/squirrel"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) user.Repository {
	return &repository{db: db}
}

func (r *repository) GetByPhone(phoneNumber user.PhoneNumber) (*user.User, error) {
	phoneStr := phoneNumber.String()
	query, args, err := squirrel.Select("id", "phone_number", "password_hash", "name", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"phone_number": phoneStr}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		logger.Error("GetByPhone:QUERY_BUILD_FAILED", map[string]interface{}{
			"phone": phoneStr,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	u := &user.User{}
	var phoneDB string
	err = r.db.QueryRow(query, args...).Scan(
		&u.ID,
		&phoneDB,
		&u.PasswordHash,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		logger.Error("GetByPhone:QUERY_EXEC_FAILED", map[string]interface{}{
			"phone_number": phoneStr,
			"error":        err.Error(),
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	u.PhoneNumber, err = user.NewPhoneNumber(phoneDB)
	if err != nil {
		logger.Error("GetByPhone:PHONE_NUMBER_INVALID", map[string]interface{}{
			"phone_number": phoneStr,
			"error":        err.Error(),
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return u, nil
}


