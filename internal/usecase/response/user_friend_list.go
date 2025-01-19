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
