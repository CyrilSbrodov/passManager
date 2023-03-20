package manager

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/CyrilSbrodov/passManager.git/client/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/client/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/client/crypto"
	"github.com/CyrilSbrodov/passManager.git/client/model"
)

type Manager struct {
	client              http.Client
	config              *config.Config
	crypto              crypto.Crypto
	logger              *loggers.Logger
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	publicKeyFromServer *rsa.PublicKey
	url                 string
	jwt                 string
}

type Managers interface {
	Register(login, password string) error
	Auth(login, password string) error
	AddCard(d *model.CryptoCard) error
	AddPassword(d *model.CryptoPassword) error
	AddText(d *model.CryptoTextData) error
	AddBinary(d *model.CryptoBinaryData) error
	GetCards() (string, error)
	GetPasswords() (string, error)
	GetText() (string, error)
	GetBinary() (string, error)
	DeleteCard(id int) error
	DeleteText(id int) error
	DeletePassword(id int) error
	DeleteBinary(id int) error
	UpdateCard(d *model.CryptoCard) error
	UpdatePassword(d *model.CryptoPassword) error
	UpdateText(d *model.CryptoTextData) error
	UpdateBinary(d *model.CryptoBinaryData) error
}

func NewManager(logger *loggers.Logger, cfg *config.Config, client http.Client) *Manager {
	c := crypto.NewRSA(*cfg, logger)
	return &Manager{
		client:              client,
		config:              cfg,
		logger:              logger,
		privateKey:          c.Private,
		publicKey:           c.Public,
		publicKeyFromServer: nil,
		url:                 "http://",
		jwt:                 "",
		crypto:              c,
	}
}

func (m *Manager) Register(login, password string) error {
	var u model.User
	u.Login = login
	u.Password = password
	uByte, err := json.Marshal(u)
	if err != nil {
		m.logger.LogErr(err, "failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/register", bytes.NewBuffer(uByte))
	if err != nil {
		m.logger.LogErr(err, "failed to request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusConflict:
		fmt.Printf("login is alrady registered")
		return fmt.Errorf("login is alrady registered")
	case http.StatusBadRequest:
		fmt.Printf("login or password is empty")
		return fmt.Errorf("login or password is empty")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "failed to read body")
		return err
	}

	var accept model.KeyAndToken
	if err = json.Unmarshal(data, &accept); err != nil {
		m.logger.LogErr(err, "failed to unmarshal publicKey")
		return err
	}

	m.publicKeyFromServer = accept.Key
	m.jwt = accept.Token

	resp.Body.Close()

	return nil
}

func (m *Manager) Auth(login, password string) error {
	var u model.User
	u.Login = login
	u.Password = password
	uByte, err := json.Marshal(u)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}
	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/login", bytes.NewBuffer(uByte))
	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}
	switch resp.StatusCode {
	case http.StatusBadRequest:
		fmt.Printf("login or password is empty")
		return fmt.Errorf("login or password is empty")
	case http.StatusUnauthorized:
		fmt.Printf("wrong login or password")
		return fmt.Errorf("wrong login or password")
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "failed to read body")
		return err
	}

	var accept model.KeyAndToken
	if err = json.Unmarshal(data, &accept); err != nil {
		m.logger.LogErr(err, "failed to unmarshal publicKey")
		return err
	}

	m.publicKeyFromServer = accept.Key
	m.jwt = accept.Token

	resp.Body.Close()

	return nil
}

func (m *Manager) AddCard(data *model.CryptoCard) error {
	m.crypto.EncryptedCard(data)
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/cards", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) AddPassword(data *model.CryptoPassword) error {
	m.crypto.EncryptedPassword(data)
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/password", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) AddText(data *model.CryptoTextData) error {
	m.crypto.EncryptedTextData(data)
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/text", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) AddBinary(data *model.CryptoBinaryData) error {
	m.crypto.EncryptedBinaryData(data)
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/binary", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) GetCards() (string, error) {
	req, err := http.NewRequest(http.MethodGet, m.url+m.config.Addr+"/api/data/cards", nil)

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return "", err
	}
	switch resp.StatusCode {
	case http.StatusNoContent:
		fmt.Printf("No cards")
		return "", fmt.Errorf("no cards")
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return "", fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return "", fmt.Errorf("unauthorized")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return "", err
	}
	defer resp.Body.Close()
	var d []model.CryptoCard
	var cards []model.CryptoCard
	if err = json.Unmarshal(data, &d); err != nil {
		m.logger.LogErr(err, "Failed to unmarshal body")
		return "", err
	}

	if err != nil {
		m.logger.LogErr(err, "Failed to decrypted data")
		return "", err
	}
	for i := 0; i < len(d); i++ {
		var card model.CryptoCard
		m.crypto.DecryptedCard(&d[i])
		card = d[i]
		cards = append(cards, card)
	}
	result := "\nyou have these cards:\n"
	for _, card := range cards {
		result += fmt.Sprintf("%v. Number: %v, Name: %s, CVC: %v \n", card.UID, string(card.Number),
			string(card.Name), string(card.CVC))
	}
	return result, nil
}

func (m *Manager) GetPasswords() (string, error) {
	req, err := http.NewRequest(http.MethodGet, m.url+m.config.Addr+"/api/data/password", nil)

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		fmt.Printf("No cards")
		return "", fmt.Errorf("no cards")
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return "", fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return "", fmt.Errorf("unauthorized")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return "", err
	}
	defer resp.Body.Close()
	var d []model.CryptoPassword
	var passwords []model.CryptoPassword
	if err = json.Unmarshal(data, &d); err != nil {
		m.logger.LogErr(err, "Failed to unmarshal body")
		return "", err
	}

	if err != nil {
		m.logger.LogErr(err, "Failed to decrypted data")
		return "", err
	}
	for i := 0; i < len(d); i++ {
		var p model.CryptoPassword
		m.crypto.DecryptedPassword(&d[i])
		p = d[i]
		passwords = append(passwords, p)
	}
	result := "\nyou have these passwords:\n"
	for _, pass := range passwords {
		result += fmt.Sprintf("%v. Password: %v \n", pass.UID, string(pass.Data))
	}
	return result, nil
}

func (m *Manager) GetText() (string, error) {
	req, err := http.NewRequest(http.MethodGet, m.url+m.config.Addr+"/api/data/text", nil)

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		fmt.Printf("No cards")
		return "", fmt.Errorf("no cards")
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return "", fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return "", fmt.Errorf("unauthorized")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return "", err
	}
	defer resp.Body.Close()
	var d []model.CryptoTextData
	var texts []model.CryptoTextData
	if err = json.Unmarshal(data, &d); err != nil {
		m.logger.LogErr(err, "Failed to unmarshal body")
		return "", err
	}

	if err != nil {
		m.logger.LogErr(err, "Failed to decrypted data")
		return "", err
	}
	for i := 0; i < len(d); i++ {
		var t model.CryptoTextData
		m.crypto.DecryptedTextData(&d[i])
		t = d[i]
		texts = append(texts, t)
	}
	result := "\nyou have these text information:\n"
	for _, text := range texts {
		result += fmt.Sprintf("%v. Text: %v \n", text.UID, string(text.Text))
	}
	return result, nil
}

func (m *Manager) GetBinary() (string, error) {
	req, err := http.NewRequest(http.MethodGet, m.url+m.config.Addr+"/api/data/binary", nil)

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		fmt.Printf("No cards")
		return "", fmt.Errorf("no cards")
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return "", fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return "", fmt.Errorf("unauthorized")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return "", err
	}
	defer resp.Body.Close()
	var d []model.CryptoBinaryData
	var binarys []model.CryptoBinaryData
	if err = json.Unmarshal(data, &d); err != nil {
		m.logger.LogErr(err, "Failed to unmarshal body")
		return "", err
	}

	if err != nil {
		m.logger.LogErr(err, "Failed to decrypted data")
		return "", err
	}
	for i := 0; i < len(d); i++ {
		var b model.CryptoBinaryData
		m.crypto.DecryptedBinaryData(&d[i])
		b = d[i]
		binarys = append(binarys, b)
	}
	result := "\nyou have these binary data:\n"
	for _, binary := range binarys {
		result += fmt.Sprintf("%v. Binary: %v \n", binary.UID, string(binary.Data))
	}
	return result, nil
}

func (m *Manager) DeleteCard(id int) error {
	var data model.CryptoCard
	data.UID = id
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/delete/cards", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) DeleteText(id int) error {
	var data model.CryptoTextData
	data.UID = id

	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/delete/text", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) DeletePassword(id int) error {
	var data model.CryptoPassword
	data.UID = id

	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/delete/password", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) DeleteBinary(id int) error {
	var data model.CryptoBinaryData
	data.UID = id

	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/delete/binary", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) UpdateCard(data *model.CryptoCard) error {
	m.crypto.EncryptedCard(data)
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/update/cards", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}
func (m *Manager) UpdatePassword(data *model.CryptoPassword) error {
	m.crypto.EncryptedPassword(data)
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/update/password", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) UpdateText(data *model.CryptoTextData) error {
	m.crypto.EncryptedTextData(data)
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/update/text", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (m *Manager) UpdateBinary(data *model.CryptoBinaryData) error {
	m.crypto.EncryptedBinaryData(data)
	uByte, err := json.Marshal(data)
	if err != nil {
		m.logger.LogErr(err, "Failed to marshal")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url+m.config.Addr+"/api/data/update/binary", bytes.NewBuffer(uByte))

	if err != nil {
		m.logger.LogErr(err, "Failed to request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.LogErr(err, "Failed to do request")
		return err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		fmt.Printf("server error")
		return fmt.Errorf("server error")
	case http.StatusUnauthorized:
		fmt.Printf("Unauthorized")
		return fmt.Errorf("unauthorized")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		m.logger.LogErr(err, "Failed to read body")
		return err
	}
	defer resp.Body.Close()
	return nil
}
