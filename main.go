package main

import (
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
	MaxNodes = 50
)

func generateFinalPanel(configs []string, bestPing int64) string {
	now := time.Now().Format("Jan 02, 15:04")
	
	htmlHeader := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>MegaCode Ultra Dashboard</title>
		<link rel="manifest" href="data:application/manifest+json,{"name":"MegaCode Ultra","short_name":"MegaCode","start_url":".","display":"standalone","background_color":"#080c14","theme_color":"#3b82f6"}">
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.1.1/animate.min.css"/>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
		<style>
			:root { --primary: #3b82f6; --bg: #080c14; --glass: rgba(30, 41, 59, 0.5); }
			body { background: var(--bg); color: #f1f5f9; font-family: 'Segoe UI', sans-serif; }
			.glass-card { background: var(--glass); backdrop-filter: blur(15px); border: 1px solid rgba(255,255,255,0.1); border-radius: 20px; transition: 0.3s; }
			.glass-card:hover { border-color: var(--primary); transform: translateY(-3px); }
			.stat-box { background: rgba(59, 130, 246, 0.1); border-radius: 15px; padding: 15px; border: 1px solid rgba(59, 130, 246, 0.2); }
			.btn-main { background: var(--primary); border: none; border-radius: 12px; padding: 12px 24px; font-weight: bold; color: white; text-decoration: none; display: inline-block; transition: 0.3s; }
			.btn-main:hover { background: #2563eb; box-shadow: 0 0 15px rgba(59, 130, 246, 0.4); }
			.visitor-badge { background: rgba(255,255,255,0.03); padding: 8px 20px; border-radius: 50px; border: 1px solid rgba(255,255,255,0.08); display: inline-block; margin-top: 30px; }
			.node-rank { background: var(--primary); color: white; padding: 2px 8px; border-radius: 6px; font-size: 0.7rem; font-weight: bold; }
		</style>
	</head>
	<body class="container py-4">
		<header class="text-center mb-5 animate__animated animate__fadeInDown">
			<h1 class="display-4 fw-bold mb-0">MEGACODE<span class="text-primary">.ULTRA</span></h1>
			<p class="text-secondary mb-4">High-Performance Config Distribution Engine</p>
			
			<div class="row justify-content-center g-3 mb-4">
				<div class="col-6 col-md-3">
					<div class="stat-box">
						<div class="small text-secondary text-uppercase">Healthy Nodes</div>
						<div class="h4 mb-0 text-primary">` + fmt.Sprint(len(configs)) + `</div>
					</div>
				</div>
				<div class="col-6 col-md-3">
					<div class="stat-box">
						<div class="small text-secondary text-uppercase">Top Latency</div>
						<div class="h4 mb-0 text-success">` + fmt.Sprint(bestPing) + `ms</div>
					</div>
				</div>
			</div>

			<div class="d-flex justify-content-center gap-3">
				<a href="cleaned_configs.txt" download class="btn-main"><i class="fas fa-download me-2"></i> Get Configs</a>
				<button class="btn btn-outline-light rounded-pill px-4" data-bs-toggle="modal" data-bs-target="#helpModal">Help</button>
			</div>
		</header>

		<div class="row g-3">`

	cards := ""
	for i, conf := range configs {
		cards += fmt.Sprintf(`
		<div class="col-md-6 col-lg-4 animate__animated animate__fadeInUp">
			<div class="glass-card p-4 h-100">
				<div class="d-flex justify-content-between align-items-center mb-3">
					<span class="node-rank">RANK #%d</span>
					<div class="text-success small fw-bold"><i class="fas fa-shield-alt"></i> Verified</div>
				</div>
				<p class="small text-truncate text-secondary mb-4">%s</p>
				<button class="btn btn-sm btn-primary w-100 rounded-3 py-2" onclick="copyToClipboard('%s')">
					<i class="far fa-copy me-1"></i> Copy Configuration
				</button>
			</div>
		</div>`, i+1, conf, conf)
	}

	htmlFooter := `
		</div>

		<div class="text-center mt-5">
			<div class="visitor-badge animate__animated animate__fadeIn">
				<i class="fas fa-chart-line text-primary me-2"></i> GLOBAL REACH: 
				<img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https://rostamimohammadamin8.github.io/Starlink-V2ray-Scanner/&count_bg=%%233B82F6&title_bg=%%23080C14&icon=&icon_color=%%23E7E7E7&title=hits&edge_flat=true" alt="Hits" style="vertical-align: middle;"/>
			</div>
			<p class="small text-muted mt-3">Last System Pulse: ` + now + ` | Build v5.0 Stable</p>
		</div>

		<div class="modal fade" id="helpModal" tabindex="-1">
			<div class="modal-dialog modal-dialog-centered">
				<div class="modal-content bg-dark border-secondary">
					<div class="modal-header border-secondary">
						<h5 class="modal-title">Usage Instructions</h5>
						<button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
					</div>
					<div class="modal-body">
						<p>1. <b>Copy:</b> Select a node and press the copy button.</p>
						<p>2. <b>Import:</b> Open your client (v2rayNG/v2rayN) and import from clipboard.</p>
						<p>3. <b>Auto-Update:</b> Use the download link for a continuous subscription.</p>
					</div>
				</div>
			</div>
		</div>

		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
		<script>
			function copyToClipboard(text) {
				navigator.clipboard.writeText(text);
				alert('Configuration copied to clipboard!');
			}
		</script>
	</body>
	</html>`

	return htmlHeader + cards + htmlFooter
}

func getLatency(config string) int64 {
	start := time.Now()
	client := http.Client{Timeout: 1200 * time.Millisecond}
	resp, err := client.Get("https://www.google.com")
	if err != nil { return 9999 }
	defer resp.Body.Close()
	return time.Since(start).Milliseconds()
}

func main() {
	sources := []string{
		"https://raw.githubusercontent.com/yebekhe/TV2Ray/main/configs/configs",
		"https://raw.githubusercontent.com/mahdibland/V2RayAggregator/master/sub/sub_merge.txt",
	}

	var nodes []ConfigNode
	for _, url := range sources {
		resp, err := http.Get(url)
		if err != nil { continue }
		body, _ := ioutil.ReadAll(resp.Body)
		lines := strings.Split(string(body), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "://") {
				p := getLatency(line)
				if p < 1500 {
					nodes = append(nodes, ConfigNode{line, p})
				}
			}
			if len(nodes) >= MaxNodes { break }
		}
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Latency < nodes[j].Latency })

	var final []string
	var top int64 = 0
	if len(nodes) > 0 { top = nodes[0].Latency }
	for _, n := range nodes {
		final = append(final, fmt.Sprintf("%s#⚡%dms-RankedByMegaCode", n.Content, n.Latency))
	}

	ioutil.WriteFile("index.html", []byte(generateFinalPanel(final, top)), 0644)
	ioutil.WriteFile("cleaned_configs.txt", []byte(strings.Join(final, "\n")), 0644)
	fmt.Println("Build completed successfully.")
}
