package service

import (
	"github.com/jorgini/filmoteka/app"
	"github.com/jorgini/filmoteka/app/models_dao"
	"github.com/sirupsen/logrus"
)

type FilmService struct {
	Actor
	tx   models_dao.Transaction
	film models_dao.Film
}

func NewFilmService(filmDao models_dao.Film, actorDao models_dao.Actor, tx models_dao.Transaction) *FilmService {
	return &FilmService{
		Actor: NewActorService(actorDao, tx),
		tx:    tx,
		film:  filmDao,
	}
}

func (f *FilmService) CreateFilm(film app.InputFilm) (int, error) {
	transaction, err := f.tx.StartTransaction()
	if err != nil {
		return 0, err
	}

	filmId, err := f.film.CreateFilm(transaction, film.Film)
	if err != nil {
		return 0, f.tx.ShutDown(transaction, err)
	}

	for i, actor := range film.Cast {
		actorId, err := f.Actor.GetActorId(actor.Name, actor.Surname)
		if err != nil {
			return 0, f.tx.ShutDown(transaction, err)
		}

		if err := f.film.AddDependency(transaction, filmId, actorId); err != nil {
			logrus.Info(i)
			return 0, f.tx.ShutDown(transaction, err)
		}
	}

	return filmId, f.tx.Commit(transaction)
}

func (f *FilmService) UpdateFilm(film app.UpdateFilmInput) error {
	transaction, err := f.tx.StartTransaction()
	if err != nil {
		return err
	}

	if film.GetValuesUpdate() != "" {
		if err := f.film.UpdateFilm(transaction, film); err != nil {
			return f.tx.ShutDown(transaction, err)
		}
	}

	actorsIds := make([]int, len(*film.Actors))
	for i, actor := range *film.Actors {
		actorId, err := f.Actor.GetActorId(actor.Name, actor.Surname)
		if err != nil {
			return f.tx.ShutDown(transaction, err)
		}
		actorsIds[i] = actorId
	}

	if err := f.film.UpdateDependencies(transaction, *film.Id, actorsIds...); err != nil {
		return f.tx.ShutDown(transaction, err)
	}

	return f.tx.Commit(transaction)
}

func (f *FilmService) GetSortedFilmList(sortBy string, page, limit int) ([]app.InputFilm, error) {
	films, err := f.film.GetSortedFilmList(sortBy, page, limit)
	if err != nil {
		return nil, err
	}

	output := make([]app.InputFilm, len(films))
	for i := range output {
		output[i].Film = films[i]
		output[i].Cast, err = f.film.GetActorsInCurFilm(films[i].Id)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func (f *FilmService) GetCurFilm(id int) (app.InputFilm, error) {
	var film app.InputFilm
	var err error
	film.Film, err = f.film.GetCurFilm(id)
	if err != nil {
		return app.InputFilm{}, err
	}

	film.Cast, err = f.film.GetActorsInCurFilm(id)
	if err != nil {
		return app.InputFilm{}, err
	}

	return film, nil
}

func (f *FilmService) GetSearchFilmList(page, limit int, fragment app.FilmSearchFragment) ([]app.InputFilm, error) {
	if fragment.Name == nil && fragment.Surname == nil {
		films, err := f.film.GetFilmListByTitle(page, limit, *fragment.Title)
		if err != nil {
			return nil, err
		}

		output := make([]app.InputFilm, len(films))
		for i := range films {
			output[i].Film = films[i]
			output[i].Cast, err = f.film.GetActorsInCurFilm(films[i].Id)
			if err != nil {
				return nil, err
			}
		}
		return output, nil
	} else {
		films, err := f.film.GetFilmListByActor(page, limit, fragment)
		if err != nil {
			return nil, err
		}

		output := make([]app.InputFilm, len(films))
		for i := range films {
			output[i].Film = films[i]
			output[i].Cast, err = f.film.GetActorsInCurFilm(films[i].Id)
			if err != nil {
				return nil, err
			}
		}
		return output, nil
	}
}

func (f *FilmService) DeleteFilmById(id int) error {
	transaction, err := f.tx.StartTransaction()
	if err != nil {
		return err
	}

	if err = f.film.DeleteFilmById(transaction, id); err != nil {
		return f.tx.ShutDown(transaction, err)
	}

	return f.tx.Commit(transaction)
}
