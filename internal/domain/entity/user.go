package entity

import (
	"errors"
	"regexp"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Name         string
	Email        string
	Password     string
	TokenVersion int
	GitHubID     string
	QiitaID      string
	ZennID       string
	XID          string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var (
	ErrEmailAlreadyUsed     = errors.New("email is already used")
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailInvalid         = errors.New("invalid email")
	ErrPasswordIncorrect    = errors.New("password is incorrect")
	ErrTokenVersionMismatch = errors.New("token version mismatch")
	ErrFailedToParseClaims  = errors.New("failed to parse claims")
	ErrTokenInValid         = errors.New("token is invalid or expired")
)

var (
	CodeEmailAlreadyUsed     = 10000
	CodeUserNotFound         = 10001
	CodeEmailInvalid         = 10002
	CodePasswordIncorrect    = 10003
	CodeTokenVersionMismatch = 10004
	CodeFailedToParseClaims  = 10005
	CodeTokenInValid         = 10006
)

func NewUser(userName, email, password string, gitHubID, qiitaID, zennID, xID *string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrEmailInvalid
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	id, err := genID()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		Name:         userName,
		Email:        email,
		Password:     hashedPassword,
		TokenVersion: 1,
		GitHubID:     safeDereference(gitHubID),
		QiitaID:      safeDereference(qiitaID),
		ZennID:       safeDereference(zennID),
		XID:          safeDereference(xID),
	}, nil
}

func (u *User) UpdateCreatedAt() {
	u.CreatedAt = time.Now()
}

func (u *User) UpdateUpdatedAt() {
	u.UpdatedAt = time.Now()
}

// DBのハッシュ化されたパスワードと入力された生のパスワードを比較する
func (u *User) VerifyPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return ErrPasswordIncorrect
		default:
			return err
		}
	}
	return nil
}

// ログに出力するときにパスワードを隠す
func (u *User) HidePassword() User {
	user := *u
	user.Password = "***"
	return user
}

func (u *User) IncrementTokenVersion() {
	u.TokenVersion++
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func genID() (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	return id, nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func safeDereference(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
