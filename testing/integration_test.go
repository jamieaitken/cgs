package testing

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jamieaitken/cgs"
	"github.com/jamieaitken/cgs/mysql"
	"github.com/jamieaitken/cgs/redis"
	instrRedis "github.com/jamieaitken/promred/redis"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.uber.org/zap"
)

var (
	redisDB   instrRedis.Redis
	redisAddr []string
	sqlDB     *sql.DB
	sqlDBAddr string
)

func TestMain(m *testing.M) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("failed to instantiate pool", zap.Error(err))
	}

	pool.MaxWait = time.Minute * 1

	redisResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Env: []string{
			"ALLOW_EMPTY_PASSWORD=yes",
		},
		Repository: "bitnami/redis",
		Tag:        "6.2.6",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	if err != nil {
		log.Fatal("failed to run redis container", zap.Error(err))
	}

	sqlResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Env: []string{
			"MYSQL_DATABASE=mysql",
			"MYSQL_ROOT_PASSWORD=secret",
		},
		Repository: "bitnami/mysql",
		Tag:        "5.7",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	if err != nil {
		log.Fatal("failed to run mysql container", zap.Error(err))
	}

	sqlDBAddr = fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", sqlResource.GetPort("3306/tcp"))

	db, err := mysql.New(sqlDBAddr)
	if err != nil {
		log.Fatal("failed to open connection to mysql store", zap.Error(err))
	}

	sqlDB = db.Client()

	redisAddr = []string{fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp"))}

	r := redis.New(redisAddr, redis.WithClientType(redis.NonFailOver))

	redisDB = r.Client()

	err = pool.Retry(func() error {
		pingErr := redisDB.Ping(context.Background(), "dockertest-backoff").Err()
		if pingErr != nil {
			return pingErr
		}

		pingErr = sqlDB.PingContext(context.Background())
		if pingErr != nil {
			return pingErr
		}

		return nil
	})
	if err != nil {
		log.Fatal("could not connect to container", zap.Error(err))
	}

	code := m.Run()

	err = cleanUp(pool, redisResource, sqlResource)
	if err != nil {
		log.Fatal("could not purge resource", zap.Error(err))
	}

	os.Exit(code)
}

func cleanUp(pool *dockertest.Pool, resources ...*dockertest.Resource) error {
	for _, resource := range resources {
		err := pool.Purge(resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestWithRedis(t *testing.T) {
	tests := []struct {
		name            string
		givenClientName string
	}{
		{
			name:            "given a correct redis address, expect a client to be made available that can ping",
			givenClientName: "narrow-test",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New(
				cgs.WithRedis(context.Background(), test.givenClientName, redisAddr,
					redis.WithClientType(redis.NonFailOver),
				),
			)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			c, err := app.Redis(test.givenClientName)
			if err != nil {
				t.Fatalf("expected nil for %s, got %v", test.givenClientName, err)
			}

			err = c.Client().Ping(context.Background(), test.givenClientName).Err()
			if err != nil {
				t.Fatalf("failed to ping redis for %s", test.givenClientName)
			}
		})
	}
}

func TestWithSQL(t *testing.T) {
	tests := []struct {
		name            string
		givenClientName string
	}{
		{
			name:            "given a correct sql address, expect a client to be made available that can ping",
			givenClientName: "narrow-test",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New(
				cgs.WithMySQL(context.Background(), test.givenClientName, sqlDBAddr),
			)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			c, err := app.MySQL(test.givenClientName)
			if err != nil {
				t.Fatalf("expected nil for %s, got %v", test.givenClientName, err)
			}

			err = c.Client().PingContext(context.Background())
			if err != nil {
				t.Fatalf("failed to ping database for %s", test.givenClientName)
			}
		})
	}
}
