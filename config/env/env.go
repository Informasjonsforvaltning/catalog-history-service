package env

import (
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func CorsOriginPatterns() []string {
	return strings.Split(getEnv("CORS_ORIGIN_PATTERNS", "*"), ",")
}

func MongoHost() string {
	return getEnv("MONGO_HOST", "localhost:27017")
}

func MongoPassword() string {
	return getEnv("MONGO_PASSWORD", "admin")
}

func MongoUsername() string {
	return getEnv("MONGO_USERNAME", "root")
}

func MongoAuthSource() string {
	return getEnv("MONGODB_AUTH", "admin")
}

func MongoReplicaSet() string {
	return getEnv("MONGODB_REPLICASET", "replicaset")
}

func KeycloakHost() string {
	return getEnv("SSO_BASE_URI", "https://auth.staging.fellesdatakatalog.digdir.no")
}

type Constants struct {
	MongoAuthParams string
	MongoCollection string
	MongoDatabase   string
}

type Paths struct {
	Resource       string
	ResourceUpdate string
	Ping           string
	Ready          string
}

type Security struct {
	TokenAudience   string
	SysAdminAuth    string
	OrgType         string
	AdminPermission string
	WritePermission string
	ReadPermission  string
}

var ConstantValues = Constants{
	MongoAuthParams: "authSource=admin&authMechanism=SCRAM-SHA-1",
	MongoCollection: "updates",
	MongoDatabase:   "catalogHistory",
}

var PathValues = Paths{
	Resource:       "/:catalogId/:resourceId/updates",
	ResourceUpdate: "/:catalogId/:resourceId/updates/:updateId",
	Ping:           "/ping",
	Ready:          "/ready",
}

var SecurityValues = Security{
	TokenAudience:   "catalog-history-service",
	SysAdminAuth:    "system:root:admin",
	OrgType:         "organization",
	AdminPermission: "admin",
	WritePermission: "write",
	ReadPermission:  "read",
}
