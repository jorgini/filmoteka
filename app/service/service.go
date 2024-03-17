package service

import (
	"github.com/jorgini/filmoteka/app"
	"github.com/jorgini/filmoteka/app/models_dao"
)

type User interface {
	CreateUser(user app.User) (int, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, error)
	DeleteUserById(id int) error
	ValidateUser(id int) (bool, error)
	UpdateUser(login, userRole string) error
}

type Actor interface {
	CreateActor(actor app.Actor) (int, error)
	UpdateActor(actor app.UpdateActorInput) error
	GetActorById(id int) (app.ActorListItem, error)
	GetActorId(name, surname string) (int, error)
	GetActorsList(page, limit int) ([]app.ActorListItem, error)
	SearchActor(page, limit int, fragment app.ActorSearchFragment) ([]app.ActorListItem, error)
	DeleteActorById(id int) error
}

type Film interface {
	Actor
	CreateFilm(film app.InputFilm) (int, error)
	UpdateFilm(film app.UpdateFilmInput) error
	GetSortedFilmList(sortBy string, page, limit int) ([]app.InputFilm, error)
	GetCurFilm(id int) (app.InputFilm, error)
	GetSearchFilmList(page, limit int, fragment app.FilmSearchFragment) ([]app.InputFilm, error)
	DeleteFilmById(id int) error
}

type Service struct {
	User
	Actor
	Film
}

func NewService(dao *models_dao.Repository) *Service {
	return &Service{
		User:  NewUserService(dao.User, dao.Transaction),
		Actor: NewActorService(dao.Actor, dao.Transaction),
		Film:  NewFilmService(dao.Film, dao.Actor, dao.Transaction),
	}
}
