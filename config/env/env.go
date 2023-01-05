package env

import "os"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func MongoHost() string {
	return getEnv("MONGO_HOST", "localhost:27017")
}

func MongoPassword() string {
	return getEnv("MONGO_PASSWORD", "admin")
}

func MongoUsername() string {
	return getEnv("MONGO_USERNAME", "admin")
}

type Constants struct {
	MongoAuthParams string
	MongoCollection string
	MongoDatabase   string
}

type Paths struct {
	Datasources string
	Ping        string
	Ready       string
}

var ConstantValues = Constants{
	MongoAuthParams: "authSource=admin&authMechanism=SCRAM-SHA-1",
	MongoCollection: "begreper",
	MongoDatabase:   "catalog-history-service",
}

var PathValues = Paths{
	Datasources: "begreper",
	Ping:        "ping",
	Ready:       "ready",
}
