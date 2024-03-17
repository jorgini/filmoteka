package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/jorgini/filmoteka/app"
	"github.com/jorgini/filmoteka/app/configs"
	"github.com/jorgini/filmoteka/app/handlers"
	"github.com/jorgini/filmoteka/app/models_dao"
	"github.com/jorgini/filmoteka/app/service"
	"github.com/sirupsen/logrus"
)

func main() {
	//logrus.SetFormatter(new(logrus.JSONFormatter))

	server := new(app.Server)

	configs.ConnectToDb()
	defer func(DB *sqlx.DB) {
		err := DB.Close()
		if err != nil {

		}
	}(configs.PsClient.DB)

	repo := models_dao.NewRepository(configs.PsClient.DB)
	serv := service.NewService(repo)
	router := handlers.NewRouter(serv)
	if err := server.Run(router); err != nil {
		logrus.Fatal(err)
		return
	}
}
