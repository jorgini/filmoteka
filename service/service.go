package service

import (
	"github.com/jorgini/filmoteka"
	"github.com/jorgini/filmoteka/models_dao"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type User interface {
	CreateUser(user filmoteka.User) (int, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, error)
	DeleteUserById(id int) error
	ValidateUser(id int) (bool, error)
	UpdateUser(login, userRole string) error
}

type Actor interface {
	CreateActor(actor filmoteka.Actor) (int, error)
	UpdateActor(actor filmoteka.UpdateActorInput) error
	GetActorById(id int) (filmoteka.ActorListItem, error)
	GetActorId(name, surname string) (int, error)
	GetActorsList(page, limit int) ([]filmoteka.ActorListItem, error)
	SearchActor(page, limit int, fragment filmoteka.ActorSearchFragment) ([]filmoteka.ActorListItem, error)
	DeleteActorById(id int) error
}

type Film interface {
	Actor
	CreateFilm(film filmoteka.InputFilm) (int, error)
	UpdateFilm(film filmoteka.UpdateFilmInput) error
	GetSortedFilmList(sortBy string, page, limit int) ([]filmoteka.InputFilm, error)
	GetCurFilm(id int) (filmoteka.InputFilm, error)
	GetSearchFilmList(page, limit int, fragment filmoteka.FilmSearchFragment) ([]filmoteka.InputFilm, error)
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
