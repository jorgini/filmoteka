package filmoteka

import (
	"encoding/json"
	"errors"
)

type Actor struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Surname  string `json:"surname" db:"surname"`
	Sex      string `json:"sex" db:"sex"`
	Birthday *Date  `json:"birthday" db:"birthday"`
}

func (a *Actor) UnmarshalJSON(data []byte) error {
	result := struct {
		Id       *int    `json:"id"`
		Name     *string `json:"name"`
		Surname  *string `json:"surname"`
		Sex      *string `json:"sex"`
		Birthday *Date   `json:"birthday"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	if result.Name == nil || result.Surname == nil || result.Sex == nil ||
		(*result.Sex != "male" && *result.Sex != "female") || result.Birthday == nil {
		return errors.New("invalid state for required field(s)")
	} else {
		if result.Id != nil {
			a.Id = *result.Id
		}
		a.Name = *result.Name
		a.Surname = *result.Surname
		a.Sex = *result.Sex
		a.Birthday = result.Birthday
	}
	return nil
}
