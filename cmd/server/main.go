package main

import (
	"fmt"
	"net/http"

	transportHTTP "github.com/dionaditya/go-production-ready-api/internal/transport"
)

type App struct{}

func (app *App) Run() error {
	fmt.Println("Setting up our REST API")
	handler := transportHTTP.NewHandler()
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8000", handler.Router); err != nil {
		fmt.Println("Failed setup routes")
	}
	return nil
}

func main() {
	app := App{}

	if err := app.Run(); err != nil {
		fmt.Println("Error starting our app")
		fmt.Println(err)
	}
}
