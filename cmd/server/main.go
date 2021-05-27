package main

import (
	"fmt"
	"net/http"

	"github.com/dionaditya/go-production-ready-api/internal/comment"
	"github.com/dionaditya/go-production-ready-api/internal/database"
	transportHTTP "github.com/dionaditya/go-production-ready-api/internal/transport"
	"github.com/joho/godotenv"
)

type App struct{}

func (app *App) Run() error {
	fmt.Println("Setting up our REST API")

	var err error
	db, err := database.NewDatabase()

	if err != nil {
		return err
	}

	err = database.MigrateDB(db)

	if err != nil {
		fmt.Println("failed auto migrate")
	}

	commentService := comment.NewService(db)

	handler := transportHTTP.NewHandler(commentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8000", handler.Router); err != nil {
		fmt.Println("Failed setup routes")
	}
	return nil
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("failed to load env")
	}
	app := App{}

	if err := app.Run(); err != nil {
		fmt.Println("Error starting our app")
		fmt.Println(err)
	}
}
