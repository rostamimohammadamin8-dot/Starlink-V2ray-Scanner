package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Deep Search Engine Started
	fmt.Println("Starting MegaCode Deep Search Engine...")
	
	// List of high-speed Starlink-standard sources
	sources := []string{
		"https://raw.githubusercontent.com/yebekhe/TV2Ray/main/configs/configs",
		"https://raw.githubusercontent.com/mahdibland/V2RayAggregator/master/sub/sub_merge.txt",
	}

	for _, url := range sources {
		start := time.Now()
		
		// Professional Timeout Configuration
		client := http.Client{
			Timeout: 2 * time.Second,
		}
		
		resp, err := client.Get(url)
		
		if err == nil && resp.StatusCode == 200 {
			latency := time.Since(start).Milliseconds()
			// Only filter for low latency (Starlink Standard)
			if latency < 150 {
				fmt.Printf("STATUS: [HEALTHY] | LATENCY: %dms | SOURCE: %s\n", latency, url)
			}
		} else {
			fmt.Printf("STATUS: [FAILED] | SOURCE: %s\n", url)
		}
	}
	
	fmt.Println("Search Operation Completed.")
}
