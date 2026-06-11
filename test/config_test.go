package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/postgresql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type TestConstants struct {
	Audience     []string
	SysAdminAuth string
}

var TestValues = TestConstants{
	Audience:     []string{"catalog-history-service"},
	SysAdminAuth: "system:root:admin",
}

func OrgAdminAuth(org string) string {
	return fmt.Sprintf("organization:%s:admin", org)
}

func OrgWriteAuth(org string) string {
	return fmt.Sprintf("organization:%s:write", org)
}

func OrgReadAuth(org string) string {
	return fmt.Sprintf("organization:%s:read", org)
}

func TestMain(m *testing.M) {
	mockJwkStore := MockJwkStore()
	os.Setenv("SSO_BASE_URI", mockJwkStore.URL)

	PostgresContainerRunner(m)
}

func PostgresContainerRunner(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	pool.MaxWait = 2 * time.Minute

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpassword",
			"POSTGRES_DB=catalog_history",
		},
		ExposedPorts: []string{"5432/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "127.0.0.1", HostPort: "0"}},
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}

	hostPort := resource.GetHostPort("5432/tcp")

	var pgPool *pgxpool.Pool
	err = pool.Retry(func() error {
		connStr := fmt.Sprintf("postgres://testuser:testpassword@%s/catalog_history?sslmode=disable", hostPort)
		pgPool, err = pgxpool.New(context.Background(), connStr)
		if err != nil {
			return err
		}
		return pgPool.Ping(context.Background())
	})
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err)
	}

	os.Setenv("POSTGRESQL_HOST", "127.0.0.1")
	os.Setenv("POSTGRESQL_PORT", resource.GetPort("5432/tcp"))
	os.Setenv("POSTGRESQL_DB", "catalog_history")
	os.Setenv("POSTGRESQL_USERNAME", "testuser")
	os.Setenv("POSTGRESQL_PASSWORD", "testpassword")

	postgresql.SetPool(pgPool)

	if _, err := pgPool.Exec(context.Background(), createTableSQL); err != nil {
		log.Fatalf("Could not create schema: %s", err)
	}

	if err := seedUpdates(context.Background(), pgPool); err != nil {
		log.Fatalf("Could not seed data: %s", err)
	}

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	pgPool.Close()

	os.Exit(code)
}

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

type seedOperation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
}

func seedUpdates(ctx context.Context, pool *pgxpool.Pool) error {
	type seedRow struct {
		id          string
		catalogId   string
		resourceId  string
		personId    string
		personEmail string
		personName  string
		datetime    time.Time
		operations  []seedOperation
	}

	rows := []seedRow{
		{
			id: "123", catalogId: "111222333", resourceId: "123456789",
			personId: "123", personEmail: "example@example.com", personName: "John Doe",
			datetime: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			operations: []seedOperation{
				{Op: "replace", Path: "/name", Value: "Jane"},
				{Op: "remove", Path: "/height"},
				{Op: "add", Path: "/name", Value: "Jane Test"},
			},
		},
		{
			id: "789", catalogId: "111222333", resourceId: "123456789",
			personId: "789", personEmail: "example3@example.com", personName: "Joe Doe",
			datetime: time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC),
			operations: []seedOperation{
				{Op: "add", Path: "/name", Value: "Joe"},
			},
		},
		{
			id: "456", catalogId: "111222333", resourceId: "123456789",
			personId: "456", personEmail: "example2@example.com", personName: "Sarah Doe",
			datetime: time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC),
			operations: []seedOperation{
				{Op: "replace", Path: "/name", Value: "Sarah"},
			},
		},
		{
			id: "012", catalogId: "111222333", resourceId: "123456789",
			personId: "012", personEmail: "example4@example.com", personName: "Bob Doe",
			datetime: time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC),
			operations: []seedOperation{
				{Op: "replace", Path: "/name", Value: "Bob"},
			},
		},
		{
			id: "113", catalogId: "123456789", resourceId: "112",
			personId: "110", personEmail: "example@example.com", personName: "Doe Doe",
			datetime: time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC),
			operations: []seedOperation{
				{Op: "replace", Path: "/name", Value: "Bob"},
			},
		},
	}

	for _, r := range rows {
		opsJSON, err := json.Marshal(r.operations)
		if err != nil {
			return fmt.Errorf("failed to marshal operations for %s: %w", r.id, err)
		}

		_, err = pool.Exec(ctx,
			`INSERT INTO updates (id, catalog_id, resource_id, person_id, person_email, person_name, datetime, operations)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			r.id, r.catalogId, r.resourceId,
			r.personId, r.personEmail, r.personName,
			r.datetime, opsJSON,
		)
		if err != nil {
			return fmt.Errorf("failed to insert seed row %s: %w", r.id, err)
		}
	}

	return nil
}
