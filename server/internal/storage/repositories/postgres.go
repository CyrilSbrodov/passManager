// Package repositories позволяет сохранять и обрабатывать данные в базе данных. Так же отдавать их клиенту по запросу.
package repositories

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/server/internal/models"
	"github.com/CyrilSbrodov/passManager.git/server/pkg/client/postgres"
)

// Store - структура репозитория.
type Store struct {
	client postgres.Client
	Hash   string
	logger loggers.Logger
}

// createTable - функция создания новых таблиц в БД.
func createTable(ctx context.Context, client postgres.Client, logger *loggers.Logger) error {
	tx, err := client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.LogErr(err, "failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)

	//создание таблиц
	q := `CREATE TABLE if not exists users (
    		id BIGINT PRIMARY KEY generated always as identity,
    		login VARCHAR(200) NOT NULL unique,
    		hashed_password VARCHAR(200) NOT NULL
		);
		CREATE UNIQUE INDEX if not exists users_login_uindex on users (login);
		CREATE TABLE if not exists text_table (
    		user_id BIGINT,
    		id BIGINT PRIMARY KEY generated always as identity,
    		FOREIGN KEY (user_id) REFERENCES users(id),
    		text bytea                         
		);
		CREATE TABLE if not exists binary_table (
    		user_id BIGINT,
    		id BIGINT PRIMARY KEY generated always as identity,
    		FOREIGN KEY (user_id) REFERENCES users(id),
    		binary_data bytea                          
		);
		CREATE TABLE if not exists passwords (
    		user_id BIGINT,
    		id BIGINT PRIMARY KEY generated always as identity,
    		FOREIGN KEY (user_id) REFERENCES users(id),
    		login bytea,
		    password bytea
		);
		CREATE TABLE if not exists cards (
    		user_id BIGINT,
    		id BIGINT PRIMARY KEY generated always as identity,
    		card_number bytea,
    		FOREIGN KEY (user_id) REFERENCES users(id),
    		card_holder bytea,
    		cvc bytea                            
		);
		CREATE UNIQUE INDEX if not exists cards_card_number_uindex on cards (card_number);`

	_, err = tx.Exec(ctx, q)
	if err != nil {
		logger.LogErr(err, "failed to create table")
		return err
	}
	return tx.Commit(ctx)
}

// NewStore - функция создания нового репозитория.
func NewStore(client postgres.Client, cfg *config.Config, logger *loggers.Logger) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := createTable(ctx, client, logger); err != nil {
		logger.LogErr(err, "failed to create table")
		return nil, err
	}
	return &Store{
		client: client,
	}, nil
}

func (s *Store) Register(u *models.User) (string, error) {
	//хэширование пароля
	hashPassword := s.hashPassword(u.Password)
	//добавление пользователя в базу
	q := `INSERT INTO users (login, hashed_password)
	   						VALUES ($1, $2) RETURNING id`
	if err := s.client.QueryRow(context.Background(), q, u.Login, hashPassword).Scan(&u.UID); err != nil {
		s.logger.LogErr(err, "Failure to insert object into table")
		return "", err
	}
	return u.UID, nil
}

func (s *Store) Login(u *models.User) (string, error) {
	var password string
	//хэширование полученного пароля
	hashPassword := s.hashPassword(u.Password)
	//получение хэш пароля, хранящегося в базе
	q := `SELECT hashed_password, id FROM users WHERE login = $1`
	if err := s.client.QueryRow(context.Background(), q, u.Login).Scan(&password, &u.UID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.LogErr(err, "Failure to select object from table")
			return "", fmt.Errorf("wrong login %s", u.Login)
		}
		fmt.Println(err)
		s.logger.LogErr(err, "Wrong login")
		return "", fmt.Errorf("wrong login %s", u.Login)
	}
	//сравнение хэш пароля полученного и хэш пароля из базы
	if hashPassword != password {
		return "", fmt.Errorf("wrong password to %s", u.Login)
	}
	return u.UID, nil
}

func (s *Store) CollectPassword(d *models.CryptoPassword, id string) (int, error) {
	q := `INSERT INTO passwords (user_id, login, password) VALUES ($1, $2, $3)`
	if _, err := s.client.Exec(context.Background(), q, id, d.Login, d.Pass); err != nil {
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	//возвращаем 200 — новые данные успешно загружены в базу.
	return 200, nil
}

func (s *Store) CollectCard(d *models.CryptoCard, id string) (int, error) {
	q := `INSERT INTO cards (user_id, card_number, card_holder, cvc) VALUES ($1, $2, $3, $4)`
	if _, err := s.client.Exec(context.Background(), q, id, d.Number, d.Name, d.CVC); err != nil {
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	//возвращаем 200 — новые данные успешно загружены в базу.
	return 200, nil
}

func (s *Store) CollectText(d *models.CryptoTextData, id string) (int, error) {
	q := `INSERT INTO text_table (user_id, text) VALUES ($1, $2)`
	if _, err := s.client.Exec(context.Background(), q, id, d.Text); err != nil {
		fmt.Println(err)
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	//возвращаем 200 — новые данные успешно загружены в базу.
	return 200, nil
}

func (s *Store) CollectBinary(d *models.CryptoBinaryData, id string) (int, error) {
	q := `INSERT INTO binary_table (user_id, binary_data) VALUES ($1, $2)`
	if _, err := s.client.Exec(context.Background(), q, id, d.Data); err != nil {
		fmt.Println(err)
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	//возвращаем 200 — новые данные успешно загружены в базу.
	return 200, nil
}

func (s *Store) GetCards(id string) (int, []models.CryptoCard, error) {
	var data []models.CryptoCard

	q := `SELECT id, card_number, card_holder, cvc	FROM cards WHERE user_id = $1`
	rows, err := s.client.Query(context.Background(), q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.LogErr(err, "Failure to select object from table")
			return 204, data, fmt.Errorf("no one cards")
		}
		s.logger.LogErr(err, "")
		return 500, data, err
	}
	//добавление всех карт в слайс
	for rows.Next() {
		var c models.CryptoCard

		err = rows.Scan(&c.UID, &c.Number, &c.Name, &c.CVC)
		if err != nil && err != pgx.ErrNoRows {
			s.logger.LogErr(err, "Failure to scan object from table")
			return 500, data, err
		}
		data = append(data, c)
	}

	return 200, data, nil
}

func (s *Store) GetPassword(id string) (int, []models.CryptoPassword, error) {
	var data []models.CryptoPassword

	q := `SELECT id, login, password FROM passwords WHERE user_id = $1`
	rows, err := s.client.Query(context.Background(), q, id)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.LogErr(err, "Failure to select object from table")
			return 204, data, fmt.Errorf("no one withdraw")
		}
		s.logger.LogErr(err, "")
		return 500, data, err
	}
	//добавление данных в слайс
	for rows.Next() {
		var p models.CryptoPassword

		err = rows.Scan(&p.UID, &p.Login, &p.Pass)
		if err != nil && err != pgx.ErrNoRows {
			s.logger.LogErr(err, "Failure to scan object from table")
			return 500, data, err
		}
		data = append(data, p)
	}

	return 200, data, nil
}

func (s *Store) GetText(id string) (int, []models.CryptoTextData, error) {
	var data []models.CryptoTextData

	q := `SELECT id, text FROM text_table WHERE user_id = $1`
	rows, err := s.client.Query(context.Background(), q, id)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.LogErr(err, "Failure to select object from table")
			return 204, data, fmt.Errorf("no one withdraw")
		}
		s.logger.LogErr(err, "")
		return 500, data, err
	}
	//добавление данных в слайс
	for rows.Next() {
		var t models.CryptoTextData

		err = rows.Scan(&t.UID, &t.Text)
		if err != nil && err != pgx.ErrNoRows {
			s.logger.LogErr(err, "Failure to scan object from table")
			return 500, data, err
		}
		data = append(data, t)
	}

	return 200, data, nil
}

func (s *Store) GetBinary(id string) (int, []models.CryptoBinaryData, error) {
	var data []models.CryptoBinaryData

	q := `SELECT id, binary_data FROM binary_table WHERE user_id = $1`
	rows, err := s.client.Query(context.Background(), q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.LogErr(err, "Failure to select object from table")
			return 204, data, fmt.Errorf("no one withdraw")
		}
		s.logger.LogErr(err, "")
		return 500, data, err
	}
	//добавление данных в слайс
	for rows.Next() {
		var b models.CryptoBinaryData

		err = rows.Scan(&b.UID, &b.Data)
		if err != nil && err != pgx.ErrNoRows {
			s.logger.LogErr(err, "Failure to scan object from table")
			return 500, data, err
		}
		data = append(data, b)
	}

	return 200, data, nil
}

func (s *Store) hashPassword(pass string) string {
	h := hmac.New(sha256.New, []byte("password"))
	h.Write([]byte(pass))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Store) DeleteCard(data *models.CryptoCard, id string) (int, error) {
	q := `DELETE FROM cards WHERE id = $1 and user_id = $2`
	if _, err := s.client.Exec(context.Background(), q, data.UID, id); err != nil {
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	return 200, nil
}

func (s *Store) DeleteText(data *models.CryptoTextData, id string) (int, error) {
	q := `DELETE FROM text_table WHERE id = $1 and user_id = $2`
	if _, err := s.client.Exec(context.Background(), q, data.UID, id); err != nil {
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	return 200, nil
}

func (s *Store) DeletePassword(data *models.CryptoPassword, id string) (int, error) {
	q := `DELETE FROM passwords WHERE id = $1 and user_id = $2`
	if _, err := s.client.Exec(context.Background(), q, data.UID, id); err != nil {
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	return 200, nil
}

func (s *Store) DeleteBinary(data *models.CryptoBinaryData, id string) (int, error) {
	q := `DELETE FROM binary_table WHERE id = $1 and user_id = $2`
	if _, err := s.client.Exec(context.Background(), q, data.UID, id); err != nil {
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	return 200, nil
}

func (s *Store) UpdateCard(data *models.CryptoCard, id string) (int, error) {
	q := `UPDATE cards SET card_number = $1, card_holder = $2, cvc = $3 WHERE id = $4 AND user_id = $5`
	if _, err := s.client.Exec(context.Background(), q, data.Number, data.Name, data.CVC, data.UID, id); err != nil {
		fmt.Println(err)
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	return 200, nil
}

func (s *Store) UpdatePassword(data *models.CryptoPassword, id string) (int, error) {
	q := `UPDATE passwords SET login = $1, password = $2 WHERE id = $3 AND user_id = $4`
	if _, err := s.client.Exec(context.Background(), q, data.Login, data.Pass, data.UID, id); err != nil {
		fmt.Println(err)
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	return 200, nil
}

func (s *Store) UpdateText(data *models.CryptoTextData, id string) (int, error) {
	q := `UPDATE text_table SET text = $1 WHERE id = $2 AND user_id = $3`
	if _, err := s.client.Exec(context.Background(), q, data.Text, data.UID, id); err != nil {
		fmt.Println(err)
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	return 200, nil
}

func (s *Store) UpdateBinary(data *models.CryptoBinaryData, id string) (int, error) {
	q := `UPDATE binary_table SET binary_data = $1 WHERE id = $2 AND user_id = $3`
	if _, err := s.client.Exec(context.Background(), q, data.Data, data.UID, id); err != nil {
		fmt.Println(err)
		s.logger.LogErr(err, "Failure to insert object into table")
		return 500, err
	}
	return 200, nil
}
