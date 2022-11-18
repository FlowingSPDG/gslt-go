package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/FlowingSPDG/gslt-go"
)

func main() {
	var (
		APIKey = flag.String("apikey", "", "Steam Web API Key.")
		APPID  = flag.Uint("appid", 730, "Game AppID to generate.")
		Memo   = flag.String("memo", "", "Memo for GSLT.")
		Num    = flag.Int("num", 1, "Number of tokens to generate")
	)
	flag.Parse()
	if *APIKey == "" {
		fmt.Println("API Key required!")
		os.Exit(1)
	}

	ts := make([]*gslt.GSLT, 0, *Num)
	mtx := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < *Num; i++ {
		wg.Add(1)
		m := fmt.Sprintf("%s-%d", *Memo, i)
		go func() {
			defer wg.Done()
			g, err := gslt.GetGSLT(*APIKey, m, uint32(*APPID))
			if err != nil {
				log.Println("Failed to get GSLT Account:", err)
				return
			}
			mtx.Lock()
			defer mtx.Unlock()
			ts = append(ts, g)
		}()
	}
	wg.Wait()
	log.Printf("%d Accounts Generated\n", len(ts))
	for i, t := range ts {
		log.Printf("Token[%d]:%#v\n", i, t)
	}
}
