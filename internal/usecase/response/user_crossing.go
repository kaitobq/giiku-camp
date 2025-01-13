package response

import "giiku-camp/internal/domain/entity"

// TODO: すれ違いAPIのレスポンスを定義する
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
