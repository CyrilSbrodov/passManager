package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/server/internal/models"
	"github.com/CyrilSbrodov/passManager.git/server/pkg/client/postgres"
)

type Store struct {
	client postgres.Client
	Hash   string
	logger loggers.Logger
}

func createTable(ctx context.Context, client postgres.Client, logger *loggers.Logger) error {
	return nil
}

func NewStore(client postgres.Client, cfg *config.Config, logger *loggers.Logger) (*Store, error) {
	return &Store{}, nil
}

func (s *Store) CollectData(d *models.Data) error {
	return nil
}

func (s *Store) GetAllData() ([]models.Data, error) {
	return nil, nil
}

func (s *Store) UpdateData(d *models.Data) error {
	return nil
}

func (s *Store) DeleteData(d *models.Data) error {
	return nil
}
