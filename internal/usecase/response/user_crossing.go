package response

import (
	"giiku-camp/internal/domain/entity"
)

type CrossedUserRes struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	GitHubID string `json:"github_id"`
	QiitaID  string `json:"qiita_id"`
	ZennID   string `json:"zenn_id"`
	XID      string `json:"x_id"`
}

func NewCrossedUserRes(user entity.User) CrossedUserRes {
	return CrossedUserRes{
		ID:       user.ID,
		Name:     user.Name,
		GitHubID: user.GitHubID,
		QiitaID:  user.QiitaID,
		ZennID:   user.ZennID,
		XID:      user.XID,
	}
}

type RegisterUserCrossingRes struct {
	Users []CrossedUserRes `json:"users"`
}

func NewRegisterUserCrossingRes(users []entity.User) (*RegisterUserCrossingRes, error) {
	var res []CrossedUserRes
	for _, user := range users {
		u := NewCrossedUserRes(user)
		res = append(res, u)
	}
	return &RegisterUserCrossingRes{
		Users: res,
	}, nil
}

type GetUserCrossingRes struct {
	Users []CrossedUserRes `json:"users"`
}

func NewGetUserCrossingRes(users []entity.User) *GetUserCrossingRes {
	var res []CrossedUserRes
	for _, user := range users {
		u := NewCrossedUserRes(user)
		res = append(res, u)
	}
	return &GetUserCrossingRes{
		Users: res,
	}
}
