package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jorgini/filmoteka/app"
	"github.com/jorgini/filmoteka/app/models_dao"
	"time"
)

const (
	salt     = "foiewnjgin2jr34fnvwi0"
	signKey  = "cwenin314cI#knewvm349,wnd2"
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type UserService struct {
	dao models_dao.User
	tx  models_dao.Transaction
}

func NewUserService(dao models_dao.User, tx models_dao.Transaction) *UserService {
	return &UserService{
		dao: dao,
		tx:  tx,
	}
}

func (u *UserService) CreateUser(user app.User) (int, error) {
	user.Password = generateHashPassword(user.Password)

	transaction, err := u.tx.StartTransaction()
	if err != nil {
		return 0, err
	}

	var id int
	id, err = u.dao.CreateUser(transaction, user)
	if err != nil {
		return 0, u.tx.ShutDown(transaction, err)
	}

	return id, u.tx.Commit(transaction)
}

func (u *UserService) GenerateToken(login, password string) (string, error) {
	user, err := u.dao.GetUser(login, generateHashPassword(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(signKey))
}

func (u *UserService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (u *UserService) ValidateUser(id int) (bool, error) {
	return u.dao.ValidateUser(id)
}

func (u *UserService) UpdateUser(login, userRole string) error {
	transaction, err := u.tx.StartTransaction()
	if err != nil {
		return err
	}

	if err = u.dao.UpdateUser(transaction, login, userRole); err != nil {
		return u.tx.ShutDown(transaction, err)
	}
	return u.tx.Commit(transaction)
}

func (u *UserService) DeleteUserById(id int) error {
	transaction, err := u.tx.StartTransaction()
	if err != nil {
		return err
	}

	if err = u.dao.DeleteUserById(transaction, id); err != nil {
		return u.tx.ShutDown(transaction, err)
	}
	return u.tx.Commit(transaction)
}

func generateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
