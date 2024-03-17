package configs

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func EnvAddr() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("ADDR")
}

func EnvHost() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("HOST")
}

func EnvPortDb() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("PORTDB")
}

func EnvDB() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("DB")
}

func EnvUserDB() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("USERDB")
}

func EnvPassword() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("PASSWORD")
}

func EnvUserTable() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("USERTABLE")
}

func EnvActorTable() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("ACTORTABLE")
}

func EnvFilmTable() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("FILMTABLE")
}

func EnvStarredTable() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return os.Getenv("STARREDTABLE")
}
