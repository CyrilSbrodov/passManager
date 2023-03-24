package repositories

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/server/pkg/client/postgres"
)

func TestPGStore(t *testing.T, cfg config.Config) (*Store, func(...string)) {
	t.Helper()

	cfg.DatabaseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	logger := loggers.NewLogger()
	client, err := postgres.NewClient(context.Background(), 5, &cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewStore(client, &cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			_, err = s.client.Exec(context.Background(), fmt.Sprintf(
				"TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			_, err = s.client.Exec(context.Background(), fmt.Sprintf(
				"DROP SCHEMA public CASCADE"))
			_, err = s.client.Exec(context.Background(), fmt.Sprintf(
				"CREATE SCHEMA public"))
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
