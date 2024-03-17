package app

import (
	"encoding/json"
	"errors"
)

type User struct {
	Id       int    `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
	UserRole string `json:"user_role" db:"user_role"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	result := struct {
		Id       *int    `json:"id"`
		UserRole *string `json:"user_role"`
		Login    *string `json:"login"`
		Password *string `json:"password"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	if result.UserRole == nil || result.Login == nil || result.Password == nil ||
		(*result.UserRole != "regular" && *result.UserRole != "admin") {
		return errors.New("invalid state for required field(s)")
	} else {
		if result.Id != nil {
			u.Id = *result.Id
		}
		u.UserRole = *result.UserRole
		u.Login = *result.Login
		u.Password = *result.Password
	}
	return nil
}
