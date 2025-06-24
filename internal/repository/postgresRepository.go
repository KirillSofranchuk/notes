package repository

import (
	"Notes/internal/constants"
	"Notes/internal/model"
	"Notes/internal/utils"
	"errors"
	"gorm.io/gorm"
	"log"
)

var (
	EntityNotFoundError = model.NewApplicationError(model.ErrorTypeNotFound, "сущность не найдена", nil)
	DataBaseError       = model.NewApplicationError(model.ErrorTypeDatabase, " внутрення ошибка БД", nil)
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) AbstractRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) SaveEntity(entity model.BusinessEntity) (int, *model.ApplicationError) {
	entity.SetTimestamp()
	log.Println(entity)
	if entity.GetId() == 0 {
		return p.createEntity(entity)
	} else {
		return p.updateEntity(entity)
	}
}

func (p *PostgresRepository) createEntity(entity model.BusinessEntity) (int, *model.ApplicationError) {
	result := p.db.Create(entity)

	if result.Error != nil {
		return -1, DataBaseError
	}

	return entity.GetId(), nil
}

func (p *PostgresRepository) updateEntity(entity model.BusinessEntity) (int, *model.ApplicationError) {
	switch e := entity.(type) {
	case *model.Note:
		result := p.db.Save(e)
		if result.Error != nil {
			return -1, DataBaseError
		}
		return e.Id, nil

	case *model.User:
		result := p.db.Save(e)
		if result.Error != nil {
			return -1, DataBaseError
		}
		return e.Id, nil

	case *model.Folder:
		result := p.db.Save(e)
		if result.Error != nil {
			return -1, DataBaseError
		}
		return e.Id, nil

	default:
		return constants.FakeId, DataBaseError
	}
}

func (p *PostgresRepository) DeleteEntity(entity model.BusinessEntity) *model.ApplicationError {
	result := p.db.Delete(entity)

	if result.Error != nil {
		return DataBaseError
	}

	return nil
}

func (p *PostgresRepository) GetUserById(id int) (*model.User, *model.ApplicationError) {
	var user model.User
	result := p.db.First(&user, id) // где id - идентификатор пользователя

	// Проверка ошибок
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, EntityNotFoundError
		} else {
			return nil, DataBaseError
		}
	}
	return &user, nil
}

func (p *PostgresRepository) GetFolderById(id int, userId int) (*model.Folder, *model.ApplicationError) {
	var folder model.Folder
	result := p.db.Where("id = ? AND user_id = ?", id, userId).First(&folder)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, DataBaseError
	}
	return &folder, nil
}

func (p *PostgresRepository) GetNoteById(id int, userId int) (*model.Note, *model.ApplicationError) {
	var note model.Note
	result := p.db.Table("notes").Where("id = ? AND user_id = ?", id, userId).First(&note)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, DataBaseError
	}
	return &note, nil
}

func (p *PostgresRepository) GetUser(login, password string) (*model.User, *model.ApplicationError) {
	var user model.User
	result := p.db.Where("login = ?", login).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, EntityNotFoundError
		}
		return nil, DataBaseError
	}

	arePasswordsEqual, err := utils.CompareHashAndPassword(user.Password, password)

	if err != nil {
		return nil, err
	}

	if arePasswordsEqual {
		return &user, nil
	}

	return &user, nil
}

func (p *PostgresRepository) GetFoldersByUserId(userId int) []*model.Folder {
	var folders []*model.Folder
	result := p.db.Where("user_id = ?", userId).Find(&folders)
	if result.Error != nil {
		return make([]*model.Folder, 0)
	}
	return folders
}

func (p *PostgresRepository) GetNotesByUserId(userId int) []*model.Note {
	var notes []*model.Note
	result := p.db.Where("user_id = ?", userId).Find(&notes)

	if result.Error != nil {
		return make([]*model.Note, 0)
	}
	return notes
}

func (p *PostgresRepository) GetUsers() []*model.User {
	var users []*model.User
	result := p.db.Find(&users)

	if result.Error != nil {
		return make([]*model.User, 0)
	}

	return users
}
