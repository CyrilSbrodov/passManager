package crypto

import (
	"bytes"
	"crypto/rsa"

	"github.com/CyrilSbrodov/passManager.git/client/cmd/loggers"
)

type Crypto interface {
	AddCryptoKey(filenamePublicKey, filenamePrivateKey, filenameCert, path string) error
	CreateNewCryptoFile(PEM bytes.Buffer, filename, path string) error
	DecryptedData(b []byte, privateKey *rsa.PrivateKey) ([]byte, error)
	EncryptedData(b []byte, publicKey *rsa.PublicKey) ([]byte, error)
	LoadPrivatePEMKey(filename string) (*rsa.PrivateKey, error)
	LoadPublicPEMKey(filename string) (*rsa.PublicKey, error)
}

type RSA struct {
	logger loggers.Logger
}

func NewRSA() *RSA {
	return &RSA{}
}

func (r *RSA) AddCryptoKey(filenamePublicKey, filenamePrivateKey, filenameCert, path string) error {
	//TODO implement me
	panic("implement me")
}

func (r *RSA) CreateNewCryptoFile(PEM bytes.Buffer, filename, path string) error {
	//TODO implement me
	panic("implement me")
}

func (r *RSA) DecryptedData(b []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RSA) EncryptedData(b []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RSA) LoadPrivatePEMKey(filename string) (*rsa.PrivateKey, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RSA) LoadPublicPEMKey(filename string) (*rsa.PublicKey, error) {
	//TODO implement me
	panic("implement me")
}
