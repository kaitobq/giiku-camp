package response

import "giiku-camp/internal/domain/entity"

type UserFriendListRes struct {
	Friends        []entity.FriendEntry        `json:"friends"`
	FriendRequests []entity.FriendRequestEntry `json:"friend_requests"`
	SentRequests   []entity.FriendRequestEntry `json:"sent_requests"`
}

func NewUserFriendListRes(friendList entity.UserFriendList) *UserFriendListRes {
	return &UserFriendListRes{
		Friends:        friendList.Friends,
		FriendRequests: friendList.FriendRequests,
		SentRequests:   friendList.SentRequests,
	}
}

type SendRequestRes struct {
	Ok bool `json:"ok"`
}

func NewSendRequestRes() (*SendRequestRes, error) {
	return &SendRequestRes{Ok: true}, nil
}

type AcceptRequestRes struct {
	Ok bool `json:"ok"`
}

func NewAcceptRequestRes() (*AcceptRequestRes, error) {
	return &AcceptRequestRes{Ok: true}, nil
}

type RejectRequestRes struct {
	Ok bool `json:"ok"`
}

func NewRejectRequestRes() (*RejectRequestRes, error) {
	return &RejectRequestRes{Ok: true}, nil
}
