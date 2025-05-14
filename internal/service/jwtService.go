package service

import (
	"Notes/config"
	"Notes/internal/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtService struct {
	cfg *config.Config
}

func NewJwtService(cfg *config.Config) JwtService {
	return JwtService{
		cfg: cfg,
	}
}

func (j JwtService) GetToken(id int) (string, *model.ApplicationError) {
	expirationTime := time.Now().Add(time.Duration(j.cfg.App.TokenTtlHours) * time.Hour)

	claims := model.Claims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "note-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.cfg.App.Secret))

	if err != nil {
		return "", model.NewApplicationError(model.ErrorTypeInternal, "Ошибка при формировании токена", err)
	}
	return signedToken, nil
}

func (j JwtService) ParseToken(tokenString string) (*model.Claims, *model.ApplicationError) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.NewApplicationError(model.ErrorTypeAuth, "Неверный метод подписи", nil)
		}
		return []byte(j.cfg.App.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, model.NewApplicationError(model.ErrorTypeAuth, "Невалидный токен", nil)
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, model.NewApplicationError(model.ErrorTypeAuth, "Невалидный токен", nil)
}
