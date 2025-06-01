package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func GetDBUser() string {
	return os.Getenv("DB_USER")
}

func GetDBPass() string {
	return os.Getenv("DB_PASS")
}

func GetDBHost() string {
	return os.Getenv("DB_HOST")
}

func GetDBPort() string {
	return os.Getenv("DB_PORT")
}

func GetDBName() string {
	return os.Getenv("DB_NAME")
}

func GetDBLog() string {
	return os.Getenv("DB_LOG")
}

func GetGinMode() string {
	return os.Getenv("GIN_MODE")
}

func GetPort() string {
	return os.Getenv("PORT")
}
