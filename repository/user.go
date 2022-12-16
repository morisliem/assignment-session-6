package repository

import (
	"database/sql"
	"session-6/model"
)

type UserRepository struct {
	DB *sql.DB
}

func (s UserRepository) RegisterUser(user model.User) error {
	query := `
    INSERT user (
        username,
        first_name,
        last_name,
        password
    ) VALUES (?,?,?,?)
    `
	_, err := s.DB.Exec(query, user.Username, user.FirstName, user.LastName, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s UserRepository) GetUserByUsername(username string) (user model.User, err error) {
	query := `
    SELECT
        username,
        first_name,
        last_name,
        password
    FROM user
    WHERE username = ?
    `

	err = s.DB.QueryRow(query, username).Scan(
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Password,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
