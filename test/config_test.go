package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

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
	os.Setenv("SSO_BASE_URI", mockJwkStore.URL)

	MongoContainerRunner(m)
}

func MongoContainerRunner(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	pool.MaxWait = 3 * time.Minute

	var resource *dockertest.Resource
	mongoPortSpec := "27017/tcp"
	for i := 0; i < 3; i++ {
		selectedHostPort := "27017"
		mongoPortSpec = selectedHostPort + "/tcp"

		resource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Repository:   "mongo",
			Tag:          "8.0",
			Env:          []string{"MONGO_INITDB_DATABASE=catalogHistory"},
			Cmd:          []string{"--replSet", "rs0", "--bind_ip_all"},
			ExposedPorts: []string{"27017"},
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(mongoPortSpec): {{HostIP: "127.0.0.1", HostPort: selectedHostPort}},
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container is automatically removed
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "address already in use") && i < 2 {
			log.Printf("MongoDB test container failed to bind port, retrying startup: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Fatalf("Could not start resource: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	seeded := false
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		containerInfo, inspectErr := pool.Client.InspectContainer(resource.Container.ID)
		if inspectErr == nil && containerInfo != nil && !containerInfo.State.Running {
			return fmt.Errorf("mongodb container stopped: exitCode=%d error=%s", containerInfo.State.ExitCode, containerInfo.State.Error)
		}

		hostPort := strings.Replace(resource.GetHostPort(mongoPortSpec), "localhost", "127.0.0.1", 1)
		log.Printf("Using MongoDB test host port: %s", hostPort)
		_ = os.Setenv("MONGO_HOST", hostPort)
		_ = os.Setenv("MONGO_USERNAME", "")
		_ = os.Setenv("MONGO_PASSWORD", "")
		_ = os.Setenv("MONGODB_REPLICASET", "rs0")
		dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				"mongodb://" + hostPort + "/?directConnection=true",
			),
		)
		if err != nil {
			return err
		}

		admin := dbClient.Database("admin")
		rsInitiateCmd := bson.D{
			{Key: "replSetInitiate", Value: bson.D{
				{Key: "_id", Value: "rs0"},
				{Key: "members", Value: bson.A{
					bson.D{
						{Key: "_id", Value: 0},
						{Key: "host", Value: hostPort},
					},
				}},
			}},
		}
		if cmdErr := admin.RunCommand(context.TODO(), rsInitiateCmd).Err(); cmdErr != nil &&
			!strings.Contains(cmdErr.Error(), "already initialized") &&
			!strings.Contains(cmdErr.Error(), "already been initiated") {
			return cmdErr
		}

		hello := bson.D{{Key: "hello", Value: 1}}
		var helloResp bson.M
		if helloErr := admin.RunCommand(context.TODO(), hello).Decode(&helloResp); helloErr != nil {
			return helloErr
		}
		if isWritablePrimary, ok := helloResp["isWritablePrimary"].(bool); !ok || !isWritablePrimary {
			return fmt.Errorf("replica set not primary yet")
		}

		db := dbClient.Database("catalogHistory")
		coll := db.Collection("updates")
		if !seeded {
			if seedErr := seedUpdates(context.TODO(), coll); seedErr != nil {
				return seedErr
			}
			seeded = true
		}

		_, err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: "113"}}).Raw()
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

func seedUpdates(ctx context.Context, coll *mongo.Collection) error {
	updates := []interface{}{
		bson.M{
			"_id":        "123",
			"catalogId":  "111222333",
			"resourceId": "123456789",
			"person": bson.M{
				"id":    "123",
				"email": "example@example.com",
				"name":  "John Doe",
			},
			"datetime": time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			"operations": bson.A{
				bson.M{"op": "replace", "path": "/name", "value": "Jane"},
				bson.M{"op": "remove", "path": "/height"},
				bson.M{"op": "add", "path": "/name", "value": "Jane Test"},
			},
		},
		bson.M{
			"_id":        "789",
			"catalogId":  "111222333",
			"resourceId": "123456789",
			"person": bson.M{
				"id":    "789",
				"email": "example3@example.com",
				"name":  "Joe Doe",
			},
			"datetime": time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC),
			"operations": bson.A{
				bson.M{"op": "add", "path": "/name", "value": "Joe"},
			},
		},
		bson.M{
			"_id":        "456",
			"catalogId":  "111222333",
			"resourceId": "123456789",
			"person": bson.M{
				"id":    "456",
				"email": "example2@example.com",
				"name":  "Sarah Doe",
			},
			"datetime": time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC),
			"operations": bson.A{
				bson.M{"op": "replace", "path": "/name", "value": "Sarah"},
			},
		},
		bson.M{
			"_id":        "012",
			"catalogId":  "111222333",
			"resourceId": "123456789",
			"person": bson.M{
				"id":    "012",
				"email": "example4@example.com",
				"name":  "Bob Doe",
			},
			"datetime": time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC),
			"operations": bson.A{
				bson.M{"op": "replace", "path": "/name", "value": "Bob"},
			},
		},
		bson.M{
			"_id":        "113",
			"catalogId":  "123456789",
			"resourceId": "112",
			"person": bson.M{
				"id":    "110",
				"email": "example@example.com",
				"name":  "Doe Doe",
			},
			"datetime": time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC),
			"operations": bson.A{
				bson.M{"op": "replace", "path": "/name", "value": "Bob"},
			},
		},
	}

	_, err := coll.InsertMany(ctx, updates)
	return err
}
