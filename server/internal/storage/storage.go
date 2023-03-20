package storage

import "github.com/CyrilSbrodov/passManager.git/server/internal/models"

type Storage interface {
	Register(u *models.User) (string, error)
	Login(u *models.User) (string, error)
	CollectCard(d *models.CryptoCard, login string) (int, error)
	CollectPassword(d *models.CryptoPassword, login string) (int, error)
	CollectText(d *models.CryptoTextData, login string) (int, error)
	CollectBinary(d *models.CryptoBinaryData, login string) (int, error)
	GetCards(id string) (int, []models.CryptoCard, error)
	GetPassword(id string) (int, []models.CryptoPassword, error)
	GetText(id string) (int, []models.CryptoTextData, error)
	GetBinary(id string) (int, []models.CryptoBinaryData, error)
	DeleteCard(d *models.CryptoCard, id string) (int, error)
	DeleteText(d *models.CryptoTextData, id string) (int, error)
	DeletePassword(d *models.CryptoPassword, id string) (int, error)
	DeleteBinary(d *models.CryptoBinaryData, id string) (int, error)
	UpdateCard(d *models.CryptoCard, id string) (int, error)
	UpdatePassword(d *models.CryptoPassword, id string) (int, error)
	UpdateText(d *models.CryptoTextData, id string) (int, error)
	UpdateBinary(d *models.CryptoBinaryData, id string) (int, error)
}
