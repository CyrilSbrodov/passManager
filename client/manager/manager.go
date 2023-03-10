package manager

import (
	"crypto/rsa"

	"github.com/CyrilSbrodov/passManager.git/client/cmd/loggers"
)

type Manager struct {
	Text                string
	Data                string
	Logger              loggers.Logger
	PrivateKey          *rsa.PrivateKey
	PublicKey           *rsa.PublicKey
	PublicKeyFromServer *rsa.PublicKey
}

type Managers interface {
	Register(login, password string) error
	Auth(login, password string) ([]byte, error)
	AddData(s string) error
	UpdateData(s string) error
	DeleteData(s string) error
	GetData() (string, error)
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Register(login, password string) error {
	return nil
}

func (m *Manager) Auth(login, password string) ([]byte, error) {
	return nil, nil
}

func (m *Manager) AddData(s string) error {
	return nil
}

func (m *Manager) UpdateData(s string) error {
	return nil
}

func (m *Manager) DeleteData(s string) error {
	return nil
}

func (m *Manager) GetData() (string, error) {
	return "", nil
}
