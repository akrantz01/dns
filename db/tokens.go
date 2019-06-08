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

func TokenFromString(tokenStr string, db *bolt.DB) (*jwt.Token, error) {
	// Retrieve token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else if _, ok := token.Header["kid"]; !ok {
			return nil, fmt.Errorf("unable to find key id in token")
		} else if _, ok := token.Header["kid"].(string); !ok {
			return nil, fmt.Errorf("token id must be a string")
		}

		// Get signing key from database
		var t Token
		if err := db.View(func(tx *bolt.Tx) error {
			data := tx.Bucket([]byte("tokens")).Get([]byte(token.Header["kid"].(string)))
			if len(data) == 0 {
				return fmt.Errorf("token not found in database")
			}
			return json.Unmarshal(data, &t)
		}); err != nil {
			return nil, err
		}

		// Decode signing key
		signingKey, err := base64.StdEncoding.DecodeString(t.SigningKey)
		if err != nil {
			return nil, fmt.Errorf("unable to decode signing key: %v", err)
		}

		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}
