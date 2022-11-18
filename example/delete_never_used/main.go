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

	// List GSLTs
	gslts, err := gslt.ListGSLT(*APIKey)
	if err != nil {
		log.Fatalln("Failed to get GSLTs:", err)
		return
	}

	// Delete with goroutine
	wg := sync.WaitGroup{}
	for _, g := range gslts {
		if g.RtLastLogon == 0 {
			wg.Add(1)
			go func(g *gslt.GSLT) {
				defer wg.Done()
				if err := g.Delete(); err != nil {
					log.Println("Failed to delete GSLT:", err)
					return
				}
			}(g)
		}
	}
	wg.Wait()
}
