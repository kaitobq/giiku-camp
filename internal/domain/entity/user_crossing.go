package entity

import "time"

type UserCrossing struct {
	ID            string
	UserID        string
	CrossedUserID string
	CreatedAt     time.Time
}

func NewUserCrossing(userID string, crossedUserID string) *UserCrossing {
	return &UserCrossing{
		ID:            genUUID(),
		UserID:        userID,
		CrossedUserID: crossedUserID,
	}
}

func (u *UserCrossing) UpdateCreatedAt() {
	u.CreatedAt = time.Now()
}
