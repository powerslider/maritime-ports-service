package main

import (
	"context"
	"log"
	"os"

	"github.com/powerslider/maritime-ports-service/pkg/portsmanaging"
	"github.com/powerslider/maritime-ports-service/pkg/transport/server"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/powerslider/maritime-ports-service/pkg/configs"
	"github.com/powerslider/maritime-ports-service/pkg/handlers"
	"github.com/powerslider/maritime-ports-service/pkg/storage/memory"
)

// @title Maritime Ports Service API
// @version 1.0
// @description API for maritime ports data.
// @termsOfService http://swagger.io/terms/

// @contact.name Tsvetan Dimitrov
// @contact.email tsvetan.dimitrov23@gmail.com

// @license.name MIT
// @license.url https://www.mit.edu/~amini/LICENSE.md

// @host 0.0.0.0:8080
// @BasePath /
func main() {
	var err error

	setEnvironment()

	portsStore := memory.NewPortsRepository()
	loader := portsmanaging.NewJSONLoader(portsStore)

	err = loader.LoadJSONFile("./fixtures/ports.json")
	if err != nil {
		log.Fatalf("cannot seed service database with ports data: %v", err)
	}

	ctx := context.Background()

	conf := configs.InitializeConfig()

	portsService := portsmanaging.NewService(portsStore)

	router := mux.NewRouter()
	router = handlers.InitializeHandlers(conf, router, portsService)

	s := server.NewServer(conf, router)
	if err = s.Run(ctx); err != nil {
		log.Fatal(err.Error())
	}
}

func setEnvironment() {
	_, foundHost := os.LookupEnv("SERVER_HOST")
	_, foundPort := os.LookupEnv("SERVER_PORT")

	if !foundHost && !foundPort {
		err := godotenv.Load(".env.dist")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
