package service

import (
	"Notes/internal/model"
	mocks "Notes/internal/service/mock"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

type folderTestArgs struct {
	userId   int
	folderId int
	title    string
}

type folderTestExpect struct {
	id    int
	error *model.ApplicationError
}

func initFolderServiceTest(t *testing.T) (AbstractFolderService, *mocks.MockAbstractRepository) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockAbstractRepository(ctrl)

	return NewConcreteFolderService(mockRepository), mockRepository
}

func TestConcreteFolderService_CreateFolder(t *testing.T) {
	folderService, repo := initFolderServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    folderTestArgs
		want    folderTestExpect
		wantErr bool
	}{
		{
			name: "invalid folder name",
			mock: func() {
			},
			args: folderTestArgs{
				userId: 0,
				title:  "",
			},
			want: folderTestExpect{
				id:    -1,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Название папки не может быть пустым", nil),
			},
			wantErr: true,
		},
		{
			name: "folder with duplicate title",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:     0,
						Title:  "duplicate title",
						UserId: 1,
						Notes:  nil,
					},
				})
			},
			args: folderTestArgs{
				userId: 1,
				title:  "duplicate title",
			},
			want: folderTestExpect{
				id:    -1,
				error: model.NewApplicationError(model.ErrorTypeValidation, folderTitleIsNotFree, nil),
			},
			wantErr: true,
		},
		{
			name: "error while saving returns error",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:     1,
						Title:  "title",
						UserId: 1,
						Notes:  nil,
					},
				})
				repo.EXPECT().SaveEntity(&model.Folder{
					Id:     0,
					Title:  "original title",
					UserId: 1,
					Notes:  nil,
				}).Return(fakeId, model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil))
			},
			args: folderTestArgs{
				userId: 1,
				title:  "original title",
			},
			want: folderTestExpect{
				id:    fakeId,
				error: model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil),
			},
			wantErr: true,
		},
		{
			name: "save folder happy path",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:     1,
						Title:  "title",
						UserId: 1,
						Notes:  nil,
					},
				})
				repo.EXPECT().SaveEntity(&model.Folder{
					Id:     0,
					Title:  "original title",
					UserId: 1,
					Notes:  nil,
				}).Return(2, nil)
			},
			args: folderTestArgs{
				userId: 1,
				title:  "original title",
			},
			want: folderTestExpect{
				id:    2,
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folderService := folderService

			tt.mock()

			got, err := folderService.CreateFolder(tt.args.userId, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("FolderService.CreateFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want.id || (err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message)) {
				t.Errorf("FolderService.CreateFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcreteFolderService_UpdateFolder(t *testing.T) {
	folderService, repo := initFolderServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    folderTestArgs
		want    folderTestExpect
		wantErr bool
	}{
		{
			name: "invalid folder name",
			mock: func() {
			},
			args: folderTestArgs{
				userId: 0,
				title:  "",
			},
			want: folderTestExpect{
				id:    -1,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Название папки не может быть пустым", nil),
			},
			wantErr: true,
		},
		{
			name: "folder with duplicate title",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:     1,
						Title:  "title",
						UserId: 1,
						Notes:  nil,
					},
					{
						Id:     1,
						Title:  "duplicate title",
						UserId: 1,
						Notes:  nil,
					},
				})
			},
			args: folderTestArgs{
				userId: 1,
				title:  "duplicate title",
			},
			want: folderTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, folderTitleIsNotFree, nil),
			},
			wantErr: true,
		},
		{
			name: "not existed folder",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:     1,
						Title:  "title",
						UserId: 1,
						Notes:  nil,
					},
				})
				repo.EXPECT().GetFolderById(2, 1).Return(&model.Folder{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			args: folderTestArgs{
				userId:   1,
				title:    "new title",
				folderId: 2,
			},
			want: folderTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "error while saving returns error",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:     1,
						Title:  "title",
						UserId: 1,
						Notes:  nil,
					},
				})
				repo.EXPECT().GetFolderById(1, 1).Return(&model.Folder{
					Id:     1,
					Title:  "title",
					UserId: 1,
					Notes:  nil,
				}, nil)
				repo.EXPECT().SaveEntity(&model.Folder{
					Id:     1,
					Title:  "original title",
					UserId: 1,
					Notes:  nil,
				}).Return(fakeId, model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil))
			},
			args: folderTestArgs{
				userId:   1,
				title:    "original title",
				folderId: 1,
			},
			want: folderTestExpect{
				error: model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil),
			},
			wantErr: true,
		},
		{
			name: "update folder happy path",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:     1,
						Title:  "title",
						UserId: 1,
						Notes:  nil,
					},
				})
				repo.EXPECT().GetFolderById(1, 1).Return(&model.Folder{
					Id:     1,
					Title:  "title",
					UserId: 1,
					Notes:  nil,
				}, nil)
				repo.EXPECT().SaveEntity(&model.Folder{
					Id:     1,
					Title:  "original title",
					UserId: 1,
					Notes:  nil,
				}).Return(1, nil)
			},
			args: folderTestArgs{
				userId:   1,
				title:    "original title",
				folderId: 1,
			},
			want: folderTestExpect{
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folderService := folderService

			tt.mock()

			err := folderService.UpdateFolder(tt.args.userId, tt.args.folderId, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("FolderService.CreateFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("FolderService.CreateFolder() unexpected error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestConcreteFolderService_DeleteFolder(t *testing.T) {
	folderService, repo := initFolderServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    folderTestArgs
		want    folderTestExpect
		wantErr bool
	}{
		{
			name: "no folder with such id",
			mock: func() {
				repo.EXPECT().GetFolderById(2, 1).Return(&model.Folder{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			args: folderTestArgs{
				userId:   1,
				folderId: 2,
			},
			want:    folderTestExpect{},
			wantErr: false,
		},
		{
			name: "no folder with such id belongs to user",
			mock: func() {
				repo.EXPECT().GetFolderById(2, 1).Return(&model.Folder{}, model.NewApplicationError(model.ErrorTypeNotFound, "Папка не найдена", nil))
			},
			args: folderTestArgs{
				userId:   1,
				folderId: 2,
			},
			want:    folderTestExpect{},
			wantErr: false,
		},
		{
			name: "folder deleted",
			mock: func() {
				repo.EXPECT().GetFolderById(2, 1).Return(&model.Folder{
					Id:        2,
					Title:     "title",
					Timestamp: time.Time{},
					UserId:    1,
					Notes:     nil,
				}, nil)
				repo.EXPECT().DeleteEntity(&model.Folder{
					Id:        2,
					Title:     "title",
					Timestamp: time.Time{},
					UserId:    1,
					Notes:     nil,
				}).Return(nil)
			},
			args: folderTestArgs{
				userId:   1,
				folderId: 2,
			},
			want:    folderTestExpect{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folderService := folderService

			tt.mock()

			err := folderService.DeleteFolder(tt.args.userId, tt.args.folderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FolderService.CreateFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("FolderService.CreateFolder() unexpected error = %v, want %v", err, tt.want)
			}
		})
	}
}
