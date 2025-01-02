package entity

import (
	"errors"
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
	ErrEmailAlreadyUsed = errors.New("email is already used")
)

var (
	CodeEmailAlreadyUsed = 10000
)

func NewUser(userName, email, password string) (*User, error) {
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
