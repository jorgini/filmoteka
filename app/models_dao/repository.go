package models_dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/jorgini/filmoteka/app"
)

type User interface {
	CreateUser(tx *sqlx.Tx, user app.User) (int, error)
	GetUser(login, password string) (app.User, error)
	DeleteUserById(tx *sqlx.Tx, id int) error
	ValidateUser(id int) (bool, error)
	UpdateUser(tx *sqlx.Tx, login, userRole string) error
}

type Actor interface {
	CreateActor(tx *sqlx.Tx, actor app.Actor) (int, error)
	UpdateActor(tx *sqlx.Tx, actor app.UpdateActorInput) error
	GetActorId(name, surname string) (int, error)
	GetActorById(id int) (app.Actor, error)
	GetActorsList(page, limit int) ([]app.Actor, error)
	SearchActor(page, limit int, name, surname *string) ([]app.Actor, error)
	GetFilmsWithCurActor(actorId int) ([]app.Film, error)
	DeleteActorById(tx *sqlx.Tx, id int) error
}

type Film interface {
	CreateFilm(tx *sqlx.Tx, film app.Film) (int, error)
	AddDependency(tx *sqlx.Tx, filmId, actorId int) error
	UpdateFilm(tx *sqlx.Tx, film app.UpdateFilmInput) error
	UpdateDependencies(tx *sqlx.Tx, filmId int, actorId ...int) error
	GetSortedFilmList(sortBy string, page, limit int) ([]app.Film, error)
	GetFilmListByTitle(page, limit int, title string) ([]app.Film, error)
	GetCurFilm(id int) (app.Film, error)
	GetActorsInCurFilm(filmId int) ([]app.InputActor, error)
	GetFilmListByActor(page, limit int, fragment app.FilmSearchFragment) ([]app.Film, error)
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
