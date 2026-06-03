package postgresql

import (
	"context"
	"fmt"

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
    datetime      TIMESTAMPTZ  NOT NULL,
    operations    JSONB        NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_updates_catalog_resource ON updates (catalog_id, resource_id);
CREATE INDEX IF NOT EXISTS idx_updates_catalog ON updates (catalog_id);
`

func ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.PostgresUsername(),
		env.PostgresPassword(),
		env.PostgresHost(),
		env.PostgresPort(),
		env.PostgresDB(),
	)
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
