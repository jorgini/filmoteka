package handlers

import (
	"github.com/jorgini/filmoteka/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var (
	sortingOption = map[string]struct{}{"title": {}, "rating": {}, "issue_date": {}}
)

func (r *Router) createNewFilm(writer http.ResponseWriter, request *http.Request) {
	ok, id := r.validateUser(writer, request)
	if !ok {
		return
	}

	var film app.InputFilm
	if err := parseBody(request.Body, &film); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	filmId, err := r.service.Film.CreateFilm(film)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("new film with id %d was created by user with id %d", filmId, id)
}

func (r *Router) updateFilm(writer http.ResponseWriter, request *http.Request) {
	ok, id := r.validateUser(writer, request)
	if !ok {
		r.sendErrorResponse(writer, http.StatusLocked, "this function locked for current user")
		return
	}

	filmId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no specified id for update film")
		return
	}

	var update app.UpdateFilmInput
	if err := parseBody(request.Body, &update); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}
	update.Id = &filmId

	if err := r.service.Film.UpdateFilm(update); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("film with id %d was updated by user with id %d", filmId, id)
}

func (r *Router) getSortedFilmList(writer http.ResponseWriter, request *http.Request) {
	sort := request.URL.Query().Get("sort_by")
	if sort == "" {
		sort = "rating"
	} else if _, ok := sortingOption[sort]; !ok {
		r.sendErrorResponse(writer, http.StatusBadRequest, "invalid parameter for sort films list")
	}

	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no page specified for sorted list")
	}

	films, err := r.service.Film.GetSortedFilmList(sort, page, limitOnPage)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeBody(writer, films); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("sorted by %s list of films was sent to user", sort)
}

func (r *Router) getCurrentFilm(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "id for get film not specified")
		return
	}

	film, err := r.service.Film.GetCurFilm(id)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	// Film here is transmitted by pointer for correct marshaling date
	if err := writeBody(writer, &film); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("film with id %d was sent to user", id)
}

func (r *Router) getSearchFilmList(writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no page specified for search list")
	}

	var input app.FilmSearchFragment
	if err := parseBody(request.Body, &input); err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}

	films, err := r.service.Film.GetSearchFilmList(page, limitOnPage, input)
	if err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeBody(writer, films); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Info("search list of films was sent to user")
}

func (r *Router) deleteFilm(writer http.ResponseWriter, request *http.Request) {
	ok, id := r.validateUser(writer, request)
	if !ok {
		return
	}

	filmId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		r.sendErrorResponse(writer, http.StatusBadRequest, "no id specified for delete film")
		return
	}

	if err := r.service.Film.DeleteFilmById(filmId); err != nil {
		r.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("film with id %d was deleted by user with id %d", filmId, id)
}
