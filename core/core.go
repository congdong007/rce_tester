package core

import (
	"fmt"
	"strings"
	"sync"

	cfgpkg "rce_tester/config"
	"rce_tester/util"
	ut "rce_tester/util"
	wk "rce_tester/worker"
)

// Run orchestrates loading data, launching workers, and saving results
func Run(cfg *cfgpkg.Config) {
	payloads, err := ut.LoadLines(cfg.Payloads)
	if err != nil {
		panic(err)
	}

	var keywords []string
	if cfg.Keywords != "" {
		keywords, err = ut.LoadLines(cfg.Keywords)
		if err != nil {
			panic(err)
		}
	}

	var host, hosturl string

	perCfg := *cfg
	if cfg.URL != "" {
		hosturl, host, _, err = ut.ResolveHostFromURL(cfg.URL)
		if err != nil {
			panic(err)
		}
		perCfg.Host = host
		perCfg.HostUrl = hosturl
	}

	// prepare targets
	var targets []string
	if cfg.Uf != "" {
		targets, err = ut.LoadLines(cfg.Uf)
		if err != nil {
			panic(err)
		}
		if len(targets) == 0 {
			panic("url file empty")
		}
	} else {
		targets = []string{cfg.URL}
	}

	fmt.Printf("Testing %d targets with %d payloads, threads=%d per target\n", len(targets), len(payloads), cfg.Threads)

	var wg sync.WaitGroup
	sem := make(chan struct{}, cfg.MaxUf)

	for idx, rawURL := range targets {
		wg.Add(1)
		sem <- struct{}{} // get sem
		go func(idx int, rawURL string) {
			defer wg.Done()
			defer func() { <-sem }() // release sem

			perCfg.URL = rawURL

			parts := strings.SplitN(rawURL, "?", 2)
			urlBase := parts[0]
			dataPart := ""
			if len(parts) == 2 {
				dataPart = parts[1]
			} else {
				dataPart = cfg.Data
			}

			fmt.Printf("=== [%d/%d] Testing %s === %s\n", idx+1, len(targets), rawURL, dataPart)

			perCfg.URL = urlBase

			resultsCh, err := wk.RunWorkersWithData(&perCfg, payloads, keywords, dataPart)
			if err != nil {
				fmt.Printf("Worker error for %s: %v\n", rawURL, err)
				return
			}

			if err := util.WriteResults(perCfg.Output, resultsCh); err != nil {
				fmt.Printf("Write results failed for %s: %v\n", rawURL, err)
			} else {
				fmt.Printf("Results appended to %s\n", perCfg.Output)
			}

		}(idx, rawURL)
	}

	wg.Wait()
	fmt.Println("All targets processed.")

}
