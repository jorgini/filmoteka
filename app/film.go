package app

import (
	"encoding/json"
	"errors"
)

type Film struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	IssueDate   Date   `json:"issue_date" db:"issue_date"`
	Rating      int    `json:"rating" db:"rating"`
}

func (f *Film) UnmarshalJSON(data []byte) error {
	result := struct {
		Id          *int    `json:"id"`
		Title       *string `json:"title"`
		Description *string `json:"description"`
		IssueDate   *Date   `json:"issue_date"`
		Rating      *int    `json:"rating"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	if result.Title == nil || result.IssueDate == nil || result.Rating == nil ||
		*result.Rating < 0 || *result.Rating > 10 {
		return errors.New("invalid state for required field(s)")
	} else {
		if result.Id != nil {
			f.Id = *result.Id
		}
		f.Title = *result.Title
		if result.Description != nil {
			f.Description = *result.Description
		}
		f.IssueDate = *result.IssueDate
		f.Rating = *result.Rating
	}
	return nil
}
