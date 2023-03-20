package models

import (
	"crypto/rsa"
)

type User struct {
	UID      string `json:"uid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Card struct {
	UID    int    `json:"UID"`
	Name   string `json:"name"`
	Number string `json:"number"`
	CVC    int    `json:"cvc"`
}

type TextData struct {
	UID  int    `json:"uid_text"`
	Text string `json:"text"`
}

type BinaryData struct {
	UID  int    `json:"uid_binary"`
	Data string `json:"data"`
}

type Password struct {
	UID  int    `json:"uid_pass"`
	Data string `json:"data_pass"`
}

type Data struct {
	Password   []Password   `json:"data_password"`
	Card       []Card       `json:"data_card"`
	TextData   []TextData   `json:"data_text"`
	BinaryData []BinaryData `json:"data_binary"`
}

type KeyAndToken struct {
	Key   *rsa.PublicKey `json:"key"`
	Token string         `json:"token"`
}

type CryptoPassword struct {
	UID   int    `json:"uid_pass"`
	Login []byte `json:"data_pass"`
	Pass  []byte `json:"pass"`
}

type CryptoBinaryData struct {
	UID  int    `json:"uid_binary"`
	Data []byte `json:"data"`
}

type CryptoTextData struct {
	UID  int    `json:"uid_text"`
	Text []byte `json:"text"`
}

type CryptoCard struct {
	UID    int    `json:"UID"`
	Name   []byte `json:"name"`
	Number []byte `json:"number"`
	CVC    []byte `json:"cvc"`
}

type CryptoData struct {
	Password   []CryptoPassword   `json:"data_password"`
	Card       []CryptoCard       `json:"data_card"`
	TextData   []CryptoTextData   `json:"data_text"`
	BinaryData []CryptoBinaryData `json:"data_binary"`
}
