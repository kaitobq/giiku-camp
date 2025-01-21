package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"giiku-camp/internal/infra/logging"
	mock_usecase "giiku-camp/internal/mock/usecase"
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserCrossingCtrl_RegisterUserCrossing(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mock_usecase.MockUserCrossingUsecase)
		expectedStatus int
	}{
		{
			name:        "正常系 - ユーザーIDが1つ",
			requestBody: `{"user_ids": ["user1"]}`,
			mockSetup: func(mock *mock_usecase.MockUserCrossingUsecase) {
				mock.EXPECT().RegisterUserCrossing(gomock.Any(), request.RegisterUserCrossingReq{
					UserIDs: []string{"user1"},
				}).Return(&response.RegisterUserCrossingRes{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "異常系 - ユーザーIDが空",
			requestBody: `{"user_ids": []}`,
			mockSetup: func(mock *mock_usecase.MockUserCrossingUsecase) {
				// モックは呼ばれない
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "異常系 - リクエストボディが不正",
			requestBody: `invalid json`,
			mockSetup: func(mock *mock_usecase.MockUserCrossingUsecase) {
				// モックは呼ばれない
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "異常系 - ユースケースでエラー発生",
			requestBody: `{"user_ids": ["user1"]}`,
			mockSetup: func(mock *mock_usecase.MockUserCrossingUsecase) {
				mock.EXPECT().RegisterUserCrossing(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logging.Init()
			// モックの準備
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_usecase.NewMockUserCrossingUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			// テスト対象のコントローラ作成
			ctrlr := NewUserCrossingCtrl(mockUsecase)

			// テスト用のGinコンテキスト作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/authenticated/crossing", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// テスト実行
			ctrlr.RegisterUserCrossing(c)

			// 結果検証
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUserCrossingCtrl_GetUserCrossing(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mock_usecase.MockUserCrossingUsecase)
		expectedStatus int
	}{
		{
			name: "正常系",
			mockSetup: func(mock *mock_usecase.MockUserCrossingUsecase) {
				mock.EXPECT().GetUserCrossing(gomock.Any()).
					Return(&response.GetUserCrossingRes{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "異常系 - ユースケースでエラー発生",
			mockSetup: func(mock *mock_usecase.MockUserCrossingUsecase) {
				mock.EXPECT().GetUserCrossing(gomock.Any()).
					Return(nil, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logging.Init()
			// モックの準備
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_usecase.NewMockUserCrossingUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			// テスト対象のコントローラ作成
			ctrlr := NewUserCrossingCtrl(mockUsecase)

			// テスト用のGinコンテキスト作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/v1/authenticated/crossing", nil)

			// テスト実行
			ctrlr.GetUserCrossing(c)

		})
	}
}
