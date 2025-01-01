package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(userName, email, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       genUUID(),
		UserName: userName,
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
