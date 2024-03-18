package handlers

import (
	"fmt"
	"github.com/jorgini/filmoteka"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	limitOnPage = 10
)

func createNewActor(r *Router, writer http.ResponseWriter, request *http.Request) {
	id, err := getUserId(request)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	var actor filmoteka.Actor
	if err := parseBody(request.Body, &actor); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	actorId, err := r.service.Actor.CreateActor(actor)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err = writeBody(writer, fmt.Sprintf("successfully create actor with id %d", actorId)); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	logrus.Infof("new actor with id %d was created by user with id %d", actorId, id)
}

func updateActor(r *Router, writer http.ResponseWriter, request *http.Request) {
	id, err := getUserId(request)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	actorId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no id specified for update actor")
		return
	}

	var input filmoteka.UpdateActorInput
	if err = parseBody(request.Body, &input); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}
	input.Id = &actorId

	if err = r.service.Actor.UpdateActor(input); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err = writeBody(writer, "successfully update"); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	logrus.Infof("actor with id %d has been updated by user with id %d", *input.Id, id)
}

func getActorsList(r *Router, writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "page not selected")
		return
	}
	if page < 1 {
		r.sendErrorResponse(writer, http.StatusBadRequest, "page out of bounds")
		return
	}

	actors, err := r.service.Actor.GetActorsList(page, limitOnPage)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if len(actors) == 0 {
		r.sendErrorResponse(writer, http.StatusBadRequest, "page out of bounds")
		return
	}

	if err := writeBody(writer, actors...); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("list of actors in page %d was sending to user", page)
}

func getActorById(r *Router, writer http.ResponseWriter, request *http.Request) {
	actorId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id not selected")
		return
	}
	if actorId < 1 {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id out of bounds")
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
	logrus.Infof("actor %s was sent to user", actor.Actor.Surname)
}

func searchActor(r *Router, writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "page not selected")
		return
	}

	var fragment filmoteka.ActorSearchFragment
	if err := parseBody(request.Body, &fragment); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	actors, err := r.service.Actor.SearchActor(page, limitOnPage, fragment)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if len(actors) == 0 {
		r.sendErrorResponse(writer, http.StatusBadRequest, "nothing find")
	}

	if err := writeBody(writer, actors...); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Info("result of search actor with fragment name was sending to user")
}

func deleteActor(r *Router, writer http.ResponseWriter, request *http.Request) {
	id, err := getUserId(request)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	inputId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id doesnt specified to delete actor")
		return
	}

	if err = r.service.Actor.DeleteActorById(inputId); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if err = writeBody(writer, "successful delete"); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	logrus.Infof("actor with id %d has been deleted by user with id %d", inputId, id)
}
