package db

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	bolt "go.etcd.io/bbolt"
)

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Tokens   int64  `json:"tokens"`
}

func NewUser(name, username, password, role string) User {
	return User{
		Name: name,
		Username: username,
		Password: password,
		Role: role,
		Tokens: 0,
	}
}

func UserFromToken(token *jwt.Token, db *bolt.DB) (User, error) {
	// Get username from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return User{}, fmt.Errorf("invalid JWT claims format")
	}

	// Get user from token
	user, err := UserFromDatabase(claims["sub"].(string), db)
	if err != nil {
		return User{}, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return user, nil
}

func UserFromDatabase(username string, db *bolt.DB) (User, error) {
	var u User

	err := db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte("users"))

		data := users.Get([]byte(username))
		if len(data) == 0 {
			return fmt.Errorf("user not found in database")
		}

		return json.Unmarshal(data, &u)
	})

	return u, err
}

func (u *User) Encode(db *bolt.DB) error {
	j, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("users")).Put([]byte(u.Username), j)
	})
}
