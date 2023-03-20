package app

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

func (a *App) deleteData(reader *bufio.Reader) {
	fmt.Printf("\n\nSelect what do you want to delete?\n\n")
	var id int
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
LoopDelete:
	for {
		switch dataSelect {
		case "1":
			fmt.Printf("\nPlease enter password id:\n")
			fmt.Fscan(reader, &id)
			if err := a.manager.DeletePassword(id); err != nil {
				fmt.Printf("\nsomething wrong, try again")
				break LoopDelete
			}
			break LoopDelete
		case "2":
			fmt.Printf("\nPlease enter card id:\n")
			fmt.Fscan(reader, &id)
			if err = a.manager.DeleteCard(id); err != nil {
				fmt.Printf("\nsomething wrong, try again")
				break LoopDelete
			}
			break LoopDelete
		case "3":
			fmt.Printf("\nPlease enter text id:\n")
			fmt.Fscan(reader, &id)
			if err := a.manager.DeleteText(id); err != nil {
				fmt.Printf("\nsomething wrong, try again")
				break LoopDelete
			}
			break LoopDelete
		case "4":
			fmt.Printf("\nPlease enter binary id:\n")
			fmt.Fscan(reader, &id)
			if err := a.manager.DeleteBinary(id); err != nil {
				fmt.Printf("\nsomething wrong, try again")
				break LoopDelete
			}
			break LoopDelete
		case "5":
			break LoopDelete
		}
	}
}
