package main

import (
	"flag"
	"log"
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

	// List accounts
	accounts, err := gslt.GetAccontList(*APIKey)
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
				if err := gslt.DeleteAccount(*APIKey, server.SteamID); err != nil {
					log.Println("Failed to delete GSLT:", err)
					return
				}
			}(server)
		}
	}
	wg.Wait()
}
