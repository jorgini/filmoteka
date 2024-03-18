package filmoteka

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type InputActor struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func (a *InputActor) UnmarshalJSON(data []byte) error {
	result := struct {
		Name    *string `json:"name"`
		Surname *string `json:"surname"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Name == nil || result.Surname == nil {
		return errors.New("invalid state for required filed(s)")
	} else {
		a.Name = *result.Name
		a.Surname = *result.Surname
	}
	return nil
}

type Cast []InputActor

func (c *Cast) UnmarshalJSON(data []byte) error {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	for _, o := range obj.([]interface{}) {
		var a InputActor
		data, err := json.Marshal(o)
		if err != nil {
			return err
		}
		if err = a.UnmarshalJSON(data); err != nil {
			return err
		}

		*c = append(*c, a)
	}
	return nil
}

type InputFilm struct {
	Film
	Cast
}

func (i *InputFilm) UnmarshalJSON(data []byte) error {
	var f Film
	actors := struct {
		Cast Cast `json:"Cast"`
	}{}
	if err := f.UnmarshalJSON(data); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &actors); err != nil {
		return err
	}
	i.Film = f
	i.Cast = actors.Cast
	return nil
}

type ActorListItem struct {
	Actor Actor
	Films []Film
}

type ActorSearchFragment struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
}

func (a *ActorSearchFragment) UnmarshalJSON(data []byte) error {
	type tmp ActorSearchFragment
	var result tmp
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Name == nil && result.Surname == nil {
		return errors.New("parameters for search not specified")
	} else {
		a.Name = result.Name
		a.Surname = result.Surname
	}
	return nil
}

type FilmSearchFragment struct {
	Title   *string `json:"title"`
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
}

func (a *FilmSearchFragment) UnmarshalJSON(data []byte) error {
	type tmp FilmSearchFragment
	var result tmp
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Name == nil && result.Surname == nil && result.Title == nil {
		return errors.New("parameters for search not specified")
	} else {
		a.Title = result.Title
		a.Name = result.Name
		a.Surname = result.Surname
	}
	return nil
}

type UpdateFilmInput struct {
	Id          *int    `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IssueDate   *Date   `json:"issue_date"`
	Rating      *int    `json:"int"`
	Actors      *Cast   `json:"cast"`
}

func (u *UpdateFilmInput) UnmarshalJSON(data []byte) error {
	type Result UpdateFilmInput
	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Title == nil && result.Description == nil && result.IssueDate == nil &&
		result.Rating == nil && result.Actors == nil {
		return errors.New("invalid state for required filed to update film")
	} else {
		u.Title = result.Title
		u.Description = result.Description
		u.IssueDate = result.IssueDate
		u.Rating = result.Rating
		u.Actors = result.Actors
	}
	return nil
}

func (u *UpdateFilmInput) GetValuesUpdate() string {
	values := make([]string, 0, 4)
	if u.Title != nil {
		values = append(values, "title=$1")
	}
	if u.Description != nil {
		values = append(values, fmt.Sprintf("description=$%d", len(values)+1))
	}
	if u.IssueDate != nil {
		values = append(values, fmt.Sprintf("issue_date=$%d", len(values)+1))
	}
	if u.Rating != nil {
		values = append(values, fmt.Sprintf("rating=$%d", len(values)+1))
	}
	return strings.Join(values, ",\n")
}

func (u *UpdateFilmInput) GetArgsUpdate() []interface{} {
	args := make([]interface{}, 0, 4)
	if u.Title != nil {
		args = append(args, *u.Title)
	}
	if u.Description != nil {
		args = append(args, *u.Description)
	}
	if u.IssueDate != nil {
		args = append(args, u.IssueDate.String())
	}
	if u.Rating != nil {
		args = append(args, *u.Rating)
	}
	return args
}

type UpdateActorInput struct {
	Id       *int    `json:"id"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Sex      *string `json:"sex"`
	Birthday *Date   `json:"birthday"`
}

func (u *UpdateActorInput) UnmarshalJSON(data []byte) error {
	type Result UpdateActorInput
	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Sex == nil && result.Name == nil && result.Surname == nil && result.Birthday == nil {
		return errors.New("invalid state for required filed id to update actor")
	} else {
		u.Name = result.Name
		u.Surname = result.Surname
		u.Sex = result.Sex
		u.Birthday = result.Birthday
	}
	return nil
}

func (u *UpdateActorInput) GetValuesUpdate() string {
	values := make([]string, 0, 4)
	if u.Name != nil {
		values = append(values, fmt.Sprintf("name=$%d", len(values)+1))
	}
	if u.Surname != nil {
		values = append(values, fmt.Sprintf("surname=$%d", len(values)+1))
	}
	if u.Sex != nil {
		values = append(values, fmt.Sprintf("sex=$%d", len(values)+1))
	}
	if u.Birthday != nil {
		values = append(values, fmt.Sprintf("birthday=$%d", len(values)+1))
	}
	return strings.Join(values, ",\n")
}

func (u *UpdateActorInput) GetArgsUpdate() []interface{} {
	args := make([]interface{}, 0, 4)
	if u.Name != nil {
		args = append(args, *u.Name)
	}
	if u.Surname != nil {
		args = append(args, *u.Surname)
	}
	if u.Sex != nil {
		args = append(args, *u.Sex)
	}
	if u.Birthday != nil {
		args = append(args, u.Birthday.String())
	}
	return args
}
