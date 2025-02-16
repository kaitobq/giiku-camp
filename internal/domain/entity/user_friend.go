package entity

import (
	"errors"
	"time"
)

type UserFriend struct {
	UserID        string         `json:"user_id"`
	RelatedUserID string         `json:"related_user_id"`
	Type          UserFriendType `json:"type"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type UserFriendType uint

const (
	UserFriendTypeFriend UserFriendType = iota
	UserFriendTypeSentRequest
	UserFriendTypeReceivedRequest
)

var (
	ErrUserFriendNotFound           = errors.New("user friend not found")
	ErrFriendRequestAlreadySent     = errors.New("friend request already sent")
	ErrFriendRequestAlreadyReceived = errors.New("friend request already received")
	ErrAlreadyFriend                = errors.New("already friend")
	ErrFriendRequestNotFound        = errors.New("friend request not found")
	ErrSentRequestNotFound          = errors.New("sent request not found")
)

var (
	CodeUserFriendNotFound           = 30000
	CodeFriendRequestAlreadySent     = 30001
	CodeFriendRequestAlreadyReceived = 30002
	CodeAlreadyFriend                = 30003
	CodeFriendRequestNotFound        = 30004
	CodeSentRequestNotFound          = 30005
)

func NewUserFriend(userID, relatedUserID string, t UserFriendType) *UserFriend {
	return &UserFriend{
		UserID:        userID,
		RelatedUserID: relatedUserID,
		Type:          t,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (u *UserFriend) IsFriend(userID string) bool {
	return u.Type == UserFriendTypeFriend && u.RelatedUserID == userID
}

func (u *UserFriend) IsSentRequest(userID string) bool {
	return u.Type == UserFriendTypeSentRequest && u.RelatedUserID == userID
}

func (u *UserFriend) IsReceivedRequest(userID string) bool {
	return u.Type == UserFriendTypeReceivedRequest && u.RelatedUserID == userID
}

func (u *UserFriend) AcceptRequest() {
	switch u.Type {
	case UserFriendTypeFriend:
	case UserFriendTypeSentRequest:
	case UserFriendTypeReceivedRequest:
		u.Type = UserFriendTypeFriend
		u.UpdatedAt = time.Now()
	}
}

func (u *UserFriend) AcceptedRequest() {
	switch u.Type {
	case UserFriendTypeFriend:
	case UserFriendTypeSentRequest:
		u.Type = UserFriendTypeFriend
		u.UpdatedAt = time.Now()
	case UserFriendTypeReceivedRequest:
	}
}
