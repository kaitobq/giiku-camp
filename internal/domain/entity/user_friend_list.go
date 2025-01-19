package entity

import (
	"errors"
	"time"
)

type FriendEntry struct {
	FriendID  string    `datastore:"friend_id"`
	CreatedAt time.Time `datastore:"created_at"`
}

type FriendRequestEntry struct {
	RequesterID string    `datastore:"requester_id"`
	CreatedAt   time.Time `datastore:"created_at"`
}

type UserFriendList struct {
	UserID         string               `datastore:"user_id"`
	Friends        []FriendEntry        `datastore:"friends,noindex"`
	FriendRequests []FriendRequestEntry `datastore:"friend_requests,noindex"`
	SentRequests   []FriendRequestEntry `datastore:"sent_requests,noindex"`
}

var (
	ErrUserFriendListNotFound = errors.New("user friend list not found")
)

func NewUserFriendList(userID string) *UserFriendList {
	return &UserFriendList{
		UserID: userID,
	}
}

func (u *UserFriendList) AddFriend(friendID string) {
	u.Friends = append(u.Friends, FriendEntry{
		FriendID:  friendID,
		CreatedAt: time.Now(),
	})
}

func (u *UserFriendList) AddFriendRequest(requesterID string) {
	u.FriendRequests = append(u.FriendRequests, FriendRequestEntry{
		RequesterID: requesterID,
		CreatedAt:   time.Now(),
	})
}

func (u *UserFriendList) AddSentRequest(requesterID string) {
	u.SentRequests = append(u.SentRequests, FriendRequestEntry{
		RequesterID: requesterID,
		CreatedAt:   time.Now(),
	})
}
