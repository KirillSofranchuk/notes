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

type noteTestArgs struct {
	userId   int
	title    string
	content  string
	tags     *[]string
	noteId   int
	folderId *int
	query    string
}

type noteTestExpect struct {
	id    int
	error *model.ApplicationError
	notes []*model.Note
}

func initNoteServiceTest(t *testing.T) (AbstractNoteService, *mocks.MockAbstractRepository) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockAbstractRepository(ctrl)

	return NewConcreteNoteService(mockRepository), mockRepository
}

func TestConcreteNoteService_CreateNote(t *testing.T) {
	noteService, repo := initNoteServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    noteTestArgs
		want    noteTestExpect
		wantErr bool
	}{
		{
			name: "note empty name",
			args: noteTestArgs{
				userId:  1,
				title:   "",
				content: "content",
				tags:    nil,
			},
			mock: func() {},
			want: noteTestExpect{
				id:    -1,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Название заметки не может быть пустым", nil),
			},
			wantErr: true,
		},
		{
			name: "note empty content",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "",
				tags:    nil,
			},
			mock: func() {},
			want: noteTestExpect{
				id:    -1,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Заметка не может быть пустой", nil),
			},
			wantErr: true,
		},
		{
			name: "note large content",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "7fHj2#pL9qR!5zXn8*Kb0$cT1%gY6@mW3^dS4&vE7)hU8(oI9-Pj0+Qk1=Al2~Bs3{Ct4}De5|Fg6Gh7/Mn8?Na9<Ob0>Pc1,Ze2.Xf3Yh4Jk5Vl6Am7Bn8Co9Dp0Eq1Fr2Gs3Ht4Iu5Jv6Kw7Lx8My9Nz0Oa1Pb2Qc3Rd4Se5Tf6Ug7Vh8Wi9Xj0Yk1Zl2$m3#n4%o5&p6*q7(r8)s9-t0_u1=v2+w3~x4{y5}z6|A7B8/C9?D0<E1>F2,G3.H4I5J6K7L8M9N0O1P2Q3R4S5T6U7V8W9X0Y1Z2a3b4c5d6e7f8g9h0i1j2k3l4m5n6o7p8q9r0s1t2u3v4w5x6y7z8A9B0C1D2E3F4G5H6I7J8K9L0M1N2O3P4Q5R6S7T8U9V0W1X2Y3Z4",
				tags:    nil,
			},
			mock: func() {},
			want: noteTestExpect{
				id:    -1,
				error: model.NewApplicationError(model.ErrorTypeValidation, fmt.Sprintf("Длина заметки не может превышать %d символов", model.MaxContentLength), nil),
			},
			wantErr: true,
		},
		{
			name: "note invalid tags count",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "content",
				tags:    &[]string{"tag1", "tag2", "tag3", "tag4"},
			},
			mock: func() {},
			want: noteTestExpect{
				id:    -1,
				error: model.NewApplicationError(model.ErrorTypeValidation, fmt.Sprintf("Нельзя добавить больше, чем %d тегов к заметке.", model.MaxTagsCount), nil),
			},
			wantErr: true,
		},
		{
			name: "note duplicate title",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "content",
				tags:    nil,
			},
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:         1,
						Title:      "title",
						Content:    "content",
						UserId:     1,
						IsFavorite: false,
						Timestamp:  time.Time{},
					},
				})
			},
			want: noteTestExpect{
				id:    -1,
				error: model.NewApplicationError(model.ErrorTypeValidation, "Заметка с таким названием уже добавлена", nil),
			},
			wantErr: true,
		},
		{
			name: "note saved",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "content",
				tags:    nil,
			},
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:         1,
						Title:      "title2",
						Content:    "content",
						UserId:     1,
						IsFavorite: false,
						Timestamp:  time.Time{},
					},
				})
				repo.EXPECT().SaveEntity(&model.Note{
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
				}).Return(2, nil)
			},
			want: noteTestExpect{
				id:    2,
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteService := noteService

			tt.mock()

			got, err := noteService.CreateNote(tt.args.userId, tt.args.title, tt.args.content, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteService.CreateNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want.id || (err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message)) {
				t.Errorf("NoteService.CreateNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcreteNoteService_UpdateNote(t *testing.T) {
	noteService, repo := initNoteServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    noteTestArgs
		want    noteTestExpect
		wantErr bool
	}{
		{
			name: "note empty name",
			args: noteTestArgs{
				userId:  1,
				title:   "",
				content: "content",
				tags:    nil,
			},
			mock: func() {},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Название заметки не может быть пустым", nil),
			},
			wantErr: true,
		},
		{
			name: "note empty content",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "",
				tags:    nil,
			},
			mock: func() {},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Заметка не может быть пустой", nil),
			},
			wantErr: true,
		},
		{
			name: "note large content",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "7fHj2#pL9qR!5zXn8*Kb0$cT1%gY6@mW3^dS4&vE7)hU8(oI9-Pj0+Qk1=Al2~Bs3{Ct4}De5|Fg6Gh7/Mn8?Na9<Ob0>Pc1,Ze2.Xf3Yh4Jk5Vl6Am7Bn8Co9Dp0Eq1Fr2Gs3Ht4Iu5Jv6Kw7Lx8My9Nz0Oa1Pb2Qc3Rd4Se5Tf6Ug7Vh8Wi9Xj0Yk1Zl2$m3#n4%o5&p6*q7(r8)s9-t0_u1=v2+w3~x4{y5}z6|A7B8/C9?D0<E1>F2,G3.H4I5J6K7L8M9N0O1P2Q3R4S5T6U7V8W9X0Y1Z2a3b4c5d6e7f8g9h0i1j2k3l4m5n6o7p8q9r0s1t2u3v4w5x6y7z8A9B0C1D2E3F4G5H6I7J8K9L0M1N2O3P4Q5R6S7T8U9V0W1X2Y3Z4",
				tags:    nil,
			},
			mock: func() {},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, fmt.Sprintf("Длина заметки не может превышать %d символов", model.MaxContentLength), nil),
			},
			wantErr: true,
		},
		{
			name: "note invalid tags count",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "content",
				tags:    &[]string{"tag1", "tag2", "tag3", "tag4"},
			},
			mock: func() {},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, fmt.Sprintf("Нельзя добавить больше, чем %d тегов к заметке.", model.MaxTagsCount), nil),
			},
			wantErr: true,
		},
		{
			name: "note duplicate title",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "content",
				tags:    nil,
			},
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:         1,
						Title:      "title",
						Content:    "content",
						UserId:     1,
						IsFavorite: false,
						Timestamp:  time.Time{},
					},
				})
			},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeValidation, "Заметка с таким названием уже добавлена", nil),
			},
			wantErr: true,
		},
		{
			name: "note not found",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "content",
				tags:    nil,
				noteId:  2,
			},
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:         1,
						Title:      "title2",
						Content:    "content",
						UserId:     1,
						IsFavorite: false,
						Timestamp:  time.Time{},
					},
				})
				repo.EXPECT().GetNoteById(2, 1).Return(&model.Note{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "note updated",
			args: noteTestArgs{
				userId:  1,
				title:   "title",
				content: "content",
				tags:    nil,
				noteId:  2,
			},
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:         1,
						Title:      "title2",
						Content:    "content",
						UserId:     1,
						IsFavorite: false,
					},
					{
						Id:         2,
						Title:      "initial title",
						Content:    "initial content",
						UserId:     1,
						IsFavorite: false,
					},
				})
				repo.EXPECT().GetNoteById(2, 1).Return(&model.Note{
					Id:         2,
					Title:      "initial title",
					Content:    "initial content",
					UserId:     1,
					IsFavorite: false,
				}, nil)
				repo.EXPECT().SaveEntity(&model.Note{
					Id:      2,
					Title:   "title",
					Content: "content",
					UserId:  1,
				}).Return(2, nil)
			},
			want: noteTestExpect{
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteService := noteService

			tt.mock()

			err := noteService.UpdateNote(tt.args.userId, tt.args.noteId, tt.args.title, tt.args.content, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteService.CreateNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("NoteService.CreateNote() unexpected error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestConcreteNoteService_DeleteNote(t *testing.T) {
	noteService, repo := initNoteServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    noteTestArgs
		want    noteTestExpect
		wantErr bool
	}{
		{
			name: "attempt to delete unexisted note don't return error",
			args: noteTestArgs{
				userId: 1,
				noteId: 1,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			want: noteTestExpect{
				error: nil,
			},
			wantErr: false,
		},
		{
			name: "delete note happy path",
			args: noteTestArgs{
				userId: 1,
				noteId: 1,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
					Timestamp:  time.Time{},
				}, nil)
				repo.EXPECT().DeleteEntity(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
					Timestamp:  time.Time{},
				}).Return(nil)
			},
			want: noteTestExpect{
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteService := noteService

			tt.mock()

			err := noteService.DeleteNote(tt.args.userId, tt.args.noteId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteService.DeleteNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("NoteService.DeleteNote() unexpected error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestConcreteNoteService_MoveToFolder(t *testing.T) {
	noteService, repo := initNoteServiceTest(t)

	testFolderId := 2

	tests := []struct {
		name    string
		mock    func()
		args    noteTestArgs
		want    noteTestExpect
		wantErr bool
	}{
		{
			name: "attempt to move unexisted note returns error",
			args: noteTestArgs{
				userId: 1,
				noteId: 1,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "attempt to move note to unexisted folder returns error",
			args: noteTestArgs{
				userId:   1,
				noteId:   1,
				folderId: &testFolderId,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
				}, nil)
				repo.EXPECT().GetFolderById(testFolderId, 1).Return(&model.Folder{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "error while save to db returns error",
			args: noteTestArgs{
				userId:   1,
				noteId:   1,
				folderId: &testFolderId,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
				}, nil)
				repo.EXPECT().GetFolderById(testFolderId, 1).Return(&model.Folder{
					Id:     2,
					Title:  "folder",
					UserId: 1,
				}, nil)
				repo.EXPECT().SaveEntity(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
					FolderId:   &testFolderId,
				}).Return(fakeId, model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil))
			},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil),
			},
			wantErr: true,
		},
		{
			name: "move note to folder",
			args: noteTestArgs{
				userId:   1,
				noteId:   1,
				folderId: &testFolderId,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
				}, nil)
				repo.EXPECT().GetFolderById(testFolderId, 1).Return(&model.Folder{
					Id:     2,
					Title:  "folder",
					UserId: 1,
				}, nil)
				repo.EXPECT().SaveEntity(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
					FolderId:   &testFolderId,
				}).Return(1, nil)
			},
			want: noteTestExpect{
				error: nil,
			},
			wantErr: false,
		},
		{
			name: "move note from folder",
			args: noteTestArgs{
				userId:   1,
				noteId:   1,
				folderId: nil,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
				}, nil)
				repo.EXPECT().SaveEntity(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
					FolderId:   nil,
				}).Return(1, nil)
			},
			want: noteTestExpect{
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteService := noteService

			tt.mock()

			err := noteService.MoveToFolder(tt.args.userId, tt.args.noteId, tt.args.folderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteService.MoveToFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("NoteService.MoveToFolder() unexpected error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestConcreteNoteService_AddToFavorites(t *testing.T) {
	noteService, repo := initNoteServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    noteTestArgs
		want    noteTestExpect
		wantErr bool
	}{
		{
			name: "attempt to add to favorites unexisted note",
			args: noteTestArgs{
				userId: 1,
				noteId: 1,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "add note to favorites happy path",
			args: noteTestArgs{
				userId: 1,
				noteId: 1,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
					Timestamp:  time.Time{},
				}, nil)
				repo.EXPECT().SaveEntity(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: true,
					Timestamp:  time.Time{},
				}).Return(1, nil)
			},
			want: noteTestExpect{
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteService := noteService

			tt.mock()

			err := noteService.AddToFavorites(tt.args.userId, tt.args.noteId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteService.AddToFavorites() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("NoteService.AddToFavorites() unexpected error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestConcreteNoteService_DeleteFromFavorites(t *testing.T) {
	noteService, repo := initNoteServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    noteTestArgs
		want    noteTestExpect
		wantErr bool
	}{
		{
			name: "attempt to delete from favorites unexisted note",
			args: noteTestArgs{
				userId: 1,
				noteId: 1,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{}, model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil))
			},
			want: noteTestExpect{
				error: model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil),
			},
			wantErr: true,
		},
		{
			name: "delete note from favorites happy path",
			args: noteTestArgs{
				userId: 1,
				noteId: 1,
			},
			mock: func() {
				repo.EXPECT().GetNoteById(1, 1).Return(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: true,
					Timestamp:  time.Time{},
				}, nil)
				repo.EXPECT().SaveEntity(&model.Note{
					Id:         1,
					Title:      "title",
					Content:    "content",
					UserId:     1,
					IsFavorite: false,
					Timestamp:  time.Time{},
				}).Return(1, nil)
			},
			want: noteTestExpect{
				error: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteService := noteService

			tt.mock()

			err := noteService.DeleteFromFavorites(tt.args.userId, tt.args.noteId)
			if (err != nil) != tt.wantErr {
				t.Errorf("NoteService.DeleteFromFavorites() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && (err.Type != tt.want.error.Type || err.Message != tt.want.error.Message) {
				t.Errorf("NoteService.DeleteFromFavorites() unexpected error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestConcreteNoteService_FindNotesByQueryPhrase(t *testing.T) {
	noteService, repo := initNoteServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    noteTestArgs
		want    noteTestExpect
		wantErr bool
	}{
		{
			name: "user with no notes receives empty list",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{})
			},
			args: noteTestArgs{
				userId: 1,
				query:  "query",
			},
			want: noteTestExpect{
				notes: []*model.Note{},
			},
			wantErr: false,
		},
		{
			name: "empty query returns all notes",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:      1,
						Title:   "title1",
						Content: "content1",
						UserId:  1,
					},
					{
						Id:      2,
						Title:   "title2",
						Content: "content2",
						UserId:  1,
					},
				})
			},
			args: noteTestArgs{
				userId: 1,
				query:  "",
			},
			want: noteTestExpect{
				notes: []*model.Note{
					{
						Id:      1,
						Title:   "title1",
						Content: "content1",
						UserId:  1,
					},
					{
						Id:      2,
						Title:   "title2",
						Content: "content2",
						UserId:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong query returns empty list",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:      1,
						Title:   "title1",
						Content: "content1",
						UserId:  1,
					},
					{
						Id:      2,
						Title:   "title2",
						Content: "content2",
						UserId:  1,
					},
				})
			},
			args: noteTestArgs{
				userId: 1,
				query:  "query",
			},
			want: noteTestExpect{
				notes: []*model.Note{},
			},
			wantErr: false,
		},
		{
			name: "find notes by title",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:      1,
						Title:   "first title",
						Content: "content1",
						UserId:  1,
					},
					{
						Id:      2,
						Title:   "second title",
						Content: "content2",
						UserId:  1,
					},
				})
			},
			args: noteTestArgs{
				userId: 1,
				query:  "first",
			},
			want: noteTestExpect{
				notes: []*model.Note{
					{
						Id:      1,
						Title:   "first title",
						Content: "content1",
						UserId:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "find notes by content",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:      1,
						Title:   "title1",
						Content: "first content",
						UserId:  1,
					},
					{
						Id:      2,
						Title:   "title2",
						Content: "content2",
						UserId:  1,
					},
				})
			},
			args: noteTestArgs{
				userId: 1,
				query:  "first",
			},
			want: noteTestExpect{
				notes: []*model.Note{
					{
						Id:      1,
						Title:   "title1",
						Content: "first content",
						UserId:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "find notes by tags",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:      1,
						Title:   "title1",
						Content: "content1",
						UserId:  1,
						Tags: &[]string{
							"first", "second",
						},
					},
					{
						Id:      2,
						Title:   "title2",
						Content: "content2",
						UserId:  1,
					},
				})
			},
			args: noteTestArgs{
				userId: 1,
				query:  "first",
			},
			want: noteTestExpect{
				notes: []*model.Note{
					{
						Id:      1,
						Title:   "title1",
						Content: "content1",
						UserId:  1,
						Tags: &[]string{
							"first", "second",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteService := noteService

			tt.mock()

			got := noteService.FindNotesByQueryPhrase(tt.args.userId, tt.args.query)
			gotJson, _ := json.Marshal(got)
			expectedJson, _ := json.Marshal(tt.want.notes)
			if fmt.Sprintf("%v", string(gotJson)) != fmt.Sprintf("%v", string(expectedJson)) {
				t.Errorf("noteService.FindNotesByQueryPhrase() = %v, want %v", string(gotJson), string(expectedJson))
			}
		})
	}
}

func TestConcreteNoteService_GetFavoriteNotes(t *testing.T) {
	noteService, repo := initNoteServiceTest(t)

	tests := []struct {
		name    string
		mock    func()
		args    noteTestArgs
		want    noteTestExpect
		wantErr bool
	}{
		{
			name: "user with no notes receives empty list",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{})
			},
			args: noteTestArgs{
				userId: 1,
			},
			want: noteTestExpect{
				notes: []*model.Note{},
			},
			wantErr: false,
		},
		{
			name: "user with no favorite notes receives empty list",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:      1,
						Title:   "title1",
						Content: "content1",
						UserId:  1,
					},
					{
						Id:      2,
						Title:   "title2",
						Content: "content2",
						UserId:  1,
					},
				})
			},
			args: noteTestArgs{
				userId: 1,
			},
			want: noteTestExpect{
				notes: []*model.Note{},
			},
			wantErr: false,
		},
		{
			name: "user with favorite notes receives data",
			mock: func() {
				repo.EXPECT().GetNotesByUserId(1).Return([]*model.Note{
					{
						Id:         1,
						Title:      "title1",
						Content:    "content1",
						UserId:     1,
						IsFavorite: true,
					},
					{
						Id:      2,
						Title:   "title2",
						Content: "content2",
						UserId:  1,
					},
				})
			},
			args: noteTestArgs{
				userId: 1,
				query:  "query",
			},
			want: noteTestExpect{
				notes: []*model.Note{
					{
						Id:         1,
						Title:      "title1",
						Content:    "content1",
						UserId:     1,
						IsFavorite: true,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteService := noteService

			tt.mock()

			got := noteService.GetFavoriteNotes(tt.args.userId)
			gotJson, _ := json.Marshal(got)
			expectedJson, _ := json.Marshal(tt.want.notes)
			if fmt.Sprintf("%v", string(gotJson)) != fmt.Sprintf("%v", string(expectedJson)) {
				t.Errorf("noteService.GetFavoriteNotes() = %v, want %v", string(gotJson), string(expectedJson))
			}
		})
	}
}
