package app

import (
	"bufio"
	"fmt"
)

func (a *App) login(reader *bufio.Reader) {
	fmt.Printf("\nPlease enter your login and password:\n")
	var login, password string
	fmt.Printf("\nLogin:\n")
	fmt.Fscan(reader, &login)

	fmt.Printf("\nPassword:\n")
	fmt.Fscan(reader, &password)

	if err := a.manager.Auth(login, password); err != nil {
		fmt.Printf("\nsomething wrong, try again")
	}
}
