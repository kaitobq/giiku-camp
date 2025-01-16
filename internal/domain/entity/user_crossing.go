package entity

import "time"

type UserCrossing struct {
	ID            string
	UserID        string
	CrossedUserID string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewUserCrossing(userID string, crossedUserID string) (*UserCrossing, error) {
	id, err := genID()
	if err != nil {
		return nil, err
	}
	return &UserCrossing{
		ID:            id,
		UserID:        userID,
		CrossedUserID: crossedUserID,
	}, nil
}

func (u *UserCrossing) UpdateCreatedAt() {
	u.CreatedAt = time.Now()
}

func (u *UserCrossing) UpdateUpdatedAt() {
	u.UpdatedAt = time.Now()
}
