package main

import (
	"flag"
	"fmt"
	"github.com/FlowingSPDG/gslt-go"
	"os"
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

	// Initialize Manager{} instance...
	Manager := gslt.Manager{}
	Manager.APIToken = *APIKey
	Manager.GetList()
	gslt1, err := Manager.Generate(*Memo, uint32(*APPID))
	if err != nil {
		fmt.Printf("Failed to generate GSLT...\nERR : %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Generated GSLT : %v\n", gslt1)

	// Or you can just generate it...
	/*
		gslt2, err := gslt.GetGSLT(*APIKey, *Memo, *APPID)
		if err != nil {
			fmt.Printf("Failed to generate GSLT...\nERR : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated GSLT2 : %v\n", gslt2.LoginToken)
		err = gslt2.Delete()
		if err != nil {
			fmt.Printf("Failed to Delete GSLT %s\nERR : %v\n", gslt2.LoginToken, err)
			os.Exit(1)
		}
		fmt.Printf("Deleted GSLT2 : %v\n", gslt2.LoginToken)
	*/
}
