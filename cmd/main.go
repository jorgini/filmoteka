package main

import (
	"context"
	"github.com/jorgini/filmoteka"
	"github.com/jorgini/filmoteka/configs"
	"github.com/jorgini/filmoteka/handlers"
	"github.com/jorgini/filmoteka/models_dao"
	"github.com/jorgini/filmoteka/service"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	server := new(filmoteka.Server)

	configs.ConnectToDb()

	repo := models_dao.NewRepository(configs.PsClient.DB)
	serv := service.NewService(repo)
	router := handlers.NewRouter(serv)

	go func() {
		if err := server.Run(router); err != nil {
			logrus.Fatal(err)
			return
		}
	}()

	logrus.Info("Filmoteka-app successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("Filmoteka-app shut down")

	if err := server.ShutDown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shut down: %s", err.Error())
	}

	if err := configs.PsClient.DB.Close(); err != nil {
		logrus.Errorf("error occured on closing db connection: %s", err.Error())
	}
}
