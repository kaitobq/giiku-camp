package entity

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		email    string
		password string
		wantErr  error
	}{
		{
			name:     "正常系",
			userName: "test user",
			email:    "test@example.com",
			password: "password123",
			wantErr:  nil,
		},
		{
			name:     "異常系: 無効なメールアドレス",
			userName: "test user",
			email:    "invalid-email",
			password: "password123",
			wantErr:  ErrEmailInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.userName, tt.email, tt.password)
			if err != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.Name != tt.userName {
					t.Errorf("NewUser().Name = %v, want %v", got.Name, tt.userName)
				}
				if got.Email != tt.email {
					t.Errorf("NewUser().Email = %v, want %v", got.Email, tt.email)
				}
				if got.ID == "" {
					t.Error("NewUser().ID should not be empty")
				}
				if got.TokenVersion != 1 {
					t.Errorf("NewUser().TokenVersion = %v, want 1", got.TokenVersion)
				}
			}
		})
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	password := "password123"
	user, err := NewUser("test user", "test@example.com", password)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "正常系",
			password: password,
			wantErr:  nil,
		},
		{
			name:     "異常系: パスワード不一致",
			password: "wrongpassword",
			wantErr:  ErrPasswordIncorrect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := user.VerifyPassword(tt.password); err != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_HidePassword(t *testing.T) {
	user, err := NewUser("test user", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	hiddenUser := user.HidePassword()
	if hiddenUser.Password != "***" {
		t.Errorf("HidePassword() = %v, want ***", hiddenUser.Password)
	}
	if user.Password == "***" {
		t.Error("Original user password should not be modified")
	}
}

func TestUser_UpdateUpdatedAt(t *testing.T) {
	user, err := NewUser("test user", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.UpdatedAt = time.Now()

	lastUpdatedAt := user.UpdatedAt
	time.Sleep(1 * time.Second)
	user.UpdateUpdatedAt()
	if user.UpdatedAt == lastUpdatedAt {
		t.Error("UpdateUpdatedAt() failed to update UpdatedAt")
	}
}

func TestUser_UpdateCreatedAt(t *testing.T) {
	user, err := NewUser("test user", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.CreatedAt = time.Now()

	lastCreatedAt := user.CreatedAt
	time.Sleep(1 * time.Second)
	user.UpdateCreatedAt()
	if user.CreatedAt == lastCreatedAt {
		t.Error("UpdateCreatedAt() failed to update CreatedAt")
	}
}

func TestUser_IncrementTokenVersion(t *testing.T) {
	user, err := NewUser("test user", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	lastTokenVersion := user.TokenVersion
	user.IncrementTokenVersion()

	if user.TokenVersion != lastTokenVersion+1 {
		t.Errorf("IncrementTokenVersion() = %v, want %v", user.TokenVersion, lastTokenVersion+1)
	}
}
