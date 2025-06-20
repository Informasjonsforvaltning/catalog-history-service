package test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

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
	os.Setenv("SSO_AUTH_URI", mockJwkStore.URL)

	MongoContainerRunner(m)
}

func MongoContainerRunner(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get directory: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "bitnami/mongodb",
		Tag:        "latest",
		Env: []string{
			"MONGODB_ROOT_PASSWORD=admin",
			"MONGODB_ADVERTISED_HOSTNAME=localhost",
			"MONGODB_REPLICA_SET_MODE=primary",
			"MONGODB_REPLICA_SET_KEY=replicaset",
		},
		ExposedPorts: []string{"27017"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017": {{HostIP: "127.0.0.1", HostPort: "27017"}},
		},
		Mounts: []string{
			currentDirectory + "/init-mongo/init-mongo.js:/docker-entrypoint-initdb.d/init-mongo.js:ro",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container is automatically removed
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				"mongodb://root:admin@localhost:27017",
			),
		)
		if err != nil {
			return err
		}
		// try to find a document added in init-mongo file
		db := dbClient.Database("catalogHistory")
		coll := db.Collection("updates")
		_, err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: "123"}}).Raw()
		return err
	})
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// run tests
	code := m.Run()

	// kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client
	if err = dbClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	os.Exit(code)
}
