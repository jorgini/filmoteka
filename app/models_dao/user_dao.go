package models_dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jorgini/filmoteka/app"
	"github.com/jorgini/filmoteka/app/configs"
)

type UserDao struct {
	db *sqlx.DB
}

func NewUserDao(db *sqlx.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (u *UserDao) CreateUser(tx *sqlx.Tx, user app.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (login, password, user_role) values ($1, $2, $3) RETURNING id",
		configs.EnvUserTable())

	row := tx.QueryRow(query, user.Login, user.Password, user.UserRole)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserDao) GetUser(login, password string) (app.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE login=$1 AND password=$2", configs.EnvUserTable())

	var user app.User
	if err := u.db.Get(&user, query, login, password); err != nil {
		return app.User{}, err
	}
	return user, nil
}

func (u *UserDao) DeleteUserById(tx *sqlx.Tx, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", configs.EnvUserTable())

	if _, err := tx.Query(query, id); err != nil {
		return err
	}
	return nil
}

func (u *UserDao) ValidateUser(id int) (bool, error) {
	query := fmt.Sprintf("SELECT user_role FROM %s WHERE id=$1", configs.EnvUserTable())

	var userRole string
	row := u.db.QueryRow(query, id)
	if err := row.Scan(&userRole); err != nil {
		return false, err
	}
	if userRole == "admin" {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *UserDao) UpdateUser(tx *sqlx.Tx, login, userRole string) error {
	query := fmt.Sprintf("UPDATE %s SET user_role=$1 WHERE login=$2", configs.EnvUserTable())

	_, err := tx.Query(query, userRole, login)
	if err != nil {
		return err
	}
	return nil
}
