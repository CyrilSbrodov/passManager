package storage

import "github.com/CyrilSbrodov/passManager.git/server/internal/models"

type Storage interface {
	Register(u *models.User) error
	Login(u *models.User) error
	CollectData(d *models.Data) error
	UpdateData(d *models.Data) error
	DeleteData(d *models.Data) error
	GetAllData() ([]models.Data, error)
}
