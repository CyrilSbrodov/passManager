package app

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/CyrilSbrodov/passManager.git/client/model"
)

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
		log.Fatalf("Error in reading input: %v", err)
	}
	dataSelect = strings.TrimSpace(dataSelect)
Loop:
	for {
		switch dataSelect {
		case "1":
			fmt.Printf("\nPlease enter a password:\n")
			fmt.Fscan(reader, &p.Data)
			if err := a.manager.AddPassword(&p); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break Loop
		case "2":
			fmt.Printf("\nPlease enter card information:\n")
			fmt.Printf("\nCard number:\n")
			fmt.Fscan(reader, &c.Number)
			fmt.Printf("\nCard holder:\n")
			fmt.Fscan(reader, &c.Name)
			fmt.Printf("\nCVC number:\n")
			fmt.Fscan(reader, &c.CVC)
			if err := a.manager.AddCard(&c); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break Loop
		case "3":
			fmt.Printf("\nPlease enter a text:\n")
			fmt.Fscan(reader, &t.Text)
			if err := a.manager.AddText(&t); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break Loop
		case "4":
			fmt.Printf("\nPlease enter a binary:\n")
			fmt.Fscan(reader, &b.Data)
			if err := a.manager.AddBinary(&b); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break Loop
		case "5":
			break Loop
		}
	}
}
