package response

import (
	"giiku-camp/internal/domain/entity"
)

type RegisterUserCrossingRes struct {
	Users []UserRes `json:"users"`
}

func NewRegisterUserCrossingRes(users []entity.User) (*RegisterUserCrossingRes, error) {
	var res []UserRes
	for _, user := range users {
		u, err := NewUserRes(&user)
		if err != nil {
			return nil, err
		}
		res = append(res, *u)
	}
	return &RegisterUserCrossingRes{
		Users: res,
	}, nil
}

type GetUserCrossingRes struct {
	Users []UserRes `json:"users"`
}

func NewGetUserCrossingRes(users []entity.User) (*GetUserCrossingRes, error) {
	var res []UserRes
	for _, user := range users {
		u, err := NewUserRes(&user)
		if err != nil {
			return nil, err
		}
		res = append(res, *u)
	}
	return &GetUserCrossingRes{
		Users: res,
	}, nil
}
