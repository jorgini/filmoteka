package models_dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jorgini/filmoteka/app"
	"github.com/jorgini/filmoteka/app/configs"
	"strings"
)

type ActorDao struct {
	db *sqlx.DB
}

func NewActorDao(db *sqlx.DB) *ActorDao {
	return &ActorDao{
		db: db,
	}
}

func getSearchCondition(name, surname *string) string {
	condition := make([]string, 0, 2)
	if name != nil {
		condition = append(condition, fmt.Sprintf("POSITION('%s' in name)>0", *name))
	}
	if surname != nil {
		condition = append(condition, fmt.Sprintf("POSITION('%s' in surname)>0", *surname))
	}
	return strings.Join(condition, " AND ")
}

func (a *ActorDao) CreateActor(tx *sqlx.Tx, actor app.Actor) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, surname, sex, birthday) values ($1, $2, $3, TO_DATE($4,'DD-MM-YYYY')) RETURNING id",
		configs.EnvActorTable())

	var id int
	row := tx.QueryRow(query, actor.Name, actor.Surname, actor.Sex, actor.Birthday.String())
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a *ActorDao) UpdateActor(tx *sqlx.Tx, actor app.UpdateActorInput) error {
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", configs.EnvActorTable(), actor.GetValuesUpdate(), *actor.Id)

	if _, err := tx.Query(query, actor.GetArgsUpdate()...); err != nil {
		return err
	}
	return nil
}

func (a *ActorDao) GetActorId(name, surname string) (int, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1 AND surname=$2", configs.EnvActorTable())

	var id int
	row := a.db.QueryRow(query, name, surname)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a *ActorDao) GetActorById(id int) (app.Actor, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", configs.EnvActorTable())

	var actor app.Actor
	err := a.db.Get(&actor, query, id)
	if err != nil {
		return app.Actor{}, err
	}
	return actor, nil
}

func (a *ActorDao) GetActorsList(page, limit int) ([]app.Actor, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", configs.EnvActorTable())

	var actors []app.Actor
	if err := a.db.Select(&actors, query, limit, (page-1)*limit); err != nil {
		return nil, err
	}

	return actors, nil
}

func (a *ActorDao) SearchActor(page, limit int, name, surname *string) ([]app.Actor, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT %d OFFSET %d",
		configs.EnvActorTable(), getSearchCondition(name, surname), limit, limit*(page-1))

	var actors []app.Actor
	if err := a.db.Select(&actors, query); err != nil {
		return nil, err
	}
	return actors, nil
}

func (a *ActorDao) GetFilmsWithCurActor(actorId int) ([]app.Film, error) {
	query := fmt.Sprintf("SELECT f.* FROM %s f, %s s WHERE s.film_id=f.id AND s.actor_id=$1",
		configs.EnvFilmTable(), configs.EnvStarredTable())

	var films []app.Film
	if err := a.db.Select(&films, query, actorId); err != nil {
		return nil, err
	}
	return films, nil
}

func (a *ActorDao) DeleteActorById(tx *sqlx.Tx, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", configs.EnvActorTable())

	if _, err := tx.Query(query, id); err != nil {
		return err
	}
	return nil
}
