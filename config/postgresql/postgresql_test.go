package postgresql

import (
	"net/url"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionString(t *testing.T) {
	t.Setenv("POSTGRESQL_HOST", "postgresql-sqlproxy")
	t.Setenv("POSTGRESQL_PORT", "5432")
	t.Setenv("POSTGRESQL_DB", "catalog_history_service_staging")
	t.Setenv("POSTGRESQL_USERNAME", "catalog_history_service_staging")
	t.Setenv("POSTGRESQL_PASSWORD", "p@ss:wo/rd%")

	cs := ConnectionString()

	_, err := url.Parse(cs)
	require.NoError(t, err)

	_, err = pgxpool.ParseConfig(cs)
	require.NoError(t, err)

	assert.Contains(t, cs, "postgresql-sqlproxy:5432")
	assert.Contains(t, cs, "sslmode=disable")
}

func TestConnectionString_trimsSecretWhitespace(t *testing.T) {
	t.Setenv("POSTGRESQL_HOST", "host\n")
	t.Setenv("POSTGRESQL_PORT", "5432")
	t.Setenv("POSTGRESQL_DB", "db")
	t.Setenv("POSTGRESQL_USERNAME", "user\n")
	t.Setenv("POSTGRESQL_PASSWORD", "secret\n")

	cs := ConnectionString()
	_, err := url.Parse(cs)
	require.NoError(t, err)
}
