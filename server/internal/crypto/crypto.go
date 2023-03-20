package crypto

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
)

type Crypto interface {
	//AddCryptoKey(filenamePublicKey, filenamePrivateKey, filenameCert, path string) error
	//CreateNewCryptoFile(PEM bytes.Buffer, filename, path string) error
	DecryptedData(b []byte, privateKey *rsa.PrivateKey) ([]byte, error)
	EncryptedData(b []byte, publicKey *rsa.PublicKey) ([]byte, error)
	LoadPrivatePEMKey(filename string) (*rsa.PrivateKey, error)
	LoadPublicPEMKey(filename string) (*rsa.PublicKey, error)
}

type RSA struct {
	logger  *loggers.Logger
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func NewRSA(cfg config.Config) *RSA {
	logger := loggers.NewLogger()
	private, public, err := addCryptoKey("public.pem", cfg.CryptoPROKey, cfg.CryptoPROKeyPath, logger)
	if err != nil {
		logger.LogErr(err, "")
		os.Exit(1)
	}
	return &RSA{
		logger:  logger,
		Private: private,
		Public:  public,
	}
}

func addCryptoKey(filenamePublicKey, filenamePrivateKey, path string, logger *loggers.Logger) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"metricService"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.LogErr(err, "")
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		logger.LogErr(err, "")
	}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	var publicKeyPEM bytes.Buffer
	pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	publicKey := privateKey.PublicKey
	if err = createNewCryptoFile(certPEM, "cert.pem", path, logger); err != nil {
		logger.LogErr(err, "filed to create new file")
		return nil, nil, err
	}
	if err = createNewCryptoFile(privateKeyPEM, filenamePrivateKey, path, logger); err != nil {
		logger.LogErr(err, "filed to create new file")
		return nil, nil, err
	}
	if err = createNewCryptoFile(publicKeyPEM, filenamePublicKey, path, logger); err != nil {
		logger.LogErr(err, "filed to create new file")
		return nil, nil, err
	}
	return privateKey, &publicKey, nil
}

func createNewCryptoFile(PEM bytes.Buffer, filename, path string, logger *loggers.Logger) error {
	file, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		logger.LogErr(err, "failed to open/create file")
		return err
	}
	writer := bufio.NewWriter(file)
	if _, err := writer.Write(PEM.Bytes()); err != nil {
		logger.LogErr(err, "failed to write buffer")
		return err
	}
	writer.Flush()
	defer file.Close()
	return nil
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
