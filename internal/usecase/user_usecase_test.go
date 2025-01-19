package usecase

import (
	"context"
	"errors"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/xcontext"
	"giiku-camp/internal/infra/logging"
	mock_repository "giiku-camp/internal/mock/domain/repository"
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_SignUp(t *testing.T) {
	logging.Init()

	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	os.Setenv("ACCESS_TOKEN_LIFE_SPAN", "1")
	defer os.Unsetenv("ACCESS_TOKEN_LIFE_SPAN")

	os.Setenv("REFRESH_TOKEN_LIFE_SPAN", "720")
	defer os.Unsetenv("REFRESH_TOKEN_LIFE_SPAN")

	tests := []struct {
		name             string
		input            request.SignUpReq
		mockSetup        func(*mock_repository.MockUserRepo, *mock_repository.MockUserFriendListRepo)
		expectedResponse response.SignUpRes
		wantErr          bool
		expectedError    error
	}{
		{
			name: "正常系: ユーザー登録成功",
			input: request.SignUpReq{
				Email:    "test@example.com",
				Name:     "test",
				Password: "password123",
			},
			mockSetup: func(mock *mock_repository.MockUserRepo, friendMock *mock_repository.MockUserFriendListRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(nil, entity.ErrUserNotFound)
				mock.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&entity.User{})).
					Return(nil)
				friendMock.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedResponse: response.SignUpRes{
				User: response.UserRes{
					ID:        "user_id",
					Email:     "test@example.com",
					Name:      "test",
					GitHubID:  "",
					QiitaID:   "",
					ZennID:    "",
					XID:       "",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				Token: response.TokenRes{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				},
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "異常系: メールアドレスが不正",
			input: request.SignUpReq{
				Email:    "test",
				Name:     "test",
				Password: "password123",
			},
			mockSetup:        func(mock *mock_repository.MockUserRepo, friendMock *mock_repository.MockUserFriendListRepo) {},
			expectedResponse: response.SignUpRes{},
			wantErr:          true,
			expectedError:    entity.ErrEmailInvalid,
		},
		{
			name: "異常系: メールアドレスが既に使用されている",
			input: request.SignUpReq{
				Email:    "test@example.com",
				Name:     "test",
				Password: "password123",
			},
			mockSetup: func(mock *mock_repository.MockUserRepo, friendMock *mock_repository.MockUserFriendListRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(&entity.User{}, nil)
			},
			expectedResponse: response.SignUpRes{},
			wantErr:          true,
			expectedError:    entity.ErrEmailAlreadyUsed,
		},
		{
			name: "異常系: ユーザー検索時に予期しないエラー",
			input: request.SignUpReq{
				Email:    "test@example.com",
				Name:     "test",
				Password: "password123",
			},
			mockSetup: func(mock *mock_repository.MockUserRepo, friendMock *mock_repository.MockUserFriendListRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(nil, errors.New("unexpected error"))
			},
			expectedResponse: response.SignUpRes{},
			wantErr:          true,
			expectedError:    errors.New("unexpected error"),
		},
		{
			name: "異常系: ユーザー作成時に予期しないエラー",
			input: request.SignUpReq{
				Email:    "test@example.com",
				Name:     "test",
				Password: "password123",
			},
			mockSetup: func(mock *mock_repository.MockUserRepo, friendMock *mock_repository.MockUserFriendListRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(nil, entity.ErrUserNotFound)
				mock.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("unexpected error"))
			},
			expectedResponse: response.SignUpRes{},
			wantErr:          true,
			expectedError:    errors.New("unexpected error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			slog.SetDefault(logger)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_repository.NewMockUserRepo(ctrl)
			friendListRepo := mock_repository.NewMockUserFriendListRepo(ctrl)
			uc := NewUserUsecase(userRepo, friendListRepo)

			tc.mockSetup(userRepo, friendListRepo)

			ctx := context.Background()
			c, _ := gin.CreateTestContext(nil)
			c.Request = (&http.Request{}).WithContext(ctx)
			res, err := uc.SignUp(c, tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("want error, but got nil")
				}
				assert.Equal(t, tc.expectedError, err)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				assert.Equal(t, tc.input.Email, res.User.Email)
				assert.Equal(t, tc.input.Name, res.User.Name)
				assert.NotEmpty(t, res.User.ID)
				assert.NotZero(t, res.User.CreatedAt)
				// UpdatedAtはrepository.Update()で更新されるのでゼロ値
				assert.NotEmpty(t, res.Token.AccessToken)
				assert.NotEmpty(t, res.Token.RefreshToken)
			}
		})
	}
}

func TestUserUsecase_SignIn(t *testing.T) {
	logging.Init()

	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	os.Setenv("ACCESS_TOKEN_LIFE_SPAN", "1")
	defer os.Unsetenv("ACCESS_TOKEN_LIFE_SPAN")

	os.Setenv("REFRESH_TOKEN_LIFE_SPAN", "720")
	defer os.Unsetenv("REFRESH_TOKEN_LIFE_SPAN")

	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate hashed password: %v", err)
	}

	tests := []struct {
		name             string
		input            request.SignInReq
		mockSetup        func(*mock_repository.MockUserRepo)
		expectedResponse response.SignInRes
		wantErr          bool
		expectedError    error
	}{
		{
			name: "正常系: ログイン成功",
			input: request.SignInReq{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mock *mock_repository.MockUserRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(&entity.User{
						ID:           "user_id",
						Name:         "test",
						Email:        "test@example.com",
						Password:     string(hashedPassword),
						TokenVersion: 2,
						GitHubID:     "",
						QiitaID:      "",
						ZennID:       "",
						XID:          "",
						CreatedAt:    time.Now().Add(-time.Hour),
						UpdatedAt:    time.Now().Add(-time.Hour),
					},
						nil)
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResponse: response.SignInRes{
				User: response.UserRes{
					ID:        "user_id",
					Email:     "test@example.com",
					Name:      "test",
					GitHubID:  "",
					QiitaID:   "",
					ZennID:    "",
					XID:       "",
					CreatedAt: time.Time{},
					UpdatedAt: time.Now(),
				},
				Token: response.TokenRes{},
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "異常系: ユーザーが見つからない",
			input: request.SignInReq{
				Email:    "test@example.com",
				Password: password,
			},
			mockSetup: func(mock *mock_repository.MockUserRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(nil, entity.ErrUserNotFound)
			},
			expectedResponse: response.SignInRes{},
			wantErr:          true,
			expectedError:    entity.ErrUserNotFound,
		},
		{
			name: "異常系: ユーザー更新時に予期しないエラー",
			input: request.SignInReq{
				Email:    "test@example.com",
				Password: password,
			},
			mockSetup: func(mock *mock_repository.MockUserRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").Return(&entity.User{
					ID:           "user_id",
					Name:         "test",
					Email:        "test@example.com",
					Password:     string(hashedPassword),
					TokenVersion: 2,
					GitHubID:     "",
					QiitaID:      "",
					ZennID:       "",
					XID:          "",
					CreatedAt:    time.Now().Add(-time.Hour),
					UpdatedAt:    time.Now().Add(-time.Hour),
				},
					nil)
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("unexpected error"))
			},
			expectedResponse: response.SignInRes{},
			wantErr:          true,
			expectedError:    errors.New("unexpected error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			slog.SetDefault(logger)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_repository.NewMockUserRepo(ctrl)
			friendListRepo := mock_repository.NewMockUserFriendListRepo(ctrl)
			uc := NewUserUsecase(userRepo, friendListRepo)

			tc.mockSetup(userRepo)

			ctx := context.Background()
			c, _ := gin.CreateTestContext(nil)
			c.Request = (&http.Request{}).WithContext(ctx)
			res, err := uc.SignIn(c, tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("want error, but got nil")
				}
				assert.Equal(t, tc.expectedError, err)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				assert.Equal(t, tc.input.Email, res.User.Email)
				assert.NotEmpty(t, res.User.ID)
				assert.NotEmpty(t, res.User.Name)
				assert.NotZero(t, res.User.CreatedAt)
				assert.NotZero(t, res.User.UpdatedAt)
				assert.NotEqual(t, tc.expectedResponse.User.UpdatedAt, res.User.UpdatedAt)
				assert.NotEmpty(t, res.Token.AccessToken)
				assert.NotEmpty(t, res.Token.RefreshToken)
			}
		})
	}
}

func TestUserUsecase_RefreshToken(t *testing.T) {
	logging.Init()

	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	os.Setenv("ACCESS_TOKEN_LIFE_SPAN", "1")
	defer os.Unsetenv("ACCESS_TOKEN_LIFE_SPAN")

	os.Setenv("REFRESH_TOKEN_LIFE_SPAN", "720")
	defer os.Unsetenv("REFRESH_TOKEN_LIFE_SPAN")

	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate hashed password: %v", err)
	}

	tests := []struct {
		name             string
		input            request.RefreshTokenReq
		genToken         bool
		mockSetup        func(*mock_repository.MockUserRepo)
		expectedResponse response.RefreshTokenRes
		wantErr          bool
		expectedError    error
	}{
		{
			name:     "正常系: トークンリフレッシュ成功",
			genToken: true,
			mockSetup: func(mock *mock_repository.MockUserRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(&entity.User{
						ID:           "user_id",
						Name:         "test",
						Email:        "test@example.com",
						Password:     string(hashedPassword),
						TokenVersion: 1,
						GitHubID:     "",
						QiitaID:      "",
						ZennID:       "",
						XID:          "",
						CreatedAt:    time.Now().Add(-time.Hour),
						UpdatedAt:    time.Now().Add(-time.Hour),
					},
						nil)
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().
					FindByID(gomock.Any(), "user_id").
					Return(&entity.User{
						ID:           "user_id",
						Name:         "test",
						Email:        "test@example.com",
						Password:     "password123",
						TokenVersion: 2,
						GitHubID:     "",
						QiitaID:      "",
						ZennID:       "",
						XID:          "",
						CreatedAt:    time.Now().Add(-time.Hour),
						UpdatedAt:    time.Now().Add(-time.Hour),
					},
						nil)
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResponse: response.RefreshTokenRes{
				Token: response.TokenRes{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				},
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name:     "異常系: トークンバージョンエラー",
			genToken: true,
			mockSetup: func(mock *mock_repository.MockUserRepo) {
				mock.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(&entity.User{
						ID:           "user_id",
						Name:         "test",
						Email:        "test@example.com",
						Password:     string(hashedPassword),
						TokenVersion: 1,
						GitHubID:     "",
						QiitaID:      "",
						ZennID:       "",
						XID:          "",
						CreatedAt:    time.Now().Add(-time.Hour),
						UpdatedAt:    time.Now().Add(-time.Hour),
					},
						nil)
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().
					FindByID(gomock.Any(), "user_id").
					Return(&entity.User{
						ID:           "user_id",
						Name:         "test",
						Email:        "test@example.com",
						Password:     "password123",
						TokenVersion: 1,
						GitHubID:     "",
						QiitaID:      "",
						ZennID:       "",
						XID:          "",
						CreatedAt:    time.Now().Add(-time.Hour),
						UpdatedAt:    time.Now().Add(-time.Hour),
					},
						nil)
			},
			expectedResponse: response.RefreshTokenRes{},
			wantErr:          true,
			expectedError:    entity.ErrTokenVersionMismatch,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			slog.SetDefault(logger)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_repository.NewMockUserRepo(ctrl)
			friendListRepo := mock_repository.NewMockUserFriendListRepo(ctrl)
			uc := NewUserUsecase(userRepo, friendListRepo)

			tc.mockSetup(userRepo)

			ctx := context.Background()
			c, _ := gin.CreateTestContext(nil)
			c.Request = (&http.Request{}).WithContext(ctx)

			if err != nil {
				t.Fatalf("failed to sign in: %v", err)
			}
			if tc.genToken {
				signInRes, err := uc.SignIn(c, request.SignInReq{
					Email:    "test@example.com",
					Password: "password123",
				})
				if err != nil {
					t.Fatalf("failed to sign in: %v", err)
				}
				tc.input.RefreshToken = signInRes.Token.RefreshToken
			}
			res, err := uc.RefreshToken(c, tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("want error, but got nil")
				}
				assert.Equal(t, tc.expectedError, err)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				assert.NotEmpty(t, res.Token.AccessToken)
				assert.NotEmpty(t, res.Token.RefreshToken)
			}
		})
	}
}

func TestUserUsecase_UpdateMe(t *testing.T) {
	logging.Init()

	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	os.Setenv("ACCESS_TOKEN_LIFE_SPAN", "1")
	defer os.Unsetenv("ACCESS_TOKEN_LIFE_SPAN")

	os.Setenv("REFRESH_TOKEN_LIFE_SPAN", "720")
	defer os.Unsetenv("REFRESH_TOKEN_LIFE_SPAN")

	testUser := &entity.User{
		ID:        "user_id",
		Email:     "test@example.com",
		Name:      "test",
		GitHubID:  "test",
		QiitaID:   "test",
		ZennID:    "test",
		XID:       "test",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	tests := []struct {
		name             string
		input            request.UpdateMeReq
		mockSetup        func(*mock_repository.MockUserRepo)
		expectedResponse response.UserRes
		wantErr          bool
		expectedError    error
	}{
		{
			name: "正常系: ユーザー情報更新成功",
			input: request.UpdateMeReq{
				Name:     "test",
				GitHubID: "test",
				QiitaID:  "test",
				ZennID:   "test",
				XID:      "test",
			},
			mockSetup: func(mock *mock_repository.MockUserRepo) {
				mock.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&entity.User{})).
					Return(nil)
			},
			expectedResponse: response.UserRes{
				ID:        "user_id",
				Email:     "test@example.com",
				Name:      "test",
				GitHubID:  "test",
				QiitaID:   "test",
				ZennID:    "test",
				XID:       "test",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "異常系: ユーザー更新時に予期しないエラー",
			input: request.UpdateMeReq{
				Name:     "test",
				GitHubID: "test",
				QiitaID:  "test",
				ZennID:   "test",
				XID:      "test",
			},
			mockSetup: func(mock *mock_repository.MockUserRepo) {
				mock.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&entity.User{})).
					Return(errors.New("unexpected error"))
			},
			expectedResponse: response.UserRes{},
			wantErr:          true,
			expectedError:    errors.New("unexpected error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			slog.SetDefault(logger)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_repository.NewMockUserRepo(ctrl)
			friendListRepo := mock_repository.NewMockUserFriendListRepo(ctrl)
			uc := NewUserUsecase(userRepo, friendListRepo)

			tc.mockSetup(userRepo)

			ctx := context.Background()
			c, _ := gin.CreateTestContext(nil)
			c.Request = (&http.Request{}).WithContext(ctx)

			xcontext.WithUser(c, testUser)

			res, err := uc.UpdateMe(c, tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("want error, but got nil")
				}
				assert.Equal(t, tc.expectedError, err)
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				assert.Equal(t, tc.expectedResponse, *res)
			}
		})
	}
}
