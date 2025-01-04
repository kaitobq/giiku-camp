package entity

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrEmailAlreadyUsed  = errors.New("email is already used")
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailInvalid      = errors.New("invalid email")
	ErrPasswordIncorrect = errors.New("password is incorrect")
)

var (
	CodeEmailAlreadyUsed  = 10000
	CodeUserNotFound      = 10001
	CodeEmailInvalid      = 10002
	CodePasswordIncorrect = 10003
)

func NewUser(userName, email, password string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrEmailInvalid
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       genUUID(),
		Name:     userName,
		Email:    email,
		Password: hashedPassword,
	}, nil
}

func (u *User) UpdateCreatedAt() {
	u.CreatedAt = time.Now()
}

func (u *User) UpdateUpdatedAt() {
	u.UpdatedAt = time.Now()
}

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

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func genUUID() string {
	id := uuid.New()
	return id.String()
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
