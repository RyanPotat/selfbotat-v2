package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID				int        `db:"id"`
	UserID    int        `db:"user_id"`
	Login     string     `db:"username"`
	Name      string     `db:"display"`
	FirstSeen time.Time  `db:"first_seen"`
}

type Channel struct {
	ID       int       `db:"id"`
	UserID   int       `db:"user_id"`
	JoinedAt time.Time `db:"joined_at"`
}

var pool *sql.DB

func InitDatabase(connStr string) error {
	var err error
	pool, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = pool.Ping()
	if err != nil {
		return err
	}

	sqlFileContent, err := os.ReadFile("bot/database/schema.sql")
	if err != nil {
		return err
	}

	_, err = pool.Query(string(sqlFileContent))
	if err != nil {
		return err
	}

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(50)
	pool.SetMaxOpenConns(50)
	return nil
}

func NewUser(id, login, name string) User {
	rows, err := pool.Query(`
		INSERT INTO users (user_id, username, display) 
		VALUES ($1, $2, $3) 
		RETURNING id, user_id, username, display
	`, 
		id, 
		login, 
		name,
	)

	if err != nil {
		return User{}
	}

	defer rows.Close()

	if rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.Login,
			&user.Name,
		)
		if err != nil {
			return User{}
		}

		return user
	}

	return User{}
}

func GetUser(id string, loginOrID bool) (User, bool) {
	queryType := "user_id"
	if loginOrID {
		queryType = "username"
	}

	queryString := fmt.Sprintf(`
		SELECT id, user_id, username, display, first_seen
		FROM users 
		WHERE %s = $1
	`, queryType)

	rows, err := pool.Query(queryString, id,
		)
	if err != nil {
		return User{}, false
	}

	defer rows.Close()

	if rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.Login,
			&user.Name,
			&user.FirstSeen,
		)
		if err != nil {
			fmt.Println("User scan error", err)
			return User{}, false
		}

		return user, true
	}

	return User{}, false
}

func NewChannel(id string) Channel {
	rows, err := pool.Query(`
  	INSERT INTO channels (user_id) 
		VALUES ($1) 
		RETURNING user_id, joined_at`, 
		id,
	)
	if err != nil {
		return Channel{}
	}

	defer rows.Close()

	if rows.Next() {
		var channel Channel
		err := rows.Scan(
			&channel.ID,
			&channel.JoinedAt,
		)
		if err != nil {
			return Channel{}
		}

		return channel
	}

	return Channel{}
}

func GetChannel(id string, loginOrID bool) Channel {	
	queryType := "user_id"
	if loginOrID {
		queryType = "username"
	}

	queryString := fmt.Sprintf(`
	  SELECT c.user_id, c.joined_at
		FROM (
			SELECT user_id, joined_at FROM channels WHERE user_id = $1
		)
		JOIN users AS u ON c.user_id = u.user_id
		WHERE u.%s = $1
	`, queryType)

	rows, err := pool.Query(queryString, id)
	if err != nil {
		return Channel{}
	}

	defer rows.Close()

	if rows.Next() {
		var channel Channel
		err := rows.Scan(
			&channel.ID,
			&channel.JoinedAt,
		)
		if err != nil {
			return Channel{}
		}

		return channel
	}

	return Channel{}
}

func RemoveChannel(id string) error {
	_, err := pool.Query(`DELETE FROM channels WHERE user_id = $1`, id)
	return err
}

func GetChannels() ([]User, error) {
	rows, err := pool.Query(`
	  SELECT u.user_id, u.username
		FROM channels AS c
		JOIN users AS u ON c.user_id = u.user_id
	`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer rows.Close()

	var channels []User
	for rows.Next() {
		var channel User
		err := rows.Scan(
			&channel.ID,
			&channel.Login,
		)
		if err != nil {
			return nil, err
		}

		channels = append(channels, channel)
	}

	return channels, nil
}
