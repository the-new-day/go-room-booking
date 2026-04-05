package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthUseCase_Register(t *testing.T) {
	const accessTokenTTL = 30 * time.Minute

	mockErr := errors.New("aaa")
	mockUserID := uuid.New()

	tests := []struct {
		name     string
		email    string
		password string

		wantErr error

		mockPasswordHash      string
		mockPasswordHasherErr error

		mockRepoErr error
	}{
		{
			name:             "success",
			email:            "example@gmail.com",
			password:         "abcde",
			mockPasswordHash: "123",
		},
		{
			name:             "email already registered",
			email:            "example@gmail.com",
			password:         "abcde",
			mockPasswordHash: "123",
			mockRepoErr:      storage.ErrAlreadyExists,
			wantErr:          domain.ErrUserWithEmailAlreadyExists,
		},
		{
			name:             "unknown storage error",
			email:            "example@gmail.com",
			password:         "abcde",
			mockPasswordHash: "123",
			mockRepoErr:      mockErr,
			wantErr:          mockErr,
		},
		{
			name:                  "hashing error",
			email:                 "example@gmail.com",
			password:              "abcde",
			mockPasswordHash:      "123",
			mockPasswordHasherErr: mockErr,
			wantErr:               ErrFailedToHashPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockHasher := mocks.NewMockPasswordHasher(t)
			mockHasher.EXPECT().
				Hash(tt.password).
				Return(tt.mockPasswordHash, tt.mockPasswordHasherErr)

			mockUser := &entity.User{
				UserID:       mockUserID,
				Role:         DefaultRoleUponRegistration,
				Email:        tt.email,
				PasswordHash: tt.mockPasswordHash,
			}
			mockRepo := mocks.NewMockUserRepository(t)
			call := mockRepo.EXPECT().
				Create(t.Context(), tt.email, tt.mockPasswordHash, DefaultRoleUponRegistration).
				Return(mockUser, tt.mockRepoErr)

			if tt.mockPasswordHasherErr == nil {
				call.Once()
			} else {
				call.Unset()
			}

			mockTokenManager := mocks.NewMockTokenManager(t)

			uc := New(mockRepo, mockTokenManager, mockHasher, accessTokenTTL)

			user, err := uc.Register(t.Context(), tt.email, tt.password)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.Equal(t, user.UserID, mockUserID)
				assert.Equal(t, user.Email, tt.email)
				assert.Equal(t, user.PasswordHash, tt.mockPasswordHash)
			}

			mockHasher.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
			mockTokenManager.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	const accessTokenTTL = 30 * time.Minute

	mockErr := errors.New("aaa")
	mockUserID := uuid.New()

	tests := []struct {
		name     string
		email    string
		password string

		wantToken string
		wantErr   error

		userPasswordHash string
		passwordMatches  bool

		mockTokenManagerErr error
		mockRepoErr         error
	}{
		{
			name:            "success",
			email:           "example@gmail.com",
			password:        "abcde",
			wantToken:       "1.2.3",
			passwordMatches: true,
		},
		{
			name:            "invalid password",
			email:           "example@gmail.com",
			password:        "abcde",
			passwordMatches: false,
			wantErr:         domain.ErrInvalidPassword,
		},
		{
			name:        "email not found",
			email:       "example@gmail.com",
			password:    "abcde",
			mockRepoErr: storage.ErrNotFound,
			wantErr:     domain.ErrEmailNotFound,
		},
		{
			name:        "unknown storage error",
			email:       "example@gmail.com",
			password:    "abcde",
			mockRepoErr: mockErr,
			wantErr:     mockErr,
		},
		{
			name:                "token manager error",
			email:               "example@gmail.com",
			password:            "abcde",
			passwordMatches:     true,
			mockTokenManagerErr: mockErr,
			wantErr:             mockErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUser := &entity.User{
				UserID:       mockUserID,
				Role:         DefaultRoleUponRegistration,
				Email:        tt.email,
				PasswordHash: tt.userPasswordHash,
			}
			mockRepo := mocks.NewMockUserRepository(t)
			mockRepo.EXPECT().
				FindByEmail(t.Context(), tt.email).
				Return(mockUser, tt.mockRepoErr).
				Once()

			mockHasher := mocks.NewMockPasswordHasher(t)
			hasherCall := mockHasher.EXPECT().
				DoesMatch(tt.userPasswordHash, tt.password).
				Return(tt.passwordMatches)

			mockTokenManager := mocks.NewMockTokenManager(t)
			tmCall := mockTokenManager.EXPECT().
				CreateToken(mockUserID.String(), string(DefaultRoleUponRegistration), accessTokenTTL).
				Return(tt.wantToken, tt.mockTokenManagerErr)

			if tt.passwordMatches {
				tmCall.Once()
			} else {
				tmCall.Unset()
			}

			if tt.mockRepoErr == nil {
				hasherCall.Once()
			} else {
				hasherCall.Unset()
			}

			uc := New(mockRepo, mockTokenManager, mockHasher, accessTokenTTL)

			token, err := uc.Login(t.Context(), tt.email, tt.password)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}

			mockHasher.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
			mockTokenManager.AssertExpectations(t)
		})
	}
}
