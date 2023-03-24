// Package app пакет для вызова бесконечного цикла с выбором возможных действий с сервером.
package app

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/CyrilSbrodov/passManager.git/client/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/client/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/client/crypto"
	"github.com/CyrilSbrodov/passManager.git/client/manager"
)

type App struct {
	crypto  crypto.Crypto
	manager manager.Managers
	client  http.Client
	logger  *loggers.Logger
	cfg     *config.Config
}

func NewApp() *App {
	cfg := config.ConfigInit()
	logger := loggers.NewLogger()
	c := crypto.NewRSA(*cfg, logger)
	client := &http.Client{}
	m := manager.NewManager(logger, cfg, *client)
	return &App{
		crypto:  c,
		manager: m,
		client:  *client,
		logger:  logger,
		cfg:     cfg,
	}
}

// Run - функция запуска клиента.
func (a *App) Run() {
	for {
		fmt.Printf("\n\nWhat do you want to do?\n\n")
		options :=
			`1. Register on PASSMANAGER.
2. Auth on PASSMANAGER.
3. Add data.
4. Update data.
5. Remove data.
6. View data on server.
7. Exit program.
		
`
		fmt.Printf(options)
		reader := bufio.NewReader(os.Stdin)

		option, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error in reading input: %v", err)
		}
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			a.register(reader)
		case "2":
			a.login(reader)
		case "3":
			a.addData(reader)
		case "4":
			a.updateData(reader)
		case "5":
			a.deleteData(reader)
		case "6":
			a.getData(reader)
		default:
			fmt.Println("Please enter a valid option in the given list!")
		case "7":
			break
		}
		if option == "7" {
			fmt.Println("Exiting PASSMANAGER.")
			break
		}
	}
}
