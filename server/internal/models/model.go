// Package models package model позволяет использовать все структуры. Между сервером и клиентом,
// т.к. на клиенте такие же структуры данных.
package models

import "crypto/rsa"

// User - структура пользователя.
type User struct {
	UID      string `json:"uid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Card - структура карты.
type Card struct {
	UID    int    `json:"UID"`
	Name   string `json:"name"`
	Number string `json:"number"`
	CVC    int    `json:"cvc"`
}

// TextData - структура текстовых данных.
type TextData struct {
	UID  int    `json:"uid_text"`
	Text string `json:"text"`
}

// BinaryData - структура бинарных данных.
type BinaryData struct {
	UID  int    `json:"uid_binary"`
	Data string `json:"data"`
}

// Password - структура пары логин/пароль
type Password struct {
	UID   int    `json:"uid_pass"`
	Login string `json:"data_login"`
	Pass  string `json:"pass"`
}

// Data - общая структура всех данных.
type Data struct {
	Password   []Password   `json:"data_password"`
	Card       []Card       `json:"data_card"`
	TextData   []TextData   `json:"data_text"`
	BinaryData []BinaryData `json:"data_binary"`
}

// KeyAndToken - структура получения ключа шифрования и токена авторизации от сервера.
type KeyAndToken struct {
	Key   *rsa.PublicKey `json:"key"`
	Token string         `json:"token"`
}

// CryptoPassword - структура зашифрованной пары логина/пароля
type CryptoPassword struct {
	UID   int    `json:"uid_pass"`
	Login []byte `json:"data_pass"`
	Pass  []byte `json:"pass"`
}

// CryptoBinaryData - структура зашифрованных бинарных данных.
type CryptoBinaryData struct {
	UID  int    `json:"uid_binary"`
	Data []byte `json:"data"`
}

// CryptoTextData - структура зашифрованных текстовых данных.
type CryptoTextData struct {
	UID  int    `json:"uid_text"`
	Text []byte `json:"text"`
}

// CryptoCard - структура зашифрованных карт.
type CryptoCard struct {
	UID    int    `json:"UID"`
	Name   []byte `json:"name"`
	Number []byte `json:"number"`
	CVC    []byte `json:"cvc"`
}

// CryptoData - общая структура всех зашифрованных данных.
type CryptoData struct {
	Password   []CryptoPassword   `json:"data_password"`
	Card       []CryptoCard       `json:"data_card"`
	TextData   []CryptoTextData   `json:"data_text"`
	BinaryData []CryptoBinaryData `json:"data_binary"`
}
