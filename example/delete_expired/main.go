package main

import (
	"flag"
	"log"
	"strconv"
	"sync"

	"github.com/FlowingSPDG/gslt-go"
)

func main() {
	var (
		APIKey = flag.String("apikey", "", "Steam Web API Key.")
	)
	flag.Parse()
	if *APIKey == "" {
		log.Fatalln("API Key required!")
		return
	}

	service := gslt.NewGameServerService(*APIKey)

	// List accounts
	accounts, err := service.GetAccountList()
	if err != nil {
		panic(err)
	}

	// Delete with goroutine
	wg := sync.WaitGroup{}
	for _, server := range accounts.Response.Servers {
		if server.IsExpired {
			wg.Add(1)
			go func(server gslt.Server) {
				defer wg.Done()
				steamid, err := strconv.ParseUint(server.SteamID, 10, 64)
				if err != nil {
					log.Println("Failed to parse steamid:", err)
					return
				}
				if err := service.DeleteAccount(steamid); err != nil {
					log.Println("Failed to delete GSLT:", err)
					return
				}
			}(server)
		}
	}
	wg.Wait()
}
