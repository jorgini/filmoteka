package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (r *Router) userIdentity(writer http.ResponseWriter, request *http.Request) bool {
	header := request.Header.Get(authorizationHeader)
	if header == "" {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "empty header")
		return false
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "invalid auth header")
		return false
	}

	if len(headerParts[1]) == 0 {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "token is empty")
		return false
	}

	userId, err := r.service.User.ParseToken(headerParts[1])
	if err != nil {
		r.sendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return false
	}

	*request = *request.Clone(context.WithValue(request.Context(), userCtx, userId))
	return true
}

func (r *Router) validateUser(writer http.ResponseWriter, request *http.Request) (bool, int) {
	id, err := getUserId(request)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return false, 0
	}

	if ok, err := r.service.User.ValidateUser(id); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return false, 0
	} else if !ok {
		r.sendErrorResponse(writer, http.StatusLocked, "this function locked for current user")
		return false, 0
	}
	return true, id
}

func getUserId(r *http.Request) (int, error) {
	id := r.Context().Value(userCtx)
	if id == nil {
		return 0, errors.New("user id not found")
	}
	return id.(int), nil
}
