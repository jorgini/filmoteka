package configs

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

type PostgresInstance struct {
	DB *sqlx.DB
}

var PsClient PostgresInstance

func ConnectToDb() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		EnvHost(), EnvPortDb(), EnvUserDB(), EnvPassword(), EnvDB())

	var err error
	PsClient.DB, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		logrus.Fatalf("connecting to db failed with: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = PsClient.DB.PingContext(ctx)
	if err != nil {
		logrus.Errorf("error occur on ping: %s", err.Error())
		err := PsClient.DB.Close()
		if err != nil {
			logrus.Fatal(err)
		}
		panic(err)
	}

	logrus.Info("Database connected successfully...")
}
