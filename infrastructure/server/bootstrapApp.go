package server

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"persons.com/api/domain/person"
	httpHandler "persons.com/api/infrastructure/api/http"
	"persons.com/api/infrastructure/cache/redis"
	"persons.com/api/infrastructure/env"
	"persons.com/api/infrastructure/events/rabbitmq"
	"persons.com/api/infrastructure/repositories/mongodb"
	"persons.com/api/infrastructure/repositories/mysql"
)

var envMap map[string]string = env.NewEnvService().GetEnvs(os.Getenv("APP_MODE"))

func StartRouter() *chi.Mux {
	redisCacheService, err := redis.GetRedisClient(envMap["CACHE_DB_URL"], 30)
	if err != nil {
		log.Fatal(err)
	}

	rabbitEventsService, err := rabbitmq.NewRabbitMqService(envMap["Q_URL"], []string{"persons"}, []string{"minors", "adults"})
	if err != nil {
		log.Fatal(err)
	}

	repository := getRepository()
	service := person.NewPersonService(repository)
	handler := httpHandler.NewHandler(service, rabbitEventsService, redisCacheService)

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	})

	router := chi.NewRouter()
	router.Use(cors.Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/api/{id}", handler.GetById)
	router.Get("/api", handler.GetAll)
	router.Post("/api", handler.Create)

	return router
}

func HttpPort() string {
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
