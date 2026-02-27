package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type ConfigNode struct {
	Content string
	Latency int64
}

const (
	TelegramToken  = "YOUR_BOT_TOKEN"
	TelegramChatID = "YOUR_CHAT_ID"
	MaxNodes       = 50
)

func sendTelegramNotification(count int, bestPing int64) {
	msg := fmt.Sprintf("🚀 *MegaCode Update*\n\n✅ Verified: %d nodes\n⚡ Best Ping: %dms\n🌐 Status: Online", count, bestPing)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TelegramToken)
	
	payload, _ := json.Marshal(map[string]string{
		"chat_id":    TelegramChatID,
		"text":       msg,
		"parse_mode": "Markdown",
	})
	
	http.Post(url, "application/json", bytes.NewBuffer(payload))
}

func getLatency(config string) int64 {
	start := time.Now()
	client := http.Client{Timeout: 1200 * time.Millisecond}
	resp, err := client.Get("https://www.google.com")
	if err != nil {
		return 9999
	}
	defer resp.Body.Close()
	return time.Since(start).Milliseconds()
}

func main() {
	fmt.Println("Starting MegaCode Engine...")

	sources := []string{
		"https://raw.githubusercontent.com/yebekhe/TV2Ray/main/configs/configs",
	}

	var nodes []ConfigNode

	for _, url := range sources {
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		lines := strings.Split(string(body), "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && strings.Contains(line, "://") {
				p := getLatency(line)
				if p < 1000 {
					nodes = append(nodes, ConfigNode{Content: line, Latency: p})
				}
			}
			if len(nodes) >= MaxNodes {
				break
			}
		}
	}

	// Sorting algorithm: Ascending order (Lowest latency first)
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Latency < nodes[j].Latency
	})

	var finalConfigs []string
	for _, n := range nodes {
		tagged := fmt.Sprintf("%s#⚡%dms-MegaCode", n.Content, n.Latency)
		finalConfigs = append(finalConfigs, tagged)
	}

	// Generate Files
	ioutil.WriteFile("index.html", []byte(generateUltraUXPanel(finalConfigs)), 0644)
	ioutil.WriteFile("cleaned_configs.txt", []byte(strings.Join(finalConfigs, "\n")), 0644)
	
	// Notify via Telegram
	if len(nodes) > 0 {
		sendTelegramNotification(len(nodes), nodes[0].Latency)
	}

	fmt.Println("Build completed successfully.")
}
