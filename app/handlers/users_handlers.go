package handlers

import (
	"encoding/json"
	"errors"
	"github.com/jorgini/filmoteka/app"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (r *Router) createNewUser(writer http.ResponseWriter, request *http.Request) {
	var user app.User
	if err := parseBody(request.Body, &user); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	id, err := r.service.User.CreateUser(user)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("a new user with an id %d has been registered", id)
}

type signInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *signInput) UnmarshalJSON(data []byte) error {
	result := struct {
		Login    *string `json:"login"`
		Password *string `json:"password"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Login == nil || result.Password == nil {
		return errors.New("missing required fields")
	} else {
		s.Login = *result.Login
		s.Password = *result.Password
	}
	return nil
}

func (r *Router) authUser(writer http.ResponseWriter, request *http.Request) {
	var input signInput
	if err := parseBody(request.Body, &input); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	token, err := r.service.User.GenerateToken(input.Login, input.Password)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return
	}

	if err := writeBody(writer, token); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("the %s user logged in", input.Login)
}

type updateInput struct {
	Login    string `json:"login"`
	UserRole string `json:"user_role"`
}

func (u *updateInput) UnmarshalJSON(data []byte) error {
	result := struct {
		Login    *string `json:"login"`
		UserRole *string `json:"user_role"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.Login == nil || result.UserRole == nil || (*result.UserRole != "regular" &&
		*result.UserRole != "admin") {
		return errors.New("invalid state for required filed(s)")
	} else {
		u.Login = *result.Login
		u.UserRole = *result.UserRole
	}
	return nil
}

func (r *Router) updateUser(writer http.ResponseWriter, request *http.Request) {
	ok, id := r.validateUser(writer, request)

	if !ok {
		return
	}

	var update updateInput
	if err := parseBody(request.Body, &update); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	err := r.service.User.UpdateUser(update.Login, update.UserRole)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("the role of the %s user has been changed to %s by a user with an id %d",
		update.Login, update.UserRole, id)
}

func (r *Router) deleteUser(writer http.ResponseWriter, request *http.Request) {
	id, err := getUserId(request)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	if err := r.service.User.DeleteUserById(id); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("user with id %d has been deleted", id)
}
