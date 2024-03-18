package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jorgini/filmoteka/service"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type Router struct {
	http.Handler
	service   *service.Service
	endpoints map[string]map[string]func(r *Router, writer http.ResponseWriter, request *http.Request)
	abort     bool
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	route := strings.Split(request.RequestURI, "?")[0]

	if _, ok := r.endpoints[request.Method]; !ok {
		r.sendErrorResponse(writer, http.StatusMethodNotAllowed, "this method not provided")
		r.abort = false
		return
	}

	if handle, ok := r.endpoints[request.Method][route]; !ok {
		r.sendErrorResponse(writer, http.StatusNotFound, fmt.Sprintf("this uri not found %s", route))
	} else {
		if request.Method == "GET" || (request.Method == "POST" && route == "/users") {
		} else if request.Method == "DELETE" && route == "/users" {
			userIdentity(r, writer, request)
		} else {
			validation(r, writer, request)
		}
		if !r.abort {
			handle(r, writer, request)
		}
	}
	r.abort = false
}

func NewRouter(service *service.Service) *Router {
	return &Router{
		service: service,
		endpoints: map[string]map[string]func(r *Router, writer http.ResponseWriter, request *http.Request){
			"POST": {"/users": createNewUser, "/actors": createNewActor, "/films": createNewFilm},
			"GET": {"/actors/list": getActorsList, "/actors": getActorById, "/actors/search": searchActor,
				"/films/list": getSortedFilmList, "/films": getCurrentFilm, "/films/search": getSearchFilmList,
				"/users": authUser},
			"PUT":    {"/actors": updateActor, "/films": updateFilm, "/users": updateUser},
			"DELETE": {"/actors": deleteActor, "/films": deleteFilm, "/users": deleteUser},
		},
	}
}

func (r *Router) AddEndPoint(method, path string, handle func(r *Router, writer http.ResponseWriter, request *http.Request)) {
	if r.endpoints == nil {
		r.endpoints = make(map[string]map[string]func(r *Router, writer http.ResponseWriter, request *http.Request))
	}

	if r.endpoints[method] == nil {
		r.endpoints[method] = make(map[string]func(r *Router, writer http.ResponseWriter, request *http.Request))
	}

	r.endpoints[method][path] = handle
}

func (r *Router) sendErrorResponse(writer http.ResponseWriter, status int, message string) {
	logrus.Error(message)
	r.abort = true
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
