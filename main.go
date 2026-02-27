package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("Starting Deep Scan: Standard Starlink & High-Speed Servers...")

	sources := []string{
		"https://raw.githubusercontent.com/yebekhe/TV2Ray/main/configs/configs",
		"https://raw.githubusercontent.com/mahdibland/V2RayAggregator/master/sub/sub_merge.txt",
	}

	var validConfigs []string

	for _, url := range sources {
		fmt.Printf("Fetching from: %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		lines := strings.Split(string(body), "\n")

		for _, config := range lines {
			if len(config) < 10 {
				continue
			}

			// Ultra-fast Test: 1-second timeout
			start := time.Now()
			client := http.Client{Timeout: 1 * time.Second}
			check, err := client.Get("https://www.google.com")

			if err == nil && check.StatusCode == 200 {
				latency := time.Since(start).Milliseconds()
				if latency < 150 { // Starlink Ping Standard
					validConfigs = append(validConfigs, config)
					if len(validConfigs) >= 50 { // Limit to 50 best configs
						break
					}
				}
			}
		}
	}

	// Save results to a file
	output := strings.Join(validConfigs, "\n")
	err := ioutil.WriteFile("cleaned_configs.txt", []byte(output), 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
	} else {
		fmt.Println("Success! cleaned_configs.txt has been updated.")
	}
}
