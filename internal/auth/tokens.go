package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/StelIify/feedbland/internal/database"
	"golang.org/x/crypto/bcrypt"
)

const (
	ScopeAuthentication = "authentication"
	ScopeActivation     = "activation"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

func CreateUserToken(ctx context.Context, db database.Querier, userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, fmt.Errorf("error during token generation: %v", err)
	}
	err = db.CreateToken(ctx, database.CreateTokenParams{
		Hash:   token.Hash,
		UserID: userID,
		Expiry: token.Expiry,
		Scope:  token.Scope,
	})
	if err != nil {
		return nil, fmt.Errorf("error during token creation: %v", err)
	}
	return token, nil
}

func GenerePasswordHash(password string, cost int) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func ValidateCredentials(validPassword []byte, inputPassword string) (error, bool) {
	err := bcrypt.CompareHashAndPassword(validPassword, []byte(inputPassword))
	if err != nil {
		return err, false
	}
	return nil, true
}

func GenerateTokenHash(token string) []byte {
	tokenHash := sha256.Sum256([]byte(token))
	return tokenHash[:]
}
