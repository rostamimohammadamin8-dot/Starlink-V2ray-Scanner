package main

import (
	"fmt"
	"html"
	"strings"
	"time"
)

type panelNode struct {
	Node rankedNode
	Link string
	Top  bool
}

type countryCount struct {
	Code  string
	Count int
}

func generateFinalPanel(cards []panelNode, countries []countryCount, healthy int, bestPing int64) string {
	now := time.Now().UTC().Format("2006-01-02 15:04 UTC")

	htmlHeader := `<!DOCTYPE html>
<html lang="en" dir="ltr">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>MegaCode Ultra Dashboard</title>
	<link rel="manifest" href='data:application/manifest+json,{"name":"MegaCode Ultra","short_name":"MegaCode","start_url":".","display":"standalone","background_color":"#080c14","theme_color":"#3b82f6"}'>
	<link id="bootstrap-css" href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.1.1/animate.min.css"/>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
	<style>
		:root { --primary: #3b82f6; --bg: #080c14; --glass: rgba(30, 41, 59, 0.5); }
		body { background: var(--bg); color: #f1f5f9; font-family: 'Inter', 'Vazirmatn', Tahoma, sans-serif; }
		.glass-card { background: var(--glass); backdrop-filter: blur(15px); border: 1px solid rgba(255,255,255,0.1); border-radius: 20px; transition: 0.3s; }
		.glass-card:hover { border-color: var(--primary); transform: translateY(-3px); }
		.stat-box { background: rgba(59, 130, 246, 0.1); border-radius: 15px; padding: 15px; border: 1px solid rgba(59, 130, 246, 0.2); }
		.btn-main { background: var(--primary); border: none; border-radius: 12px; padding: 12px 24px; font-weight: bold; color: white; text-decoration: none; display: inline-block; transition: 0.3s; }
		.btn-main:hover { background: #2563eb; box-shadow: 0 0 15px rgba(59, 130, 246, 0.4); }
		.visitor-badge { background: rgba(255,255,255,0.03); padding: 8px 20px; border-radius: 50px; border: 1px solid rgba(255,255,255,0.08); display: inline-block; margin-top: 30px; }
		.node-rank { background: var(--primary); color: white; padding: 2px 8px; border-radius: 6px; font-size: 0.7rem; font-weight: bold; }
		.node-tag { background: rgba(56, 189, 248, 0.12); color: #7dd3fc; border-radius: 6px; padding: 1px 6px; margin-inline-end: 4px; display: inline-block; }
		.form-select, .form-select:focus { background-color: #0f172a; color: #f1f5f9; border-color: rgba(255,255,255,0.15); }
		.node-link { direction: ltr; unicode-bidi: plaintext; }
	</style>
</head>
<body class="container py-4">
	<header class="text-center mb-5 animate__animated animate__fadeInDown">
		<h1 class="display-4 fw-bold mb-0" dir="ltr">MEGACODE<span class="text-primary">.ULTRA</span></h1>
		<p class="text-secondary mb-4" data-i18n="subtitle">Configs tested through real proxy connections</p>

		<div class="row justify-content-center g-3 mb-4">
			<div class="col-12 col-sm-6 col-md-3">
				<label for="lang-select" class="form-label small text-secondary"><i class="fas fa-language"></i> <span data-i18n="language">Language</span></label>
				<select id="lang-select" class="form-select"></select>
			</div>
			<div class="col-12 col-sm-6 col-md-3">
				<label for="country-select" class="form-label small text-secondary"><i class="fas fa-globe"></i> <span data-i18n="country">Country</span></label>
				<select id="country-select" class="form-select">
					<option value="" data-i18n="all_countries">All countries (best 50)</option>` + countryOptions(countries) + `
				</select>
			</div>
			<div class="col-12 col-md-3 d-flex align-items-end justify-content-center">
				<div class="form-check form-switch mb-2">
					<input class="form-check-input" type="checkbox" role="switch" id="static-only">
					<label class="form-check-label" for="static-only" data-i18n="static_only">Static IP only</label>
				</div>
			</div>
		</div>

		<div class="row justify-content-center g-3 mb-4">
			<div class="col-6 col-md-3">
				<div class="stat-box">
					<div class="small text-secondary text-uppercase" data-i18n="healthy">Healthy configs</div>
					<div class="h4 mb-0 text-primary" id="visible-count">` + fmt.Sprint(healthy) + `</div>
				</div>
			</div>
			<div class="col-6 col-md-3">
				<div class="stat-box">
					<div class="small text-secondary text-uppercase" data-i18n="latency">Best ping</div>
					<div class="h4 mb-0 text-success" dir="ltr">` + fmt.Sprint(bestPing) + `ms</div>
				</div>
			</div>
		</div>

		<div class="d-flex flex-wrap justify-content-center gap-3">
			<a id="main-sub" href="cleaned_configs.txt" download class="btn-main"><i class="fas fa-download me-2"></i> <span data-i18n="get">Get configs</span></a>` + categoryLinks() + `
			<button class="btn btn-outline-light rounded-pill px-4" data-bs-toggle="modal" data-bs-target="#helpModal" data-i18n="help">Help</button>
		</div>
	</header>

	<div class="row g-3" id="nodes">`

	var cardsHTML strings.Builder
	for _, card := range cards {
		node := card.Node
		escaped := html.EscapeString(card.Link)
		top, static := "0", "0"
		if card.Top {
			top = "1"
		}
		if node.StaticIP {
			static = "1"
		}
		exitIP := node.ExitIP
		if exitIP == "" {
			exitIP = "-"
		}
		fmt.Fprintf(&cardsHTML, `
		<div class="col-md-6 col-lg-4 node-card" data-country="%s" data-top="%s" data-static="%s">
			<div class="glass-card p-4 h-100">
				<div class="d-flex justify-content-between align-items-center mb-3">
					<span class="node-rank">#<span class="node-index"></span></span>
					<div class="text-success small fw-bold"><i class="fas fa-shield-alt"></i> <span class="country-name" data-code="%s">%s</span> <span dir="ltr">%dms</span></div>
				</div>
				<p class="small mb-2">%s</p>
				<p class="small text-secondary mb-2"><span data-i18n="exit_ip">Exit IP</span>: <span dir="ltr">%s</span></p>
				<p class="small text-truncate text-secondary mb-4 node-link">%s</p>
				<button class="btn btn-sm btn-primary w-100 rounded-3 py-2 copy-btn" data-config="%s">
					<i class="far fa-copy me-1"></i> <span data-i18n="copy">Copy configuration</span>
				</button>
			</div>
		</div>`, node.Location, top, static, node.Location, node.Location, node.Delay, tagBadges(node), html.EscapeString(exitIP), escaped, escaped)
	}

	htmlFooter := `
	</div>
	<p id="empty-message" class="text-center text-secondary my-5 d-none" data-i18n="no_nodes">No configs match this filter in the latest test run.</p>

	<div class="text-center mt-5">
		<div class="visitor-badge animate__animated animate__fadeIn">
			<i class="fas fa-chart-line text-primary me-2"></i> GLOBAL REACH:
			<img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https://rostamimohammadamin8.github.io/Starlink-V2ray-Scanner/&count_bg=%233B82F6&title_bg=%23080C14&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=true" alt="Hits" style="vertical-align: middle;"/>
		</div>
		<p class="small text-muted mt-3"><span data-i18n="last_update">Last update</span>: <span dir="ltr">` + now + `</span></p>
	</div>

	<div class="modal fade" id="helpModal" tabindex="-1">
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content bg-dark border-secondary">
				<div class="modal-header border-secondary">
					<h5 class="modal-title" data-i18n="help_title">Usage instructions</h5>
					<button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
				</div>
				<div class="modal-body">
					<p data-i18n="help1">Copy: choose a config and press the copy button.</p>
					<p data-i18n="help2">Import: open your client (v2rayNG, v2rayN, Hiddify) and import from the clipboard.</p>
					<p data-i18n="help3">Auto-update: add the download link as a subscription in your client.</p>
				</div>
			</div>
		</div>
	</div>

	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
	<script src="assets/i18n.js"></script>
	<script src="assets/panel.js"></script>
</body>
</html>
`

	return htmlHeader + cardsHTML.String() + htmlFooter
}

func countryOptions(countries []countryCount) string {
	var options strings.Builder
	for _, c := range countries {
		fmt.Fprintf(&options, `
					<option value="%s" data-count="%d">%s (%d)</option>`, c.Code, c.Count, c.Code, c.Count)
	}
	return options.String()
}

func tagBadges(node rankedNode) string {
	var badges strings.Builder
	for _, tag := range nodeTagList(node) {
		title := ""
		if tag.Key == "static" {
			title = ` data-i18n-title="static_hint"`
		}
		fmt.Fprintf(&badges, `<span class="node-tag" data-i18n="tag_%s"%s>%s</span>`, tag.Key, title, html.EscapeString(tag.Label))
	}
	return badges.String()
}

func categoryLinks() string {
	links := ""
	for _, c := range categories {
		links += fmt.Sprintf(`
			<a href="%s/%s" download class="btn btn-outline-light rounded-pill px-3" data-i18n="tag_%s">%s</a>`, CategoryDir, c.File, c.Key, html.EscapeString(c.Name))
	}
	return links
}
