package entity

import (
	"errors"
	"fmt"
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
	ErrFriendRequestNotFound  = errors.New("friend request not found")
	ErrSentRequestNotFound    = errors.New("sent request not found")
)

var (
	CodeFriendRequestNotFound = 30000
	CodeSentRequestNotFound   = 30001
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

func (u *UserFriendList) RemoveFriend(friendID string) {
	for i, friend := range u.Friends {
		if friend.FriendID == friendID {
			u.Friends = append(u.Friends[:i], u.Friends[i+1:]...)
			return
		}
	}
}

func (u *UserFriendList) RemoveFriendRequest(requesterID string) {
	for i, request := range u.FriendRequests {
		if request.RequesterID == requesterID {
			u.FriendRequests = append(u.FriendRequests[:i], u.FriendRequests[i+1:]...)
			return
		}
	}
}

func (u *UserFriendList) RemoveSentRequest(requesterID string) {
	for i, request := range u.SentRequests {
		if request.RequesterID == requesterID {
			u.SentRequests = append(u.SentRequests[:i], u.SentRequests[i+1:]...)
			return
		}
	}
}

func (u *UserFriendList) HasFriend(friendID string) bool {
	for _, friend := range u.Friends {
		if friend.FriendID == friendID {
			return true
		}
	}
	return false
}

func (u *UserFriendList) HasFriendRequest(requesterID string) bool {
	fmt.Println("u.FriendRequests", u.FriendRequests, "requesterID", requesterID)
	for _, request := range u.FriendRequests {
		if request.RequesterID == requesterID {
			return true
		}
	}
	return false
}

func (u *UserFriendList) HasSentRequest(requesterID string) bool {
	for _, request := range u.SentRequests {
		if request.RequesterID == requesterID {
			return true
		}
	}
	return false
}
