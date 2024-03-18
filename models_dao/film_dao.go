package models_dao

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jorgini/filmoteka"
	"github.com/jorgini/filmoteka/configs"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"strings"
)

type FilmDao struct {
	db *sqlx.DB
}

func NewFilmDao(db *sqlx.DB) *FilmDao {
	return &FilmDao{
		db: db,
	}
}

const (
	searchQuery = `
		SELECT DISTINCT f.* 
		FROM %s s
		INNER JOIN %s f ON (s.film_id=f.id) 
		INNER JOIN %s a ON (s.actor_id=a.id)
		WHERE %s
		LIMIT $1
		OFFSET $2
		`
)

func getFilmSearchCondition(fragment filmoteka.FilmSearchFragment) string {
	condition := make([]string, 0, 3)
	if fragment.Title != nil {
		condition = append(condition, fmt.Sprintf("POSITION('%s' in f.title)>0", *fragment.Title))
	}
	if fragment.Name != nil {
		condition = append(condition, fmt.Sprintf("POSITION('%s' in a.name)>0", *fragment.Name))
	}
	if fragment.Surname != nil {
		condition = append(condition, fmt.Sprintf("POSITION('%s' in a.surname)>0", *fragment.Surname))
	}
	return strings.Join(condition, " AND ")
}

func (f *FilmDao) CreateFilm(tx *sqlx.Tx, film filmoteka.Film) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (title,description,issue_date,rating) values ($1, $2, TO_DATE($3,'DD-MM-YYYY'), $4) RETURNING id",
		configs.EnvFilmTable())

	var id int
	row := tx.QueryRow(query, film.Title, film.Description, film.IssueDate.String(), film.Rating)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (f *FilmDao) AddDependency(tx *sqlx.Tx, filmId, actorId int) error {
	query := fmt.Sprintf("INSERT INTO %s (actor_id,film_id) values ($1,$2)", configs.EnvStarredTable())

	if _, err := tx.Exec(query, actorId, filmId); err != nil {
		logrus.Infof(query)
		return err
	}
	return nil
}

func (f *FilmDao) UpdateFilm(tx *sqlx.Tx, film filmoteka.UpdateFilmInput) error {
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", configs.EnvFilmTable(), film.GetValuesUpdate(), *film.Id)

	if _, err := tx.Exec(query, film.GetArgsUpdate()...); err != nil {
		logrus.Info(query)
		return err
	}
	return nil
}

func (f *FilmDao) UpdateDependencies(tx *sqlx.Tx, filmId int, actorIds ...int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE film_id=$1 AND array_position($2, actor_id)=NULL",
		configs.EnvStarredTable())

	if _, err := tx.Exec(query, filmId, pq.Array(actorIds)); err != nil {
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (film_id, actor_id) values ($1, $2) ON CONFLICT DO NOTHING",
		configs.EnvStarredTable())

	for _, aid := range actorIds {
		if _, err := tx.Exec(query, filmId, aid); err != nil {
			logrus.Info(query)
			return err
		}
	}
	return nil
}

func (f *FilmDao) GetSortedFilmList(sortBy string, page, limit int) ([]filmoteka.Film, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s DESC LIMIT $1 OFFSET $2", configs.EnvFilmTable(), sortBy)

	var films []filmoteka.Film
	if err := f.db.Select(&films, query, limit, limit*(page-1)); err != nil {
		return nil, err
	}
	return films, nil
}

func (f *FilmDao) GetFilmListByTitle(page, limit int, title string) ([]filmoteka.Film, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE POSITION('%s' in title) > 0 LIMIT $1 OFFSET $2",
		configs.EnvFilmTable(), title)

	var films []filmoteka.Film
	if err := f.db.Select(&films, query, limit, limit*(page-1)); err != nil {
		return nil, err
	}
	return films, nil
}

func (f *FilmDao) GetCurFilm(id int) (filmoteka.Film, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", configs.EnvFilmTable())

	var film []filmoteka.Film
	if err := f.db.Select(&film, query, id); err != nil {
		return filmoteka.Film{}, err
	} else if len(film) == 0 {
		return filmoteka.Film{}, errors.New("no film find")
	}
	return film[0], nil
}

func (f *FilmDao) GetActorsInCurFilm(filmId int) ([]filmoteka.InputActor, error) {
	query := fmt.Sprintf("SELECT a.name, a.surname FROM %s a INNER JOIN %s s ON a.id = s.actor_id WHERE s.film_id=$1",
		configs.EnvActorTable(), configs.EnvStarredTable())

	var actors []filmoteka.InputActor
	if err := f.db.Select(&actors, query, filmId); err != nil {
		return nil, err
	}
	return actors, nil
}

func (f *FilmDao) GetFilmListByActor(page, limit int, fragment filmoteka.FilmSearchFragment) ([]filmoteka.Film, error) {
	query := fmt.Sprintf(searchQuery, configs.EnvStarredTable(), configs.EnvFilmTable(),
		configs.EnvActorTable(), getFilmSearchCondition(fragment))

	var films []filmoteka.Film
	if err := f.db.Select(&films, query, limit, limit*(page-1)); err != nil {
		return nil, err
	}
	return films, nil
}

func (f *FilmDao) DeleteFilmById(tx *sqlx.Tx, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", configs.EnvFilmTable())

	if _, err := tx.Exec(query, id); err != nil {
		return err
	}

	return nil
}
