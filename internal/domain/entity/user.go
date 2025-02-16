package entity

import (
	"errors"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type User struct {
	ID           string
	AppleID      string
	Name         string
	TokenVersion int
	Description  string
	GitHubID     string
	QiitaID      string
	ZennID       string
	XID          string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrPasswordIncorrect    = errors.New("password is incorrect")
	ErrTokenVersionMismatch = errors.New("token version mismatch")
	ErrFailedToParseClaims  = errors.New("failed to parse claims")
	ErrTokenInValid         = errors.New("token is invalid or expired")
	ErrAppleIDAlreadyUsed   = errors.New("apple id is already used")
)

var (
	CodeUserNotFound         = 10001
	CodePasswordIncorrect    = 10003
	CodeTokenVersionMismatch = 10004
	CodeFailedToParseClaims  = 10005
	CodeTokenInValid         = 10006
	CodeAppleIDAlreadyUsed   = 10007
)

func NewUser(userName string, gitHubID, qiitaID, zennID, xID *string) (*User, error) {
	id, err := genID()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		Name:         userName,
		TokenVersion: 1,
		GitHubID:     safeDereference(gitHubID),
		QiitaID:      safeDereference(qiitaID),
		ZennID:       safeDereference(zennID),
		XID:          safeDereference(xID),
	}, nil
}

func (u *User) IsAuthedByApple() bool {
	return u.AppleID != ""
}

func (u *User) SetAppleID(appleID string) {
	u.AppleID = appleID
}

func (u *User) GetAppleID() string {
	return u.AppleID
}

func (u *User) UpdateCreatedAt() {
	u.CreatedAt = time.Now()
}

func (u *User) UpdateUpdatedAt() {
	u.UpdatedAt = time.Now()
}

func (u *User) IncrementTokenVersion() {
	u.TokenVersion++
}

func genID() (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	return id, nil
}

func safeDereference(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
