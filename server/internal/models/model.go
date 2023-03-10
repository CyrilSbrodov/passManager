package models

type User struct {
	ID       int
	Login    string
	Password string
}

type Card struct {
	Name   string
	Number []int
	CVC    int
}

type TextData string

type BinaryData string

type Data struct {
	Card       Card
	TextData   TextData
	BinaryData BinaryData
}
