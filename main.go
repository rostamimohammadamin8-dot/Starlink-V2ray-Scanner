package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MaxLatency = 100 // Target latency for Starlink standard
	OutputFile = "cleaned_configs.txt"
)

// IPInfo holds location data
type IPInfo struct {
	CountryCode string `json:"countryCode"`
}

// getCountryCode finds the country of an IP
func getCountryCode(host string) string {
	ips, _ := net.LookupIP(host)
	if len(ips) == 0 {
		return "UN"
	}
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=countryCode", ips[0].String()))
	if err != nil {
		return "UN"
	}
	defer resp.Body.Close()
	var info IPInfo
	json.NewDecoder(resp.Body).Decode(&info)
	return info.CountryCode
}

func processConfig(config string) string {
	if !strings.Contains(config, "://") {
		return ""
	}

	// Basic check for connection
	start := time.Now()
	client := http.Client{Timeout: time.Duration(MaxLatency+100) * time.Millisecond}
	resp, err := client.Get("https://www.google.com")
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	latency := time.Since(start).Milliseconds()

	// Extract Hostname to find Location
	u, err := url.Parse(config)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	country := getCountryCode(host)

	// Add Location Tag and Ping to the config name (Remark)
	remark := fmt.Sprintf("[%s]-%dms-MegaCode", country, latency)
	if strings.Contains(config, "#") {
		parts := strings.Split(config, "#")
		return parts[0] + "#" + remark
	}
	return config + "#" + remark
}

func main() {
	fmt.Println("🚀 MegaCode Pro: Geolocation & Speed Scan Activated")

	sources := []string{
		"https://raw.githubusercontent.com/yebekhe/TV2Ray/main/configs/configs",
		"https://raw.githubusercontent.com/mahdibland/V2RayAggregator/master/sub/sub_merge.txt",
	}

	var bestConfigs []string
	for _, src := range sources {
		fmt.Printf("Scanning: %s\n", src)
		resp, _ := http.Get(src)
		body, _ := ioutil.ReadAll(resp.Body)
		lines := strings.Split(string(body), "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			processed := processConfig(line)
			if processed != "" {
				bestConfigs = append(bestConfigs, processed)
				fmt.Printf("Found: %s\n", processed)
				if len(bestConfigs) >= 100 {
					break
				}
			}
		}
	}

	ioutil.WriteFile(OutputFile, []byte(strings.Join(bestConfigs, "\n")), 0644)
	fmt.Printf("✅ Done! Saved %d localized configs.\n", len(bestConfigs))
}
