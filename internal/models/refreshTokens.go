package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"snippetbox.glebich/internal/jwtAuth"
)

type RefreshToken struct {
	ID      int
	Value   string
	UserId  int
	Expires time.Time
}

type RefreshTokenModel struct {
	DB *sql.DB
}

func (m *RefreshTokenModel) Insert(value string, expires int, userId int) error {
	/*
		hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(value), 12)
		if err != nil {
			return err
		}
	*/
	stmt := `INSERT INTO refresh_tokens(value, expires, user_id) 
	VALUES($1, CURRENT_TIMESTAMP + $2 * INTERVAL '1 day', $3)`
	_, err := m.DB.Exec(stmt, value, expires, userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *RefreshTokenModel) CheckRefreshToken(value string) (*jwtAuth.Sub, error) {
	/*
		hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(value), 12)
		if err != nil {
			return nil, err
		}
	*/
	stmt := `SELECT value, user_id, expires FROM refresh_tokens WHERE value = $1`

	var row struct {
		hashedRefreshTokenFromDB string
		userId                   int
		expires                  time.Time
	}
	err := m.DB.QueryRow(stmt, value).Scan(&row.hashedRefreshTokenFromDB, &row.userId, &row.expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	if row.expires.Unix() < time.Now().Unix() {
		m.Delete(row.userId)
		return nil, fmt.Errorf("expired token")
	}

	stmt = `SELECT id, name, email FROM users WHERE id = $1`

	user := &jwtAuth.Sub{}
	err = m.DB.QueryRow(stmt, row.userId).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *RefreshTokenModel) Delete(userId int) error {
	stmt := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := m.DB.Exec(stmt, userId)
	if err != nil {
		return err
	}
	return nil
}

/*
CREATE TABLE refresh_tokens (id SERIAL PRIMARY KEY, value CHAR(60) NOT NULL, expires TIMESTAMP NOT NULL, user_id INTEGER NOT NULL UNIQUE, FOREIGN KEY (user_id) REFERENCES users(id));
// либо можно составной первичный ключ
*/
