package repository

import "Notes/internal/model"

type AbstractRepository interface {
	SaveEntity(entity model.BusinessEntity)
}
