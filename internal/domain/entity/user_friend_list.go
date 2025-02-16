package entity

import (
	"time"
)

type FriendEntry struct {
	FriendID   string    `datastore:"friend_id" json:"friend_id"`
	FriendName string    `datastore:"friend_name" json:"friend_name"`
	CreatedAt  time.Time `datastore:"created_at" json:"created_at"`
}

type FriendRequestEntry struct {
	RequesterID   string    `datastore:"requester_id" json:"requester_id"`
	RequesterName string    `datastore:"requester_name" json:"requester_name"`
	CreatedAt     time.Time `datastore:"created_at" json:"created_at"`
}

type UserFriendList struct {
	UserID         string               `datastore:"user_id" json:"user_id"`
	Friends        []FriendEntry        `datastore:"friends,noindex" json:"friends"`
	FriendRequests []FriendRequestEntry `datastore:"friend_requests,noindex" json:"friend_requests"`
	SentRequests   []FriendRequestEntry `datastore:"sent_requests,noindex" json:"sent_requests"`
}

var (
// ErrUserFriendListNotFound = errors.New("user friend list not found")
// ErrFriendRequestNotFound  = errors.New("friend request not found")
// ErrSentRequestNotFound    = errors.New("sent request not found")
// ErrFriendRequestAlreadySent     = errors.New("friend request already sent")
// ErrFriendRequestAlreadyReceived = errors.New("friend request already received")
// ErrAlreadyFriend                = errors.New("already friend")
)

var (
// CodeFriendRequestNotFound      = 30000
// CodeSentRequestNotFound        = 30001
// CodeFriendRequestAlreadyExists = 30002
// CodeFriendRequestAlreadySent     = 30003
// CodeFriendRequestAlreadyReceived = 30004
// CodeAlreadyFriend                = 30005
)

func NewUserFriendList(userID string) *UserFriendList {
	return &UserFriendList{
		UserID: userID,
	}
}

func (u *UserFriendList) AddFriend(friend *User) {
	u.Friends = append(u.Friends, FriendEntry{
		FriendID:   friend.ID,
		FriendName: friend.Name,
		CreatedAt:  time.Now(),
	})
}

func (u *UserFriendList) AddFriendRequest(requester *User) {
	u.FriendRequests = append(u.FriendRequests, FriendRequestEntry{
		RequesterID:   requester.ID,
		RequesterName: requester.Name,
		CreatedAt:     time.Now(),
	})
}

func (u *UserFriendList) AddSentRequest(requester *User) {
	u.SentRequests = append(u.SentRequests, FriendRequestEntry{
		RequesterID:   requester.ID,
		RequesterName: requester.Name,
		CreatedAt:     time.Now(),
	})
}

func (u *UserFriendList) RemoveFriend(friend *User) {
	for i, fr := range u.Friends {
		if fr.FriendID == friend.ID {
			u.Friends = append(u.Friends[:i], u.Friends[i+1:]...)
			return
		}
	}
}

func (u *UserFriendList) RemoveFriendRequest(requester *User) {
	for i, request := range u.FriendRequests {
		if request.RequesterID == requester.ID {
			u.FriendRequests = append(u.FriendRequests[:i], u.FriendRequests[i+1:]...)
			return
		}
	}
}

func (u *UserFriendList) RemoveSentRequest(requester *User) {
	for i, request := range u.SentRequests {
		if request.RequesterID == requester.ID {
			u.SentRequests = append(u.SentRequests[:i], u.SentRequests[i+1:]...)
			return
		}
	}
}

func (u *UserFriendList) HasFriend(friend *User) bool {
	for _, fr := range u.Friends {
		if fr.FriendID == friend.ID {
			return true
		}
	}
	return false
}

func (u *UserFriendList) HasFriendRequest(requester *User) bool {
	for _, request := range u.FriendRequests {
		if request.RequesterID == requester.ID {
			return true
		}
	}
	return false
}

func (u *UserFriendList) HasSentRequest(requester *User) bool {
	for _, request := range u.SentRequests {
		if request.RequesterID == requester.ID {
			return true
		}
	}
	return false
}
