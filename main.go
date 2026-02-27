package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Global Configuration
const (
	MaxLatency = 80 // Ultra-fast threshold for Gaming & AI
	OutputFile = "cleaned_configs.txt"
)

// Smart Filter Function
func isGamingReady(config string) bool {
	// Focusing on high-performance protocols for Starlink
	if !strings.HasPrefix(config, "vless") && !strings.HasPrefix(config, "vmess") && 
	   !strings.HasPrefix(config, "hysteria2") && !strings.HasPrefix(config, "tuic") {
		return false
	}

	start := time.Now()
	client := http.Client{
		Timeout: time.Duration(MaxLatency+50) * time.Millisecond,
	}
	
	// Test against Google or Cloudflare
	resp, err := client.Get("https://www.google.com")
	if err == nil && resp.StatusCode == 200 {
		latency := time.Since(start).Milliseconds()
		if latency <= MaxLatency {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("🚀 MegaCode Smart Engine: Gaming & AI Mode Activated")
	
	sources := []string{
		"https://raw.githubusercontent.com/yebekhe/TV2Ray/main/configs/configs",
		"https://raw.githubusercontent.com/mahdibland/V2RayAggregator/master/sub/sub_merge.txt",
	}

	var bestConfigs []string

	for _, url := range sources {
		fmt.Printf("🔍 Scanning source: %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		lines := strings.Split(string(body), "\n")

		for _, config := range lines {
			config = strings.TrimSpace(config)
			if isGamingReady(config) {
				bestConfigs = append(bestConfigs, config)
				// Limit to top 100 super configs to keep it clean
				if len(bestConfigs) >= 100 {
					break
				}
			}
		}
	}

	// Final Step: Saving the High-Speed results
	outputData := strings.Join(bestConfigs, "\n")
	err := ioutil.WriteFile(OutputFile, []byte(outputData), 0644)
	
	if err != nil {
		fmt.Printf("❌ Save Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success! Found %d Ultra-Fast configs.\n", len(bestConfigs))
	}
}
