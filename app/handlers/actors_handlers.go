package handlers

import (
	"github.com/jorgini/filmoteka/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	limitOnPage = 10
)

func (r *Router) createNewActor(writer http.ResponseWriter, request *http.Request) {
	ok, id := r.validateUser(writer, request)
	if !ok {
		return
	}

	var actor app.Actor
	if err := parseBody(request.Body, &actor); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	actorId, err := r.service.Actor.CreateActor(actor)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("new actor with id %d was created by user with id %d", actorId, id)
}

func (r *Router) updateActor(writer http.ResponseWriter, request *http.Request) {
	ok, id := r.validateUser(writer, request)
	if !ok {
		return
	}

	actorId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no id specified for update actor")
		return
	}

	var input app.UpdateActorInput
	if err := parseBody(request.Body, &input); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}
	input.Id = &actorId

	if err := r.service.Actor.UpdateActor(input); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("actor with id %d has been updated by user with id %d", *input.Id, id)
}

func (r *Router) getActorsList(writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "page not selected")
		return
	}

	actors, err := r.service.Actor.GetActorsList(page, limitOnPage)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeBody(writer, actors...); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("list of actors in page %d was sending to user", page)
}

func (r *Router) getActorById(writer http.ResponseWriter, request *http.Request) {
	actorId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id not selected")
		return
	}

	actor, err := r.service.Actor.GetActorById(actorId)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeBody(writer, actor); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("actor %s was sent to user", actor)
}

func (r *Router) searchActor(writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "page not selected")
		return
	}

	var fragment app.ActorSearchFragment
	if err := parseBody(request.Body, &fragment); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	actors, err := r.service.Actor.SearchActor(page, limitOnPage, fragment)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeBody(writer, actors...); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Info("result of search actor with fragment name was sending to user")
}

func (r *Router) deleteActor(writer http.ResponseWriter, request *http.Request) {
	ok, id := r.validateUser(writer, request)
	if !ok {
		return
	}

	inputId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id doesnt specified to delete actor")
		return
	}

	if err := r.service.Actor.DeleteActorById(inputId); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("actor with id %d has been deleted by user with id %d", inputId, id)
}
