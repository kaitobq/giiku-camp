package response

import "giiku-camp/internal/domain/entity"

type UserFriendRes struct {
	Friends          []UserRes `json:"friends"`
	SentRequests     []UserRes `json:"sent_requests"`
	ReceivedRequests []UserRes `json:"received_requests"`
}

func NewUserFriendRes(userFriends []entity.UserFriend, users []entity.User) (*UserFriendRes, error) {
	res := &UserFriendRes{}
	userMap := map[string]entity.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	for _, userFriend := range userFriends {
		user, ok := userMap[userFriend.RelatedUserID]
		if !ok {
			return nil, entity.ErrUserNotFound
		}
		switch userFriend.Type {
		case entity.UserFriendTypeFriend:
			res.Friends = append(res.Friends, *NewUserRes(&user))
		case entity.UserFriendTypeSentRequest:
			res.SentRequests = append(res.SentRequests, *NewUserRes(&user))
		case entity.UserFriendTypeReceivedRequest:
			res.ReceivedRequests = append(res.ReceivedRequests, *NewUserRes(&user))
		}
	}
	return res, nil
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
