package postgresql

import (
	"context"
	"net"
	"net/url"
	"strings"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/env"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS updates (
    id            VARCHAR(255) PRIMARY KEY,
    catalog_id    VARCHAR(255) NOT NULL,
    resource_id   VARCHAR(255) NOT NULL,
    person_id     VARCHAR(255) NOT NULL,
    person_email  VARCHAR(255) NOT NULL,
    person_name   VARCHAR(255) NOT NULL,
    datetime      TIMESTAMPTZ,
    operations    JSONB        NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_updates_catalog_resource ON updates (catalog_id, resource_id);
CREATE INDEX IF NOT EXISTS idx_updates_catalog ON updates (catalog_id);
`

func ConnectionString() string {
	user := strings.TrimSpace(env.PostgresUsername())
	password := strings.TrimSpace(env.PostgresPassword())
	host := strings.TrimSpace(env.PostgresHost())
	port := strings.TrimSpace(env.PostgresPort())
	db := strings.TrimSpace(env.PostgresDB())

	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, password),
		Host:     net.JoinHostPort(host, port),
		Path:     "/" + db,
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

var pool *pgxpool.Pool

func Pool() *pgxpool.Pool {
	if pool == nil {
		var err error
		pool, err = pgxpool.New(context.Background(), ConnectionString())
		if err != nil {
			logrus.Fatalf("Unable to create connection pool: %v", err)
		}

		if _, err := pool.Exec(context.Background(), createTableSQL); err != nil {
			logrus.Errorf("Failed to initialize schema: %v", err)
			logging.LogAndPrintError(err)
		}
	}
	return pool
}

func SetPool(p *pgxpool.Pool) {
	pool = p
}
