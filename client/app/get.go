package app

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

func (a *App) getData(reader *bufio.Reader) {
	fmt.Printf("\n\nSelect what do you want to load?\n\n")

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
LoopSecond:
	for {
		switch dataSelect {
		case "1":
			d, err := a.manager.GetPasswords()
			if err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			fmt.Printf(d + "\n")
			break LoopSecond
		case "2":
			d, err := a.manager.GetCards()
			if err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			fmt.Printf(d + "\n")
			break LoopSecond
		case "3":
			d, err := a.manager.GetText()
			if err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			fmt.Printf(d + "\n")
			break LoopSecond
		case "4":
			d, err := a.manager.GetBinary()
			if err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			fmt.Printf(d + "\n")
			break LoopSecond
		case "5":
			break LoopSecond
		}
	}
}
