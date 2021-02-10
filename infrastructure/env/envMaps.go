package env

import "os"

var devMap map[string]string = map[string]string{
	"DB_TYPE":       os.Getenv("DB_TYPE"),
	"DB_USER":       os.Getenv("DB_USER"),
	"DB_PASSWORD":   os.Getenv("DB_PASSWORD"),
	"DB_NAME":       os.Getenv("DB_NAME"),
	"DB_URL":        os.Getenv("DB_URL"),
	"MONGO_DB":      os.Getenv("MONGO_DB"),
	"MONGO_TIMEOUT": os.Getenv("MONGO_TIMEOUT"),
	"CACHE_DB_URL":  os.Getenv("CACHE_DB_URL"),
	"APP_PORT":      os.Getenv("APP_PORT"),
	"Q_NAME":        os.Getenv("Q_NAME"),
	"Q_URL":         os.Getenv("Q_URL"),
}

var prodMap map[string]string = map[string]string{
	"DB_TYPE":       os.Getenv("DB_TYPE"),
	"DB_USER":       os.Getenv("DB_USER"),
	"DB_PASSWORD":   os.Getenv("DB_PASSWORD"),
	"DB_NAME":       os.Getenv("DB_NAME"),
	"DB_URL":        os.Getenv("DB_URL"),
	"MONGO_DB":      os.Getenv("MONGO_DB"),
	"MONGO_TIMEOUT": os.Getenv("MONGO_TIMEOUT"),
	"CACHE_DB_URL":  os.Getenv("CACHE_DB_URL"),
	"APP_PORT":      os.Getenv("APP_PORT"),
	"Q_NAME":        os.Getenv("Q_NAME"),
	"Q_URL":         os.Getenv("Q_URL"),
}
