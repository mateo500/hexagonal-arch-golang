package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"persons.com/api/domain/person"
	"persons.com/api/infrastructure/api"
	"persons.com/api/infrastructure/repository/mysql"
	"persons.com/api/infrastructure/repository/redis"
)

func main() {
	repository := getRepository() // repository <- (domain -> service) -> handler & serializer -> Transporter(http, grpc, soap, websockets etc.)
	service := person.NewPersonService(repository)
	handler := api.NewHandler(service)

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
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func getRepository() person.PersonRepository {
	switch os.Getenv("DB_TYPE") {
	case "redis":
		redisUrl := os.Getenv("DB_URL")
		repository, err := redis.NewRedisRepository(redisUrl)
		if err != nil {
			log.Fatal(err)
		}
		return repository

	case "mysql":
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		database, err := mysql.NewMysqlClient(dbUser, dbPass, dbName)
		if err != nil {
			log.Fatal(err)
		}

		return mysql.NewMysqlRepository(database)
	}

	return nil
}
