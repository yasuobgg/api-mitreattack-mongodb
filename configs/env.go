package configs

import (
    "log"
    "os"
    "fmt"

    "github.com/joho/godotenv" 
)

func EnvMongoURI() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    dbHost, dbPort, dbName := os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME")

	connection_string := fmt.Sprintf("mongodb://%s:%s/%s", dbHost, dbPort, dbName)

    return connection_string
}