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

func PostgresHost() string {
	return getEnv("POSTGRESQL_HOST", "localhost")
}

func PostgresPort() string {
	return getEnv("POSTGRESQL_PORT", "5432")
}

func PostgresDB() string {
	return getEnv("POSTGRESQL_DB", "catalog_history")
}

func PostgresUsername() string {
	return getEnv("POSTGRESQL_USERNAME", "admin")
}

func PostgresPassword() string {
	return getEnv("POSTGRESQL_PASSWORD", "admin")
}

func KeycloakHost() string {
	return getEnv("SSO_BASE_URI", "https://auth.staging.fellesdatakatalog.digdir.no")
}

type Paths struct {
	Resource       string
	ResourceUpdate string
	ConceptUpdates string
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

var PathValues = Paths{
	Resource:       "/:catalogId/:resourceId/updates",
	ResourceUpdate: "/:catalogId/:resourceId/updates/:updateId",
	ConceptUpdates: "/:catalogId/concepts/updates",
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
