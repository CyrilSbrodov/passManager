package app

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/CyrilSbrodov/passManager.git/client/model"
)

func (a *App) updateData(reader *bufio.Reader) {
	fmt.Printf("\n\nSelect what do you want to update?\n\n")

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
LoopUpdate:
	for {
		switch dataSelect {
		case "1":
			fmt.Printf("\nPlease enter password's id:\n")
			fmt.Fscan(reader, &p.UID)
			fmt.Printf("\nPlease enter login:\n")
			fmt.Fscan(reader, &p.Login)
			fmt.Printf("\nPlease enter password:\n")
			fmt.Fscan(reader, &p.Pass)
			if err := a.manager.UpdatePassword(&p); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break LoopUpdate
		case "2":
			fmt.Printf("\nPlease enter new card information:\n")
			fmt.Printf("\nCard's id:\n")
			fmt.Fscan(reader, &c.UID)
			fmt.Printf("\nCard number:\n")
			fmt.Fscan(reader, &c.Number)
			fmt.Printf("\nCard holder:\n")
			fmt.Fscan(reader, &c.Name)
			fmt.Printf("\nCVC number:\n")
			fmt.Fscan(reader, &c.CVC)
			if err := a.manager.UpdateCard(&c); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break LoopUpdate
		case "3":
			fmt.Printf("\nPlease enter a text's id:\n")
			fmt.Fscan(reader, &t.UID)
			fmt.Printf("\nPlease enter a text:\n")
			fmt.Fscan(reader, &t.Text)
			if err := a.manager.UpdateText(&t); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break LoopUpdate
		case "4":
			fmt.Printf("\nPlease enter a bibary's id:\n")
			fmt.Fscan(reader, &b.UID)
			fmt.Printf("\nPlease enter a binary:\n")
			fmt.Fscan(reader, &b.Data)
			if err := a.manager.UpdateBinary(&b); err != nil {
				fmt.Printf("\nsomething wrong, try again")
			}
			break LoopUpdate
		case "5":
			break LoopUpdate
		}
	}
}
