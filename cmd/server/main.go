package main

import (
	"fmt"
	"net/http"

	"github.com/dionaditya/go-production-ready-api/internal/database"
	transportHTTP "github.com/dionaditya/go-production-ready-api/internal/transport"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type App struct {
	Name    string
	version string
}

func (app *App) Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"AppName": app.Name,
		"Version": app.version,
	}).Info("Setting our app")

	var err error
	db, err := database.NewDatabase()

	if err != nil {
		return err
	}

	err = database.MigrateDB(db)

	if err != nil {
		log.Error("Failed to set up database")
	}

	handler := transportHTTP.NewHandler(db)

	handler.SetupRoutes()

	if err := http.ListenAndServe(":8000", handler.Router); err != nil {
		log.Error("failed to set up server")
	}
	return nil
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Error("failed to parse .env")
	}

	app := App{
		Name:    "comment-api",
		version: "1.0.0",
	}

	if err := app.Run(); err != nil {
		log.Error("failed to start app")
		fmt.Println(err)
	}
}
