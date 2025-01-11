// internal/controller/user_ctrl_test.go

package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"giiku-camp/internal/controller/render"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/infra/logging"
	mock_usecase "giiku-camp/internal/mock/usecase"
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// setupRouter はモックされた usecase を使用して Gin エンジンを設定します。
func setupRouter(userUsecaseMock *mock_usecase.MockUserUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mocked UserCtrl
	userCtrl := NewUserCtrl(userUsecaseMock)

	// Define routes
	r := router.Group("/api/v1/auth")
	{
		r.POST("/signup", userCtrl.SignUp)
		r.POST("/signin", userCtrl.SignIn)
		r.POST("/refresh", userCtrl.RefreshToken)
	}

	return router
}

func TestUserCtrl_SignUp(t *testing.T) {
	logging.Init()
	// テストケースの定義
	tests := []struct {
		name               string
		inputBody          request.SignUpReq
		mockSetup          func(*mock_usecase.MockUserUsecase)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "正常系: ユーザー登録成功",
			inputBody: request.SignUpReq{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					SignUp(gomock.Any(), gomock.Eq(request.SignUpReq{
						Name:     "Test User",
						Email:    "test@example.com",
						Password: "password123",
					})).
					Return(&response.SignUpRes{
						User: response.UserRes{
							Name:      "Test User",
							Email:     "test@example.com",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Token: response.TokenRes{
							AccessToken:  "access-token",
							RefreshToken: "refresh-token",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: response.SignUpRes{
				User: response.UserRes{
					Name:      "Test User",
					Email:     "test@example.com",
					CreatedAt: time.Time{}, // テスト内ではゼロ値に設定
					UpdatedAt: time.Time{},
				},
				Token: response.TokenRes{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				},
			},
		},
		{
			name: "異常系: 無効なメールアドレス",
			inputBody: request.SignUpReq{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					SignUp(gomock.Any(), gomock.Eq(request.SignUpReq{
						Name:     "Test User",
						Email:    "invalid-email",
						Password: "password",
					})).
					Return((*response.SignUpRes)(nil), entity.ErrEmailInvalid)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   render.Error{Message: entity.ErrEmailInvalid.Error(), Code: entity.CodeEmailInvalid, Status: http.StatusBadRequest},
		},
		{
			name: "異常系: Emailが既に使用されている",
			inputBody: request.SignUpReq{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					SignUp(gomock.Any(), gomock.Eq(request.SignUpReq{
						Name:     "Test User",
						Email:    "test@example.com",
						Password: "password123",
					})).
					Return((*response.SignUpRes)(nil), entity.ErrEmailAlreadyUsed)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: render.Error{
				Message: entity.ErrEmailAlreadyUsed.Error(),
				Code:    entity.CodeEmailAlreadyUsed,
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "異常系: 内部エラー",
			inputBody: request.SignUpReq{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					SignUp(gomock.Any(), gomock.Eq(request.SignUpReq{
						Name:     "Test User",
						Email:    "test@example.com",
						Password: "password123",
					})).
					Return((*response.SignUpRes)(nil), errors.New("internal server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: render.Error{
				Message: "internal server error",
				Status:  http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ロガーを手動で設定（テスト用にログ出力を無視）
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			slog.SetDefault(logger)

			// gomock コントローラの作成
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モックの生成
			userUsecaseMock := mock_usecase.NewMockUserUsecase(ctrl)

			// モックの設定
			tc.mockSetup(userUsecaseMock)

			// ルーターのセットアップ
			router := setupRouter(userUsecaseMock)

			// リクエストボディのエンコード
			bodyBytes, err := json.Marshal(tc.inputBody)
			assert.NoError(t, err)

			// テストリクエストの作成
			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// レスポンスの記録
			w := httptest.NewRecorder()

			// リクエストの実行
			router.ServeHTTP(w, req)

			// レスポンスの検証
			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if w.Code == http.StatusOK {
				var res response.SignUpRes
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)

				// CreatedAtとUpdatedAtは動的な値なので、比較をスキップ
				res.User.CreatedAt = time.Time{}
				res.User.UpdatedAt = time.Time{}

				assert.Equal(t, tc.expectedResponse, res)
			} else {
				var res render.Error
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, res)
			}
		})
	}
}

func TestUserCtrl_SignIn(t *testing.T) {
	logging.Init()

	tests := []struct {
		name               string
		inputBody          request.SignInReq
		mockSetup          func(*mock_usecase.MockUserUsecase)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "正常系: ユーザー認証成功",
			inputBody: request.SignInReq{
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					SignIn(gomock.Any(), gomock.Eq(request.SignInReq{
						Email:    "test@example.com",
						Password: "password",
					})).
					Return(&response.SignInRes{
						User: response.UserRes{
							Name:      "Test User",
							Email:     "test@example.com",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Token: response.TokenRes{
							AccessToken:  "access-token",
							RefreshToken: "refresh-token",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: response.SignInRes{
				User: response.UserRes{
					Name:      "Test User",
					Email:     "test@example.com",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				Token: response.TokenRes{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				},
			},
		},
		{
			name: "異常系: ユーザーが見つからない",
			inputBody: request.SignInReq{
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					SignIn(gomock.Any(), gomock.Eq(request.SignInReq{
						Email:    "test@example.com",
						Password: "password",
					})).
					Return((*response.SignInRes)(nil), entity.ErrUserNotFound)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: render.Error{
				Message: entity.ErrUserNotFound.Error(),
				Code:    entity.CodeUserNotFound,
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "異常系: パスワードが間違っている",
			inputBody: request.SignInReq{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					SignIn(gomock.Any(), gomock.Eq(request.SignInReq{
						Email:    "test@example.com",
						Password: "wrongpassword",
					})).
					Return((*response.SignInRes)(nil), entity.ErrPasswordIncorrect)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: render.Error{
				Message: entity.ErrPasswordIncorrect.Error(),
				Code:    entity.CodePasswordIncorrect,
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "異常系: 内部エラー",
			inputBody: request.SignInReq{
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					SignIn(gomock.Any(), gomock.Eq(request.SignInReq{
						Email:    "test@example.com",
						Password: "password",
					})).
					Return((*response.SignInRes)(nil), errors.New("internal server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: render.Error{
				Message: "internal server error",
				Status:  http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ロガーを手動で設定（テスト用にログ出力を無視）
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			slog.SetDefault(logger)

			// gomock コントローラの作成
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モックの生成
			userUsecaseMock := mock_usecase.NewMockUserUsecase(ctrl)

			// モックの設定
			tc.mockSetup(userUsecaseMock)

			// ルーターのセットアップ
			router := setupRouter(userUsecaseMock)

			// リクエストボディのエンコード
			bodyBytes, err := json.Marshal(tc.inputBody)
			assert.NoError(t, err)

			// テストリクエストの作成
			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signin", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// レスポンスの記録
			w := httptest.NewRecorder()

			// リクエストの実行
			router.ServeHTTP(w, req)

			// レスポンスの検証
			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if w.Code == http.StatusOK {
				var res response.SignInRes
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)

				// CreatedAtとUpdatedAtは動的な値なので、比較をスキップ
				res.User.CreatedAt = time.Time{}
				res.User.UpdatedAt = time.Time{}

				assert.Equal(t, tc.expectedResponse, res)
			} else {
				var res render.Error
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, res)
			}
		})
	}
}

func TestUserCtrl_RefreshToken(t *testing.T) {
	logging.Init()

	tests := []struct {
		name               string
		inputBody          request.RefreshTokenReq
		mockSetup          func(*mock_usecase.MockUserUsecase)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "正常系: トークンリフレッシュ成功",
			inputBody: request.RefreshTokenReq{
				RefreshToken: "refresh-token",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					RefreshToken(gomock.Any(), gomock.Eq(request.RefreshTokenReq{
						RefreshToken: "refresh-token",
					})).
					Return(&response.RefreshTokenRes{
						Token: response.TokenRes{
							AccessToken:  "new-access-token",
							RefreshToken: "new-refresh-token",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: response.RefreshTokenRes{
				Token: response.TokenRes{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				},
			},
		},
		{
			name: "異常系: トークンが無効",
			inputBody: request.RefreshTokenReq{
				RefreshToken: "invalid-refresh-token",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					RefreshToken(gomock.Any(), gomock.Eq(request.RefreshTokenReq{
						RefreshToken: "invalid-refresh-token",
					})).
					Return((*response.RefreshTokenRes)(nil), entity.ErrTokenInValid)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: render.Error{
				Message: entity.ErrTokenInValid.Error(),
				Code:    entity.CodeTokenInValid,
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "異常系: トークンのパースに失敗",
			inputBody: request.RefreshTokenReq{
				RefreshToken: "invalid-refresh-token",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					RefreshToken(gomock.Any(), gomock.Eq(request.RefreshTokenReq{
						RefreshToken: "invalid-refresh-token",
					})).
					Return((*response.RefreshTokenRes)(nil), entity.ErrFailedToParseClaims)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: render.Error{
				Message: entity.ErrFailedToParseClaims.Error(),
				Code:    entity.CodeFailedToParseClaims,
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "異常系: トークンバージョンが一致しない",
			inputBody: request.RefreshTokenReq{
				RefreshToken: "refresh-token",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					RefreshToken(gomock.Any(), gomock.Eq(request.RefreshTokenReq{
						RefreshToken: "refresh-token",
					})).
					Return((*response.RefreshTokenRes)(nil), entity.ErrTokenVersionMismatch)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: render.Error{
				Message: entity.ErrTokenVersionMismatch.Error(),
				Code:    entity.CodeTokenVersionMismatch,
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "異常系: 内部エラー",
			inputBody: request.RefreshTokenReq{
				RefreshToken: "refresh-token",
			},
			mockSetup: func(mock *mock_usecase.MockUserUsecase) {
				mock.EXPECT().
					RefreshToken(gomock.Any(), gomock.Eq(request.RefreshTokenReq{
						RefreshToken: "refresh-token",
					})).
					Return((*response.RefreshTokenRes)(nil), errors.New("internal server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: render.Error{
				Message: "internal server error",
				Status:  http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ロガーを手動で設定（テスト用にログ出力を無視）
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			slog.SetDefault(logger)

			// gomock コントローラの作成
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モックの生成
			userUsecaseMock := mock_usecase.NewMockUserUsecase(ctrl)

			// モックの設定
			tc.mockSetup(userUsecaseMock)

			// ルーターのセットアップ
			router := setupRouter(userUsecaseMock)

			// リクエストボディのエンコード
			bodyBytes, err := json.Marshal(tc.inputBody)
			assert.NoError(t, err)

			// テストリクエストの作成
			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// レスポンスの記録
			w := httptest.NewRecorder()

			// リクエストの実行
			router.ServeHTTP(w, req)

			// レスポンスの検証
			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if w.Code == http.StatusOK {
				var res response.RefreshTokenRes
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedResponse, res)
			} else {
				var res render.Error
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, res)
			}
		})
	}
}
