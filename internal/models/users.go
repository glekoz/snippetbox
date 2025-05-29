package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users(name, email, hashed_password, created)
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
	RETURNING id`
	var id int64
	err = m.DB.QueryRow(stmt, name, email, hashedPassword).Scan(&id)
	if err != nil {
		// нужно, если потом откажусь от отдельной проверки Exist
		if strings.Contains(err.Error(), "pq: duplicate key value") {
			return 0, ErrDuplicateEntry
		} else {
			return 0, err
		}
	}
	return int(id), nil
}

/*
// Нельзя использовать это как проверку на уникальность email
// перед регистрацией нового пользователя, потому что
// при одновременной регистрации с одной и той же почтой
// регистрация у одного пройдет успешно, а у второго
// произойдет 500 ошибка

	func (m *UserModel) Exist(email string) (bool, error) {
		stmt := `SELECT id FROM users WHERE email = $1`
		var id int64
		err := m.DB.QueryRow(stmt, email).Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			} else {
				return true, err
			}
		}
		return true, nil
	}
*/
func (m *UserModel) Get(email, password string) (*User, error) {
	stmt := `SELECT id, name, email, hashed_password FROM users WHERE email = $1`
	row := m.DB.QueryRow(stmt, email)

	u := &User{}
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrWrongCredentials
		} else {
			return nil, err
		}

	}
	return u, nil
}

/*
CREATE TABLE users (id SERIAL NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL, email VARCHAR(255) NOT NULL, hashed_password CHAR(60) NOT NULL, created TIMESTAMP NOT NULL);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
*/
