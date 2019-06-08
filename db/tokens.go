package db

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	bolt "go.etcd.io/bbolt"
	"time"
)

type Token struct {
	SigningKey string `json:"signing-key"`
	Username   string `json:"username"`
}

func NewToken(user User, db *bolt.DB) (string, error) {
	// Generate key
	signingKey := make([]byte, 128)
	if _, err := rand.Read(signingKey); err != nil {
		return "", fmt.Errorf("failed to generate JWT signing key: "+err.Error())
	}

	// Create claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + (60*60*24),
		Issuer: "dns.krantz.dev",
		IssuedAt: time.Now().Unix(),
		Subject: user.Username,
	}

	// Generate token
	user.Tokens++
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token.Header["kid"] = fmt.Sprintf("%s-%v", user.Username, user.Tokens)

	// Sign token
	signed, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	// Encode to JSON
	t := Token{
		SigningKey: base64.StdEncoding.EncodeToString(signingKey),
		Username: user.Username,
	}
	j, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	// Save to database
	if err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("tokens")).Put([]byte(fmt.Sprintf("%s-%v", user.Username, user.Tokens)), j)
	}); err != nil {
		return "", err
	}

	// Save updates to number of tokens
	if err := user.Encode(db); err != nil {
		return "", err
	}

	return signed, nil
}

func TokenFromDatabase(username string, id int64, db *bolt.DB) (Token, error) {
	var t Token

	err := db.View(func(tx *bolt.Tx) error {
		tokens := tx.Bucket([]byte("tokens"))

		data := tokens.Get([]byte(fmt.Sprintf("%s-%v", username, id)))
		if len(data) == 0 {
			return fmt.Errorf("token not found in database")
		}

		return json.Unmarshal(data, &t)
	})

	return t, err
}
