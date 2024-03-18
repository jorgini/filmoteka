package handlers

import (
	"fmt"
	"github.com/jorgini/filmoteka"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var (
	sortingOption = map[string]struct{}{"title": {}, "rating": {}, "issue_date": {}}
)

func createNewFilm(r *Router, writer http.ResponseWriter, request *http.Request) {
	id, err := getUserId(request)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	var film filmoteka.InputFilm
	if err := parseBody(request.Body, &film); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	filmId, err := r.service.Film.CreateFilm(film)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err = writeBody(writer, fmt.Sprintf("successfully create film with id %d", filmId)); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	logrus.Infof("new film with id %d was created by user with id %d", filmId, id)
}

func updateFilm(r *Router, writer http.ResponseWriter, request *http.Request) {
	id, err := getUserId(request)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	filmId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no id specified to update film")
		return
	}

	var update filmoteka.UpdateFilmInput
	if err = parseBody(request.Body, &update); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}
	update.Id = &filmId

	if err = r.service.Film.UpdateFilm(update); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err = writeBody(writer, "successfully update"); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	logrus.Infof("film with id %d was updated by user with id %d", filmId, id)
}

func getSortedFilmList(r *Router, writer http.ResponseWriter, request *http.Request) {
	sort := request.URL.Query().Get("sort_by")
	if sort == "" {
		sort = "rating"
	} else if _, ok := sortingOption[sort]; !ok {
		r.sendErrorResponse(writer, http.StatusBadRequest, "invalid parameter to sort films list")
		return
	}

	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no page specified for sorted list")
		return
	}
	if page < 0 {
		r.sendErrorResponse(writer, http.StatusBadRequest, "page out of bounds")
		return
	}

	films, err := r.service.Film.GetSortedFilmList(sort, page, limitOnPage)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if len(films) == 0 {
		r.sendErrorResponse(writer, http.StatusBadRequest, "page out of bounds")
		return
	}

	if err := writeBody(writer, films); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("sorted by %s list of films was sent to user", sort)
}

func getCurrentFilm(r *Router, writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id for get film not specified")
		return
	}
	if id < 1 {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id out of bounds")
		return
	}

	film, err := r.service.Film.GetCurFilm(id)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeBody(writer, film); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("film with id %d was sent to user", id)
}

func getSearchFilmList(r *Router, writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no page specified for search list")
		return
	}

	var input filmoteka.FilmSearchFragment
	if err := parseBody(request.Body, &input); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	films, err := r.service.Film.GetSearchFilmList(page, limitOnPage, input)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if len(films) == 0 {
		r.sendErrorResponse(writer, http.StatusBadRequest, "nothing find")
		return
	}

	if err := writeBody(writer, films); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Info("search list of films was sent to user")
}

func deleteFilm(r *Router, writer http.ResponseWriter, request *http.Request) {
	id, err := getUserId(request)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	filmId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id doesnt specified to delete film")
		return
	}

	if err = r.service.Film.DeleteFilmById(filmId); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err = writeBody(writer, "successfully delete"); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	logrus.Infof("film with id %d was deleted by user with id %d", filmId, id)
}
