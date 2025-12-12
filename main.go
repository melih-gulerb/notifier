package main

import (
	"notifier/configs"
	"notifier/routes"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	appLogger, nr := configs.InitLogging()
	e := configs.InitEcho(nr)

	routes.SetupMailRoutes(e, appLogger)

	e.Logger.Fatal(e.Start(":4006"))
}
