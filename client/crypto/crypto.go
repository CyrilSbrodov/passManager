package crypto

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/CyrilSbrodov/passManager.git/client/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/client/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/client/model"
)

type Crypto interface {
	DecryptedData(b []byte, privateKey *rsa.PrivateKey) ([]byte, error)
	EncryptedData(b []byte, publicKey *rsa.PublicKey) ([]byte, error)
	EncryptedCard(d *model.CryptoCard)
	EncryptedPassword(d *model.CryptoPassword)
	EncryptedTextData(d *model.CryptoTextData)
	EncryptedBinaryData(d *model.CryptoBinaryData)
	DecryptedCard(d *model.CryptoCard)
	DecryptedPassword(d *model.CryptoPassword)
	DecryptedTextData(d *model.CryptoTextData)
	DecryptedBinaryData(d *model.CryptoBinaryData)
}

type RSA struct {
	logger  *loggers.Logger
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func NewRSA(cfg config.Config, logger *loggers.Logger) *RSA {
	var private *rsa.PrivateKey
	var public *rsa.PublicKey
	var err error
	private, public, err = LoadPrivateAndPublicPEMKey(cfg.CryptoPROKey, "public.pem", logger)

	if err != nil {
		fmt.Println("new")
		private, public, err = addCryptoKey("public.pem", cfg.CryptoPROKey, cfg.CryptoPROKeyPath, logger)
		if err != nil {
			logger.LogErr(err, "")
			os.Exit(1)
		}
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
			Organization: []string{"passManager"},
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

func (r *RSA) DecryptedData(msg []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	msgLen := len(msg)
	step := privateKey.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(sha256.New(), rng, privateKey, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

func (r *RSA) EncryptedData(b []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	label := []byte("OAEP Encrypted")
	msgLen := len(b)
	step := publicKey.Size() - 2*sha256.New().Size() - 2
	rng := rand.Reader
	var encryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, b[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func LoadPrivateAndPublicPEMKey(private, public string, logger *loggers.Logger) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKeyFile, err := os.Open(private)
	if err != nil {
		logger.LogErr(err, "filed to open file")
		return nil, nil, err
	}
	defer privateKeyFile.Close()

	pemFileInfo, err := privateKeyFile.Stat()
	if err != nil {
		logger.LogErr(err, "filed to read stat from file")
	}
	size := pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pemBytes)

	if err != nil {
		logger.LogErr(err, "filed to read bytes from file")
		return nil, nil, err
	}

	data, _ := pem.Decode(pemBytes)

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		logger.LogErr(err, "filed to import key")
		return nil, nil, err
	}

	publicKeyFile, err := os.Open(public)
	if err != nil {
		logger.LogErr(err, "filed to open file")
		return nil, nil, err
	}
	defer publicKeyFile.Close()

	pemFileInfo, err = publicKeyFile.Stat()
	if err != nil {
		logger.LogErr(err, "filed to read stat from file")
	}
	size = pemFileInfo.Size()
	pemBytes = make([]byte, size)
	buffer = bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pemBytes)

	if err != nil {
		logger.LogErr(err, "filed to read bytes from file")
		return nil, nil, err
	}

	data, _ = pem.Decode(pemBytes)

	publicKeyImported, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		logger.LogErr(err, "filed to import key")
		return nil, nil, err
	}

	return privateKeyImported, publicKeyImported, nil
}

//func LoadPublicPEMKey(filename string, logger *loggers.Logger) (*rsa.PublicKey, error) {
//	publicKeyFile, err := os.Open(filename)
//	//"../../internal/crypto/" +
//	if err != nil {
//		logger.LogErr(err, "filed to open file")
//		return nil, err
//	}
//	defer publicKeyFile.Close()
//
//	pemFileInfo, err := publicKeyFile.Stat()
//	if err != nil {
//		logger.LogErr(err, "filed to read stat from file")
//	}
//	size := pemFileInfo.Size()
//	pemBytes := make([]byte, size)
//	buffer := bufio.NewReader(publicKeyFile)
//	_, err = buffer.Read(pemBytes)
//
//	if err != nil {
//		logger.LogErr(err, "filed to read bytes from file")
//		return nil, err
//	}
//
//	data, _ := pem.Decode(pemBytes)
//
//	public, err := x509.ParsePKCS1PublicKey(data.Bytes)
//	if err != nil {
//		logger.LogErr(err, "filed to import key")
//		return nil, err
//	}
//	return public, nil
//}

func (r *RSA) EncryptedCard(d *model.CryptoCard) {
	var err error
	d.Number, err = r.EncryptedData(d.Number, r.Public)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
	d.Name, err = r.EncryptedData(d.Name, r.Public)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
	d.CVC, err = r.EncryptedData(d.CVC, r.Public)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
}
func (r *RSA) EncryptedPassword(d *model.CryptoPassword) {
	var err error
	d.Data, err = r.EncryptedData(d.Data, r.Public)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
}
func (r *RSA) EncryptedTextData(d *model.CryptoTextData) {
	var err error
	d.Text, err = r.EncryptedData(d.Text, r.Public)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
}
func (r *RSA) EncryptedBinaryData(d *model.CryptoBinaryData) {
	var err error
	d.Data, err = r.EncryptedData(d.Data, r.Public)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
}

func (r *RSA) DecryptedCard(d *model.CryptoCard) {
	var err error
	d.Number, err = r.DecryptedData(d.Number, r.Private)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
	d.Name, err = r.DecryptedData(d.Name, r.Private)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
	d.CVC, err = r.DecryptedData(d.CVC, r.Private)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
}

func (r *RSA) DecryptedPassword(d *model.CryptoPassword) {
	var err error
	d.Data, err = r.DecryptedData(d.Data, r.Private)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
}

func (r *RSA) DecryptedTextData(d *model.CryptoTextData) {
	var err error
	d.Text, err = r.DecryptedData(d.Text, r.Private)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
}

func (r *RSA) DecryptedBinaryData(d *model.CryptoBinaryData) {
	var err error
	d.Data, err = r.DecryptedData(d.Data, r.Private)
	if err != nil {
		r.logger.LogErr(err, "")
		os.Exit(1)
	}
}
