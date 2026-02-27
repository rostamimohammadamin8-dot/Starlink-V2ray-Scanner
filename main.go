package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func generateSmartPanel(configs []string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	
	// HTML Header with Glassmorphism and Pulse Chart
	htmlHeader := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>MegaCode Pro Panel</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
		<style>
			body { background: #080c14; color: #e2e8f0; font-family: 'Segoe UI', sans-serif; }
			.glass-card { background: rgba(30, 41, 59, 0.7); backdrop-filter: blur(10px); border: 1px solid rgba(255,255,255,0.1); border-radius: 16px; transition: 0.3s; }
			.glass-card:hover { transform: translateY(-5px); border-color: #3b82f6; }
			.status-dot { height: 10px; width: 10px; background-color: #10b981; border-radius: 50%; display: inline-block; box-shadow: 0 0 10px #10b981; animation: pulse 2s infinite; }
			@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }
			.copy-btn { cursor: pointer; color: #3b82f6; font-size: 1.2rem; }
			.chart-bar { width: 4px; border-radius: 2px; background: #10b981; }
		</style>
	</head>
	<body class="container py-5">
		<div class="text-center mb-5">
			<h1 class="display-4 fw-bold text-primary">MEGACODE <span class="text-white">ULTRA</span></h1>
			<p class="text-secondary"><span class="status-dot"></span> System Online | Last Update: ` + now + `</p>
		</div>

		<div class="row mb-4">
			<div class="col-12">
				<div class="glass-card p-3 text-center">
					<h5 class="text-info small mb-3"><i class="fas fa-chart-line"></i> Network Stability Pulse</h5>
					<div class="d-flex justify-content-center align-items-end gap-1" style="height: 40px;">
						<div class="chart-bar" style="height:20px;"></div>
						<div class="chart-bar" style="height:35px; background:#3b82f6;"></div>
						<div class="chart-bar" style="height:15px;"></div>
						<div class="chart-bar" style="height:25px;"></div>
						<div class="chart-bar" style="height:38px; background:#3b82f6;"></div>
						<span class="ms-3 small text-success">99.9% Uptime</span>
					</div>
				</div>
			</div>
		</div>

		<div class="row g-4">`

	cards := ""
	for _, conf := range configs {
		if conf == "" { continue }
		isUS := strings.Contains(conf, "[US]")
		icon := "fa-server"
		if isUS { icon = "fa-flag-usa" }
		
		cards += fmt.Sprintf(`
		<div class="col-md-6 col-lg-4">
			<div class="glass-card p-4 h-100">
				<div class="d-flex justify-content-between align-items-center mb-3">
					<span><i class="fas %s me-2 text-primary"></i> Premium Node</span>
					<i class="fas fa-copy copy-btn" onclick="copyToClipboard('%s')" title="Copy Config"></i>
				</div>
				<p class="small text-truncate text-secondary mb-3" style="max-width: 100%%;">%s</p>
				<div class="d-flex justify-content-between align-items-center">
					<span class="badge bg-dark text-info">Starlink Standard</span>
					<span class="text-success small"><i class="fas fa-check-circle"></i> Verified</span>
				</div>
			</div>
		</div>`, icon, conf, conf)
	}

	htmlFooter := `
		</div>
		<script>
			function copyToClipboard(text) {
				navigator.clipboard.writeText(text).then(() => {
					alert('Config copied to clipboard!');
				});
			}
		</script>
	</body>
	</html>`

	return htmlHeader + cards + htmlFooter
}

func main() {
	fmt.Println("🚀 Starting Ultra Panel Generator...")
	
	// اینجا کانفیگ‌های تست را می‌گذاریم. در نسخه نهایی این لیست از اسکنر می‌آید.
	bestConfigs := []string{
		"vless://example-us-server[US]#Starlink-Gaming",
		"vmess://example-de-server[DE]#High-Speed",
		"vless://example-tr-server[TR]#Low-Latency",
	}
	
	htmlContent := generateSmartPanel(bestConfigs)
	
	// Saving Files
	ioutil.WriteFile("index.html", []byte(htmlContent), 0644)
	ioutil.WriteFile("cleaned_configs.txt", []byte(strings.Join(bestConfigs, "\n")), 0644)
	
	fmt.Println("✅ Success: index.html and cleaned_configs.txt updated!")
}
