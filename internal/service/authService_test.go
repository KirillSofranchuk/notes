package service

import (
	"Notes/internal/model"
	mocks "Notes/internal/service/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

type authTestArgs struct {
	login    string
	password string
}

func initAuthServiceTests(t *testing.T) (AbstractAuthService, *mocks.MockAbstractRepository, *mocks.MockAbstractJwtService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockAbstractRepository(ctrl)
	mockJwtService := mocks.NewMockAbstractJwtService(ctrl)

	return NewConcreteAuthService(mockRepository, mockJwtService), mockRepository, mockJwtService
}

func TestConcreteAuthService_AuthUser(t *testing.T) {
	authService, repo, jwtService := initAuthServiceTests(t)

	tests := []struct {
		name    string
		mock    func()
		args    authTestArgs
		want    string
		wantErr bool
	}{
		{
			name: "user not found",
			mock: func() {
				repo.EXPECT().GetUser("login", "password").Return(nil, model.NewApplicationError(model.ErrorTypeNotFound, "Пользователь не найден", nil))
			},
			args: authTestArgs{
				login:    "login",
				password: "password",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "user found and token not valid",
			mock: func() {
				repo.EXPECT().GetUser("login", "password").Return(&model.User{
					Login:    "login",
					Password: "password",
					Id:       1,
				}, nil)
				jwtService.EXPECT().GetToken(1).Return("", model.NewApplicationError(model.ErrorTypeInternal, "Ошибка при формировании токена", nil))
			},
			args: authTestArgs{
				login:    "login",
				password: "password",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "user found and token valid",
			mock: func() {
				repo.EXPECT().GetUser("login", "password").Return(&model.User{
					Login:    "login",
					Password: "password",
					Id:       1,
				}, nil)
				jwtService.EXPECT().GetToken(1).Return("valid token", nil)
			},
			args: authTestArgs{
				login:    "login",
				password: "password",
			},
			want:    "valid token",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := authService

			tt.mock()

			got, err := authService.AuthUser(tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.AuthUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("AuthService.AuthUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcreteAuthService_ValidateToken(t *testing.T) {
	authService, repo, jwtService := initAuthServiceTests(t)

	tests := []struct {
		name    string
		mock    func()
		args    string
		want    *model.Claims
		wantErr bool
	}{
		{
			name: "invalid token signature",
			mock: func() {
				jwtService.EXPECT().ParseToken("token with invalid signature").Return(nil, model.NewApplicationError(model.ErrorTypeAuth, "Неверный метод подписи", nil))
			},
			args:    "token with invalid signature",
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid token",
			mock: func() {
				jwtService.EXPECT().ParseToken("invalid token").Return(nil, model.NewApplicationError(model.ErrorTypeAuth, "Невалидный токен", nil))
			},
			args:    "invalid token",
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid token no user",
			mock: func() {
				jwtService.EXPECT().ParseToken("valid token").Return(&model.Claims{UserId: 1}, nil)
				repo.EXPECT().GetUserById(1).Return(nil, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			args:    "valid token",
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid token user exists",
			mock: func() {
				jwtService.EXPECT().ParseToken("valid token").Return(&model.Claims{UserId: 1}, nil)
				repo.EXPECT().GetUserById(1).Return(&model.User{Id: 1}, nil)
			},
			args:    "valid token",
			want:    &model.Claims{UserId: 1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := authService

			tt.mock()

			got, err := authService.ValidateToken(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got == nil || tt.want == nil {
				if got != tt.want {
					t.Errorf("AuthService.ValidateToken() = %v, want %v", got, tt.want)
				}
				return
			}

			if got.UserId != tt.want.UserId {
				t.Errorf("AuthService.ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
