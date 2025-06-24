package service

import (
	"Notes/internal/constants"
	"Notes/internal/model"
	mocks "Notes/internal/service/mock"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

type userTestArgs struct {
	userId   int
	login    string
	password string
	name     string
	surname  string
}

type userTestExpect struct {
	id    int
	error *model.ApplicationError
	user  *model.User
}

func initUserServiceTest(t *testing.T) (AbstractUserService, *mocks.MockAbstractRepository, *mocks.MockAbstractHashService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockAbstractRepository(ctrl)
	mockHashService := mocks.NewMockAbstractHashService(ctrl)

	return NewConcreteUserService(mockRepository, mockHashService), mockRepository, mockHashService
}

func TestConcreteUserService_CreateUser(t *testing.T) {
	userService, repo, hash := initUserServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    userTestArgs
		want    userTestExpect
		wantErr bool
	}{
		{
			name: "empty name returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "",
				password: "",
				name:     "",
				surname:  "",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Имя не может быть пустым", nil),
			},
			wantErr: true,
		},
		{
			name: "empty surname returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "",
				password: "",
				name:     "name",
				surname:  "",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Фамилия не может быть пустой", nil),
			},
			wantErr: true,
		},
		{
			name: "short login returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "l",
				password: "",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, fmt.Sprintf("Логин слишком короткий. Пожалуйста, создайте логин длинной не меньше %d символов", model.MinLoginLength), nil),
			},
			wantErr: true,
		},
		{
			name: "short password returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "pas",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, fmt.Sprintf("Пароль слишком короткий. Пожалуйста, создайте пароль длиной не менее %d символов", model.MinPasswordLength), nil),
			},
			wantErr: true,
		},
		{
			name: "password without upper cases returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "passwordpasw",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Пароль должен содержать букву верхнего регистра.", nil),
			},
			wantErr: true,
		},
		{
			name: "password without lower cases returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "PASSWORDPASSS",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Пароль должен содержать букву нижнего регистра.", nil),
			},
			wantErr: true,
		},
		{
			name: "password without numbers returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Пароль должен содержать число.", nil),
			},
			wantErr: true,
		},
		{
			name: "password without special symbols returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss1",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Пароль должен содержать спецсимволы.", nil),
			},
			wantErr: true,
		},
		{
			name: "duplicate login returns error",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{
					{
						Id:        1,
						Name:      "name1234",
						Surname:   "surname4321",
						Login:     "login1234",
						Password:  "SecurePassword123$",
						Timestamp: time.Time{},
					},
				})
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeValidation, loginAlreadyUsedMessage, nil),
			},
			wantErr: true,
		},
		{
			name: "error while hash password returns error",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{})
				hash.EXPECT().GetHash("Passwordpasss123$").Return("", model.NewApplicationError(model.ErrorTypeInternal, "Ошибка при создании хэша", nil))
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeInternal, "Ошибка при создании хэша", nil),
			},
			wantErr: true,
		},
		{
			name: "error while saving user returns error",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{})
				hash.EXPECT().GetHash("Passwordpasss123$").Return("hashed_password", nil)
				repo.EXPECT().SaveEntity(&model.User{
					Id:       0,
					Name:     "name",
					Surname:  "surname",
					Login:    "login1234",
					Password: "hashed_password",
				}).Return(constants.FakeId, model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil))
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil),
			},
			wantErr: true,
		},
		{
			name: "save user happy path",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{})
				hash.EXPECT().GetHash("Passwordpasss123$").Return("hashed_password", nil)
				repo.EXPECT().SaveEntity(&model.User{
					Id:       0,
					Name:     "name",
					Surname:  "surname",
					Login:    "login1234",
					Password: "hashed_password",
				}).Return(1, nil)
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    1,
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := userService

			tt.mock()

			got, err := userService.CreateUser(tt.args.login, tt.args.password, tt.args.name, tt.args.surname)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want.id || (err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message)) {
				t.Errorf("userService.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcreteUserService_UpdateUser(t *testing.T) {
	userService, repo, hash := initUserServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    userTestArgs
		want    userTestExpect
		wantErr bool
	}{
		{
			name: "empty name returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "",
				password: "",
				name:     "",
				surname:  "",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Имя не может быть пустым", nil),
			},
			wantErr: true,
		},
		{
			name: "empty surname returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "",
				password: "",
				name:     "name",
				surname:  "",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Фамилия не может быть пустой", nil),
			},
			wantErr: true,
		},
		{
			name: "short login returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "l",
				password: "",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, fmt.Sprintf("Логин слишком короткий. Пожалуйста, создайте логин длинной не меньше %d символов", model.MinLoginLength), nil),
			},
			wantErr: true,
		},
		{
			name: "short password returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "pas",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, fmt.Sprintf("Пароль слишком короткий. Пожалуйста, создайте пароль длиной не менее %d символов", model.MinPasswordLength), nil),
			},
			wantErr: true,
		},
		{
			name: "password without upper cases returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "passwordpasw",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Пароль должен содержать букву верхнего регистра.", nil),
			},
			wantErr: true,
		},
		{
			name: "password without lower cases returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "PASSWORDPASSS",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Пароль должен содержать букву нижнего регистра.", nil),
			},
			wantErr: true,
		},
		{
			name: "password without numbers returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Пароль должен содержать число.", nil),
			},
			wantErr: true,
		},
		{
			name: "password without special symbols returns error",
			mock: func() {
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss1",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Пароль должен содержать спецсимволы.", nil),
			},
			wantErr: true,
		},
		{
			name: "duplicate login returns error",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{
					{
						Id:        1,
						Name:      "name1234",
						Surname:   "surname4321",
						Login:     "login1234",
						Password:  "SecurePassword123$",
						Timestamp: time.Time{},
					},
				})
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, loginAlreadyUsedMessage, nil),
			},
			wantErr: true,
		},
		{
			name: "error while hash password returns error",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{})
				hash.EXPECT().GetHash("Passwordpasss123$").Return("", model.NewApplicationError(model.ErrorTypeInternal, "Ошибка при создании хэша", nil))
			},
			args: userTestArgs{
				userId:   0,
				login:    "login1234",
				password: "Passwordpasss123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				id:    constants.FakeId,
				error: model.NewApplicationError(model.ErrorTypeInternal, "Ошибка при создании хэша", nil),
			},
			wantErr: true,
		},
		{
			name: "error unexisted user returns error",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{})
				hash.EXPECT().GetHash("Passwordpasss123$").Return("hashed_password", nil)
				repo.EXPECT().GetUserById(1).Return(&model.User{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			args: userTestArgs{
				userId:   1,
				login:    "login1234",
				password: "Passwordpasss123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "error while saving user returns error",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{})
				hash.EXPECT().GetHash("New_password123$").Return("New_hashed_password123$", nil)
				repo.EXPECT().GetUserById(1).Return(&model.User{
					Id:       1,
					Name:     "initial name",
					Surname:  "initial surname",
					Login:    "initial_login",
					Password: "initial_hashed_password123$",
				}, nil)
				repo.EXPECT().SaveEntity(&model.User{
					Id:       1,
					Name:     "name",
					Surname:  "surname",
					Login:    "new_login",
					Password: "New_hashed_password123$",
				}).Return(constants.FakeId, model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil))
			},
			args: userTestArgs{
				userId:   1,
				login:    "new_login",
				password: "New_password123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil),
			},
			wantErr: true,
		},
		{
			name: "update user happy path",
			mock: func() {
				repo.EXPECT().GetUsers().Return([]*model.User{})
				hash.EXPECT().GetHash("New_password123$").Return("New_hashed_password123$", nil)
				repo.EXPECT().GetUserById(1).Return(&model.User{
					Id:       1,
					Name:     "initial name",
					Surname:  "initial surname",
					Login:    "initial_login",
					Password: "initial_hashed_password123$",
				}, nil)
				repo.EXPECT().SaveEntity(&model.User{
					Id:       1,
					Name:     "name",
					Surname:  "surname",
					Login:    "new_login",
					Password: "New_hashed_password123$",
				}).Return(1, nil)
			},
			args: userTestArgs{
				userId:   1,
				login:    "new_login",
				password: "New_password123$",
				name:     "name",
				surname:  "surname",
			},
			want: userTestExpect{
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := userService

			tt.mock()

			err := userService.UpdateUser(tt.args.userId, tt.args.login, tt.args.password, tt.args.name, tt.args.surname)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("userService.CreateUser() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestConcreteUserService_GetUser(t *testing.T) {
	userService, repo, _ := initUserServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    userTestArgs
		want    userTestExpect
		wantErr bool
	}{
		{
			name: "attempt to get unexisted user",
			mock: func() {
				repo.EXPECT().GetUserById(1).Return(nil, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			args: userTestArgs{
				userId: 1,
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "get user happy path",
			mock: func() {
				repo.EXPECT().GetUserById(1).Return(&model.User{
					Id:       1,
					Name:     "name",
					Surname:  "surname",
					Login:    "login",
					Password: "password",
				}, nil)
			},
			args: userTestArgs{
				userId: 1,
			},
			want: userTestExpect{
				user: &model.User{
					Id:       1,
					Name:     "name",
					Surname:  "surname",
					Login:    "login",
					Password: "password",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := userService

			tt.mock()

			got, err := userService.GetUser(tt.args.userId)

			if (err != nil) != tt.wantErr {
				t.Errorf("NoteService.CreateNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotJson, _ := json.Marshal(got)
			expectedJson, _ := json.Marshal(tt.want.user)
			if fmt.Sprintf("%v", string(gotJson)) != fmt.Sprintf("%v", string(expectedJson)) {
				t.Errorf("userService.GetUser() = %v, want %v", string(gotJson), string(expectedJson))
			}
		})
	}
}

func TestConcreteUserService_DeleteUser(t *testing.T) {
	userService, repo, _ := initUserServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    userTestArgs
		want    userTestExpect
		wantErr bool
	}{
		{
			name: "attempt to delete unexisted user",
			mock: func() {
				repo.EXPECT().GetUserById(1).Return(&model.User{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			args: userTestArgs{
				userId: 1,
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "error while saving returns error",
			mock: func() {
				repo.EXPECT().GetUserById(1).Return(&model.User{
					Id:       1,
					Name:     "name",
					Surname:  "surname",
					Login:    "login",
					Password: "hashed_password",
				}, nil)
				repo.EXPECT().DeleteEntity(&model.User{
					Id:       1,
					Name:     "name",
					Surname:  "surname",
					Login:    "login",
					Password: "hashed_password",
				}).Return(model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil))
			},
			args: userTestArgs{
				userId: 1,
			},
			want: userTestExpect{
				error: model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil),
			},
			wantErr: true,
		},
		{
			name: "delete user happy path",
			mock: func() {
				repo.EXPECT().GetUserById(1).Return(&model.User{
					Id:       1,
					Name:     "name",
					Surname:  "surname",
					Login:    "login",
					Password: "hashed_password",
				}, nil)
				repo.EXPECT().DeleteEntity(&model.User{
					Id:       1,
					Name:     "name",
					Surname:  "surname",
					Login:    "login",
					Password: "hashed_password",
				}).Return(nil)
			},
			args: userTestArgs{
				userId: 1,
			},
			want: userTestExpect{
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := userService

			tt.mock()

			err := userService.DeleteUser(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("userService.DeleteUser() = %v, want %v", err, tt.want)
			}
		})
	}
}
