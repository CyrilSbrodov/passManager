package crypto

import (
	"bytes"
	"crypto/rsa"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
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
