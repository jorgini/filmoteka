package models_dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/jorgini/filmoteka"
)

type User interface {
	CreateUser(tx *sqlx.Tx, user filmoteka.User) (int, error)
	GetUser(login, password string) (filmoteka.User, error)
	DeleteUserById(tx *sqlx.Tx, id int) error
	ValidateUser(id int) (bool, error)
	UpdateUser(tx *sqlx.Tx, login, userRole string) error
}

type Actor interface {
	CreateActor(tx *sqlx.Tx, actor filmoteka.Actor) (int, error)
	UpdateActor(tx *sqlx.Tx, actor filmoteka.UpdateActorInput) error
	GetActorId(name, surname string) (int, error)
	GetActorById(id int) (filmoteka.Actor, error)
	GetActorsList(page, limit int) ([]filmoteka.Actor, error)
	SearchActor(page, limit int, name, surname *string) ([]filmoteka.Actor, error)
	GetFilmsWithCurActor(actorId int) ([]filmoteka.Film, error)
	DeleteActorById(tx *sqlx.Tx, id int) error
}

type Film interface {
	CreateFilm(tx *sqlx.Tx, film filmoteka.Film) (int, error)
	AddDependency(tx *sqlx.Tx, filmId, actorId int) error
	UpdateFilm(tx *sqlx.Tx, film filmoteka.UpdateFilmInput) error
	UpdateDependencies(tx *sqlx.Tx, filmId int, actorId ...int) error
	GetSortedFilmList(sortBy string, page, limit int) ([]filmoteka.Film, error)
	GetFilmListByTitle(page, limit int, title string) ([]filmoteka.Film, error)
	GetCurFilm(id int) (filmoteka.Film, error)
	GetActorsInCurFilm(filmId int) ([]filmoteka.InputActor, error)
	GetFilmListByActor(page, limit int, fragment filmoteka.FilmSearchFragment) ([]filmoteka.Film, error)
	DeleteFilmById(tx *sqlx.Tx, id int) error
}

type Transaction interface {
	StartTransaction() (*sqlx.Tx, error)
	ShutDown(tx *sqlx.Tx, err error) error
	Commit(tx *sqlx.Tx) error
}

type Repository struct {
	User
	Actor
	Film
	Transaction
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:        NewUserDao(db),
		Actor:       NewActorDao(db),
		Film:        NewFilmDao(db),
		Transaction: NewTransaction(db),
	}
}
