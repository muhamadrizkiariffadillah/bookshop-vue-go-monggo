package config

import (
	"errors"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type ConfigDTO struct {
	port         string
	secret_key   string
	database_url string
}

var env ConfigDTO

// Load .env file and set values to env struct
func LoadEnvVariable() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	env = ConfigDTO{
		port:         os.Getenv("PORT"),
		secret_key:   os.Getenv("SECRET_KEY"),
		database_url: os.Getenv("MONGODB_URL"),
	}
}

func GetEnvProperties(key string) string {
	value, err := accessField(key)
	if err != nil || value == "" {
		log.Printf("Warning: %s is not set, using default value if applicable\n", key)
		return ""
	}
	return value
}

// Access struct field dynamically
func accessField(key string) (string, error) {
	value := reflect.ValueOf(env)
	t := value.Type()

	field, ok := t.FieldByName(key)
	if !ok {
		return "", errors.New("property not found: " + key)
	}

	return value.FieldByIndex(field.Index).String(), nil
}
