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

func validation(r *Router, writer http.ResponseWriter, request *http.Request) {
	header := request.Header.Get(authorizationHeader)
	if header == "" {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "empty header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := r.service.User.ParseToken(headerParts[1])
	if err != nil {
		r.sendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return
	}

	if ok, err := r.service.User.ValidateUser(userId); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	} else if !ok {
		r.sendErrorResponse(writer, http.StatusLocked, "this function locked for current user")
		return
	}

	*request = *request.Clone(context.WithValue(request.Context(), userCtx, userId))
	return
}

func userIdentity(r *Router, writer http.ResponseWriter, request *http.Request) {
	header := request.Header.Get(authorizationHeader)
	if header == "" {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "empty header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		r.sendErrorResponse(writer, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := r.service.User.ParseToken(headerParts[1])
	if err != nil {
		r.sendErrorResponse(writer, http.StatusUnauthorized, err.Error())
		return
	}

	*request = *request.Clone(context.WithValue(request.Context(), userCtx, userId))
}

func getUserId(r *http.Request) (int, error) {
	id := r.Context().Value(userCtx)
	if id == nil {
		return 0, errors.New("user id not found")
	}
	return id.(int), nil
}
