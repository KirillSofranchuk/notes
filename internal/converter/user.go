package converter

import (
	"project/internal/model/api"
	"project/internal/model/db"
)

func ApiUserToDb(api *api.UserRequest) *db.User {
	return &db.User{
		Name:  api.Name,
		Login: api.Login,
		// Пока не хэшируем пароль, так как просто накидываем модели и маппинг.
		Password: api.Password,
	}
}
