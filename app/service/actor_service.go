package service

import (
	"github.com/jorgini/filmoteka/app"
	"github.com/jorgini/filmoteka/app/models_dao"
)

type ActorService struct {
	dao models_dao.Actor
	tx  models_dao.Transaction
}

func NewActorService(dao models_dao.Actor, tx models_dao.Transaction) *ActorService {
	return &ActorService{
		dao: dao,
		tx:  tx,
	}
}

func (a *ActorService) CreateActor(actor app.Actor) (int, error) {
	transaction, err := a.tx.StartTransaction()
	if err != nil {
		return 0, err
	}

	var id int
	id, err = a.dao.CreateActor(transaction, actor)
	if err != nil {
		return 0, a.tx.ShutDown(transaction, err)
	}
	return id, a.tx.Commit(transaction)
}

func (a *ActorService) UpdateActor(actor app.UpdateActorInput) error {
	transaction, err := a.tx.StartTransaction()
	if err != nil {
		return err
	}

	if err = a.dao.UpdateActor(transaction, actor); err != nil {
		return a.tx.ShutDown(transaction, err)
	}

	return a.tx.Commit(transaction)
}

func (a *ActorService) GetActorId(name, surname string) (int, error) {
	return a.dao.GetActorId(name, surname)
}

func (a *ActorService) GetActorsList(page, limit int) ([]app.ActorListItem, error) {
	actors, err := a.dao.GetActorsList(page, limit)
	if err != nil {
		return nil, err
	}

	list := make([]app.ActorListItem, len(actors))
	for i := range actors {
		list[i].Actor = actors[i]
		list[i].Films, err = a.dao.GetFilmsWithCurActor(actors[i].Id)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (a *ActorService) GetActorById(id int) (app.ActorListItem, error) {
	actor, err := a.dao.GetActorById(id)
	if err != nil {
		return app.ActorListItem{}, err
	}

	films, err := a.dao.GetFilmsWithCurActor(actor.Id)
	if err != nil {
		return app.ActorListItem{}, err
	}
	return app.ActorListItem{Actor: actor, Films: films}, nil
}

func (a *ActorService) SearchActor(page, limit int, fragment app.ActorSearchFragment) ([]app.ActorListItem, error) {
	actors, err := a.dao.SearchActor(page, limit, fragment.Name, fragment.Surname)
	if err != nil {
		return nil, err
	}

	list := make([]app.ActorListItem, len(actors))
	for i := range actors {
		list[i].Actor = actors[i]
		list[i].Films, err = a.dao.GetFilmsWithCurActor(actors[i].Id)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (a *ActorService) DeleteActorById(id int) error {
	transaction, err := a.tx.StartTransaction()
	if err != nil {
		return err
	}

	if err = a.dao.DeleteActorById(transaction, id); err != nil {
		return a.tx.ShutDown(transaction, err)
	}

	return a.tx.Commit(transaction)
}
