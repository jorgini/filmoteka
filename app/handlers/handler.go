package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jorgini/filmoteka/app/service"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type Router struct {
	http.Handler
	service *service.Service
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	route := strings.Split(request.RequestURI, "?")[0]

	switch request.Method {
	case "POST":
		switch route {
		case "/actors":
			if r.userIdentity(writer, request) {
				r.createNewActor(writer, request)
			}
		case "/films":
			if r.userIdentity(writer, request) {
				r.createNewFilm(writer, request)
			}
		case "/users":
			r.createNewUser(writer, request)
		default:
			r.sendErrorResponse(writer, http.StatusNotFound, "this uri not found")
		}
	case "GET":
		switch route {
		case "/actors/list":
			r.getActorsList(writer, request)
		case "/actors":
			r.getActorById(writer, request)
		case "/actors/search":
			r.searchActor(writer, request)
		case "/films/list":
			r.getSortedFilmList(writer, request)
		case "/films":
			r.getCurrentFilm(writer, request)
		case "/films/search":
			r.getSearchFilmList(writer, request)
		case "/users":
			r.authUser(writer, request)
		default:
			r.sendErrorResponse(writer, http.StatusNotFound, fmt.Sprintf("this uri not found %s", route))
		}
	case "PUT":
		if r.userIdentity(writer, request) {
			switch route {
			case "/actors":
				r.updateActor(writer, request)
			case "/films":
				r.updateFilm(writer, request)
			case "/users":
				r.updateUser(writer, request)
			default:
				r.sendErrorResponse(writer, http.StatusNotFound, "this uri not found")
			}
		}
	case "DELETE":
		if r.userIdentity(writer, request) {
			switch route {
			case "/actors":
				r.deleteActor(writer, request)
			case "/films":
				r.deleteFilm(writer, request)
			case "/users":
				r.deleteUser(writer, request)
			default:
				r.sendErrorResponse(writer, http.StatusNotFound, "this uri not found")
			}
		}
	default:
		r.sendErrorResponse(writer, http.StatusMethodNotAllowed, "this method not allowed")
	}
}

func NewRouter(service *service.Service) *Router {
	return &Router{
		service: service,
	}
}

func (r *Router) sendErrorResponse(writer http.ResponseWriter, status int, message string) {
	logrus.Error(message)
	writer.WriteHeader(status)
	if _, err := writer.Write([]byte(message)); err != nil {
		writer.WriteHeader(http.StatusBadGateway)
	}
}

func parseBody(body io.ReadCloser, value interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(value); err != nil {
		return err
	}
	return nil
}

func writeBody[T any](writer http.ResponseWriter, values ...T) error {
	enc := json.NewEncoder(writer)
	for i := range values {
		if err := enc.Encode(values[i]); err != nil {
			return err
		}
	}
	return nil
}
