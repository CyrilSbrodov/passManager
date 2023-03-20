// Package app пакет для вызова бесконечного цикла с выбором возможных действий с сервером.
// Данный пакет предоставляет возможность добавлять данные на сервер.
package app

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/CyrilSbrodov/passManager.git/client/model"
)

// addData функция добавления данных.
func (a *App) addData(reader *bufio.Reader) {
	fmt.Printf("\n\nSelect what do you want to save?\n\n")

	var c model.CryptoCard
	var p model.CryptoPassword
	var t model.CryptoTextData
	var b model.CryptoBinaryData

	data :=
		`1. Password.
2. Card.
3. Text data.
4. Binary data.
5. Return.`

	fmt.Printf(data + "\n")
	dataSelect, err := reader.ReadString('\n')
	if err != nil {
		a.logger.LogErr(err, "Error in reading input")
	}
	dataSelect = strings.TrimSpace(dataSelect)

loop:
	for {
		switch dataSelect {
		case "1":
			fmt.Printf("\nPlease enter login:\n")
			p.Login, err = reader.ReadBytes('\n')
			a.checkError(err)
			fmt.Printf("\nPlease enter password:\n")
			p.Pass, err = reader.ReadBytes('\n')
			a.checkError(err)
			if err := a.manager.AddPassword(&p); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break loop
		case "2":
			fmt.Printf("\nPlease enter card information:\n")
			fmt.Printf("\nCard number:\n")
			c.Number, err = reader.ReadBytes('\n')
			a.checkError(err)
			fmt.Printf("\nCard holder:\n")
			c.Name, err = reader.ReadBytes('\n')
			a.checkError(err)
			fmt.Printf("\nCVC number:\n")
			c.CVC, err = reader.ReadBytes('\n')
			a.checkError(err)
			if err := a.manager.AddCard(&c); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break loop
		case "3":
			fmt.Printf("\nPlease enter a text:\n")
			t.Text, err = reader.ReadBytes('\n')
			a.checkError(err)
			if err := a.manager.AddText(&t); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break loop
		case "4":
			fmt.Printf("\nPlease enter a binary:\n")
			b.Data, err = reader.ReadBytes('\n')
			a.checkError(err)
			if err := a.manager.AddBinary(&b); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break loop
		case "5":
			break loop
		}
	}
}

//checkError проверка на ошибку.
func (a *App) checkError(err error) {
	if err != nil {
		a.logger.LogErr(err, "Error in reading input")
	}
}
