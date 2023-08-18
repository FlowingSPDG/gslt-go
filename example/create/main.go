package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FlowingSPDG/gslt-go"
)

func main() {
	var (
		APIKey = flag.String("apikey", "", "Steam Web API Key.")
		APPID  = flag.Uint("appid", 730, "Game AppID to generate.")
		Memo   = flag.String("memo", "", "Memo for GSLT.")
	)
	flag.Parse()
	if *APIKey == "" {
		fmt.Println("API Key required!")
		os.Exit(1)
	}

	service := gslt.NewGameServerService(*APIKey)

	gslt1, err := service.CreateAccount(uint32(*APPID), *Memo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated : %v\n", gslt1)
}
