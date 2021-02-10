package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"persons.com/api/domain/person"
	httpApi "persons.com/api/infrastructure/api/http"
	"persons.com/api/infrastructure/env"
	"persons.com/api/infrastructure/repositories/mongodb"
	"persons.com/api/infrastructure/repositories/mysql"
)

var envMap map[string]string = env.NewEnvService().GetEnvs("dev")

func main() {
	repository := getRepository()
	service := person.NewPersonService(repository)
	handler := httpApi.NewHandler(service)
	//app flow: Domain -> Service -> Repository -> Serializers(json, messagePack, grpc, soap, etc) -> Handlers(controllers) -> Transporter(http, websockets, GraphQl etc.)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{id}", handler.GetById)
	router.Get("/", handler.GetAll)
	router.Post("/", handler.Create)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :5000")
		errs <- http.ListenAndServe(httpPort(), router)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "5000"
	if envMap["APP_PORT"] != "" {
		port = envMap["APP_PORT"]
	}
	return fmt.Sprintf(":%s", port)
}

func getRepository() person.PersonRepository {

	switch os.Getenv("DB_TYPE") {
	case "mysql":
		dbUser := envMap["DB_USER"]
		dbPass := envMap["DB_PASSWORD"]
		dbName := envMap["DB_NAME"]

		database, err := mysql.NewMysqlClient(dbUser, dbPass, dbName)
		if err != nil {
			log.Fatal(err)
		}

		return mysql.NewMysqlRepository(database)

	case "mongodb":
		dbURL := envMap["DB_URL"]
		dbName := envMap["MONGO_DB"]
		mongoTimeout, err := strconv.Atoi(envMap["MONGO_TIMEOUT"])
		if err != nil {
			log.Fatal(err)
		}

		repository, err := mongodb.NewMongoRepository(dbURL, dbName, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}

		return repository

	}

	return nil
}
