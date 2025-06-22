package service

import (
	"Notes/internal/model"
	mocks "Notes/internal/service/mock"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func initNotebookServiceTest(t *testing.T) (AbstractNotebookService, *mocks.MockAbstractRepository) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockAbstractRepository(ctrl)

	return NewConcreteNotebookService(mockRepository), mockRepository
}

func TestConcreteNotebookService_GetUserNotebook(t *testing.T) {
	notebookService, repo := initNotebookServiceTest(t)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		mock    func()
		args    int
		want    model.Notebook
		wantErr bool
	}{
		{
			name: "no notes and folders",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{})
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{})
			},
			args: 1,
			want: model.Notebook{
				Folders: []model.Folder{},
				Notes:   []model.Note{},
			},
			wantErr: false,
		},
		{
			name: "no notes and one folder",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:        1,
						Title:     "title",
						Timestamp: fixedTime,
						UserId:    1,
						Notes:     nil,
					},
				})
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{})
			},
			args: 1,
			want: model.Notebook{
				Folders: []model.Folder{
					{
						Id:        1,
						Title:     "title",
						Timestamp: fixedTime,
						UserId:    1,
						Notes:     []model.Note{},
					},
				},
				Notes: []model.Note{},
			},
			wantErr: false,
		},
		{
			name: "one folder and one note without folder",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:        1,
						Title:     "title",
						Timestamp: fixedTime,
						UserId:    1,
						Notes:     nil,
					},
				})
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:         1,
						Title:      "title",
						Content:    "content",
						UserId:     1,
						IsFavorite: false,
						Timestamp:  fixedTime,
						Tags:       nil,
					},
				})
			},
			args: 1,
			want: model.Notebook{
				Folders: []model.Folder{
					{
						Id:        1,
						Title:     "title",
						Timestamp: fixedTime,
						UserId:    1,
						Notes:     []model.Note{},
					},
				},
				Notes: []model.Note{
					{
						Id:         1,
						Title:      "title",
						Content:    "content",
						UserId:     1,
						IsFavorite: false,
						Timestamp:  fixedTime,
						Tags:       nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "one folder and one note in folder",
			mock: func() {
				repo.EXPECT().GetFoldersByUserId(1).Return([]*model.Folder{
					{
						Id:        1,
						Title:     "title",
						Timestamp: fixedTime,
						UserId:    1,
						Notes:     nil,
					},
				})
				note := &model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
					Timestamp:  fixedTime,
					Tags:       nil,
				}
				folderId := 1
				note.FolderId = &folderId
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					note,
				})

			},
			args: 1,
			want: model.Notebook{
				Folders: []model.Folder{
					{
						Id:        1,
						Title:     "title",
						Timestamp: fixedTime,
						UserId:    1,
						Notes: []model.Note{
							{
								Id:         1,
								Title:      "title",
								Content:    "content",
								UserId:     1,
								IsFavorite: false,
								Timestamp:  fixedTime,
								Tags:       nil,
							},
						},
					},
				},
				Notes: []model.Note{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notebookService := notebookService

			tt.mock()

			got := notebookService.GetUserNotebook(tt.args)

			gotJson, _ := json.Marshal(got)
			expectedJson, _ := json.Marshal(tt.want)
			if fmt.Sprintf("%v", string(gotJson)) != fmt.Sprintf("%v", string(expectedJson)) {
				t.Errorf("notebookService.GetUserNotebook() = %v, want %v", string(gotJson), string(expectedJson))
			}
		})
	}
}
