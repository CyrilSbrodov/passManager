// Package app пакет для вызова бесконечного цикла с выбором возможных действий с сервером.
// Данный пакет предоставляет возможность регистрации пользователя на сервере и получения доступа к остальным возможностям.
package app

import (
	"bufio"
	"fmt"
	"strings"
)

func (a *App) register(reader *bufio.Reader) {
	fmt.Printf("\nPlease enter your login and password:\n")
	fmt.Printf("\nLogin:\n")
	login, err := reader.ReadString('\n')
	a.checkError(err)
	fmt.Printf("\nPassword:\n")
	password, err := reader.ReadString('\n')
	a.checkError(err)
	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)
	if err := a.manager.Register(login, password); err != nil {
		fmt.Printf("\nsomething wrong, try again")
	}
}
