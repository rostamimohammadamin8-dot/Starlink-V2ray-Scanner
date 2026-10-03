package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	MaxNodes        = 50
	MaxPerSource    = 150
	MaxGamingDelay  = 200
	MaxResponseSize = 32 << 20
	CandidatesFile  = "candidates.txt"
	ResultsFile     = "results.jsonl"
	CategoryDir     = "subs"
)

// The test URLs must match the --test-urls list passed to xray-knife in the workflow.
const (
	baseTestLabel      = "cloudflare.com"
	openAITestLabel    = "api.openai.com"
	anthropicTestLabel = "api.anthropic.com"
	geminiTestLabel    = "gemini.google.com"
	youtubeTestLabel   = "youtube.com"
	instagramTestLabel = "instagram.com"
)

// Countries where ChatGPT, Claude or Gemini refuse service.
var aiRestrictedCountries = map[string]bool{
	"IR": true, "CN": true, "HK": true, "MO": true, "RU": true,
	"BY": true, "KP": true, "SY": true, "CU": true,
}

// Sources were selected by sampling up to 60 configs from each public
// collector and keeping those where at least 40% passed a real proxy test.
var sources = []string{
	"https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/v2ray-base64.txt",
	"https://raw.githubusercontent.com/igareck/vpn-configs-for-russia/main/BLACK_VLESS_RUS.txt",
	"https://raw.githubusercontent.com/F0rc3Run/F0rc3Run/main/Best-Results/proxies.txt",
	"https://raw.githubusercontent.com/Idolvpn/Automate-V2ray-Config-Collector/main/configs/lite_mix_sub.txt",
	"https://raw.githubusercontent.com/Pawdroid/Free-servers/main/sub",
	"https://raw.githubusercontent.com/Danialsamadi/v2go/main/AllConfigsSub.txt",
	"https://raw.githubusercontent.com/3yed-61/V2rayCollector/main/sub/vmess",
	"https://raw.githubusercontent.com/hans-thomas/v2ray-subscription/master/servers.txt",
	"https://raw.githubusercontent.com/Mahdi0024/ProxyCollector/master/sub/proxies.txt",
	"https://raw.githubusercontent.com/YawStar/Proxy-Hunter/main/configs/proxy_configs.txt",
	"https://raw.githubusercontent.com/cbusifabcap/daily_free_vpn/main/Z.txt",
	"https://raw.githubusercontent.com/RKPchannel/RKP_bypass_configs/main/blacklist.txt",
	"https://raw.githubusercontent.com/ShatakVPN/ConfigForge-V2Ray/main/configs/all.txt",
	"https://raw.githubusercontent.com/zhuhaiuk/free-nodes/main/nodes.txt",
}

var supportedSchemes = map[string]bool{
	"vmess": true, "vless": true, "trojan": true, "ss": true,
	"hysteria2": true, "hy2": true, "tuic": true,
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

type endpointResult struct {
	Label   string `json:"label"`
	Outcome string `json:"outcome"`
	Code    int    `json:"code"`
	Delay   int64  `json:"delay"`
}

type testResult struct {
	Link      string           `json:"link"`
	Location  string           `json:"location"`
	Endpoints []endpointResult `json:"endpoints"`
}

type rankedNode struct {
	Link     string
	Location string
	Delay    int64
	AI       bool
	YouTube  bool
	Insta    bool
	Gaming   bool
}

type category struct {
	Name    string
	File    string
	Matches func(rankedNode) bool
}

var categories = []category{
	{"AI (ChatGPT, Claude, Gemini)", "ai.txt", func(n rankedNode) bool { return n.AI }},
	{"Gaming (low ping)", "gaming.txt", func(n rankedNode) bool { return n.Gaming }},
	{"YouTube", "youtube.txt", func(n rankedNode) bool { return n.YouTube }},
	{"Instagram", "instagram.txt", func(n rankedNode) bool { return n.Insta }},
}

func generateFinalPanel(nodes []rankedNode, configs []string, bestPing int64) string {
	now := time.Now().Format("Jan 02, 15:04")

	htmlHeader := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>MegaCode Ultra Dashboard</title>
		<link rel="manifest" href='data:application/manifest+json,{"name":"MegaCode Ultra","short_name":"MegaCode","start_url":".","display":"standalone","background_color":"#080c14","theme_color":"#3b82f6"}'>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.1.1/animate.min.css"/>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
		<style>
			:root { --primary: #3b82f6; --bg: #080c14; --glass: rgba(30, 41, 59, 0.5); }
			body { background: var(--bg); color: #f1f5f9; font-family: 'Inter', sans-serif; }
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
				<a href="cleaned_configs.txt" download class="btn-main"><i class="fas fa-download me-2"></i> Get Configs</a>` + categoryLinks() + `
				<button class="btn btn-outline-light rounded-pill px-4" data-bs-toggle="modal" data-bs-target="#helpModal">Help</button>
			</div>
		</header>

		<div class="row g-3">`

	cards := ""
	for i, conf := range configs {
		escaped := html.EscapeString(conf)
		tags := nodeTags(nodes[i])
		cards += fmt.Sprintf(`
		<div class="col-md-6 col-lg-4 animate__animated animate__fadeInUp">
			<div class="glass-card p-4 h-100">
				<div class="d-flex justify-content-between align-items-center mb-3">
					<span class="node-rank">RANK #%d</span>
					<div class="text-success small fw-bold"><i class="fas fa-shield-alt"></i> %s</div>
				</div>
				<p class="small text-info mb-2">%s</p>
				<p class="small text-truncate text-secondary mb-4">%s</p>
				<button class="btn btn-sm btn-primary w-100 rounded-3 py-2" data-config="%s" onclick="copyToClipboard(this.dataset.config)">
					<i class="far fa-copy me-1"></i> Copy Configuration
				</button>
			</div>
		</div>`, i+1, html.EscapeString(fmt.Sprintf("%s %dms", nodes[i].Location, nodes[i].Delay)), html.EscapeString(tags), escaped, escaped)
	}

	htmlFooter := `
		</div>

		<div class="text-center mt-5">
			<div class="visitor-badge animate__animated animate__fadeIn">
				<i class="fas fa-chart-line text-primary me-2"></i> GLOBAL REACH: 
				<img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https://rostamimohammadamin8.github.io/Starlink-V2ray-Scanner/&count_bg=%233B82F6&title_bg=%23080C14&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=true" alt="Hits" style="vertical-align: middle;"/>
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

func fetchSource(source string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "v2rayNG/1.8")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseSize))
	if err != nil {
		return nil, err
	}
	return parseSubscription(body), nil
}

func parseSubscription(body []byte) []string {
	if links := extractLinks(string(body)); len(links) > 0 {
		return links
	}
	compact := strings.Join(strings.Fields(string(body)), "")
	for _, encoding := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		if decoded, err := encoding.DecodeString(compact); err == nil {
			return extractLinks(string(decoded))
		}
	}
	return nil
}

func extractLinks(text string) []string {
	var links []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if isValidConfig(line) {
			links = append(links, line)
		}
	}
	return links
}

func isValidConfig(link string) bool {
	scheme, rest, found := strings.Cut(link, "://")
	if !found || !supportedSchemes[strings.ToLower(scheme)] || rest == "" {
		return false
	}
	if strings.EqualFold(scheme, "vmess") {
		_, err := decodeVmess(link)
		return err == nil
	}
	parsed, err := url.Parse(link)
	return err == nil && parsed.Hostname() != "" && parsed.Port() != ""
}

func decodeVmess(link string) (map[string]json.RawMessage, error) {
	_, payload, _ := strings.Cut(link, "://")
	payload, _, _ = strings.Cut(payload, "#")
	for _, encoding := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		decoded, err := encoding.DecodeString(payload)
		if err != nil {
			continue
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(decoded, &fields); err != nil {
			return nil, err
		}
		var address string
		if err := json.Unmarshal(fields["add"], &address); err != nil || address == "" {
			return nil, errors.New("vmess config has no address")
		}
		return fields, nil
	}
	return nil, errors.New("vmess payload is not valid base64")
}

// connectionKey identifies a server regardless of its remark, so re-labelled copies are dropped.
func connectionKey(link string) string {
	if strings.HasPrefix(strings.ToLower(link), "vmess://") {
		fields, err := decodeVmess(link)
		if err != nil {
			return link
		}
		var parts []string
		for _, name := range []string{"add", "port", "id", "net", "type", "host", "path", "tls", "sni"} {
			parts = append(parts, strings.Trim(string(fields[name]), `"`))
		}
		return "vmess://" + strings.Join(parts, "|")
	}
	base, _, _ := strings.Cut(link, "#")
	return base
}

func collect() error {
	type sourceResult struct {
		source string
		links  []string
		err    error
	}
	results := make([]sourceResult, len(sources))
	var wg sync.WaitGroup
	for i, source := range sources {
		wg.Add(1)
		go func(i int, source string) {
			defer wg.Done()
			links, err := fetchSource(source)
			results[i] = sourceResult{source, links, err}
		}(i, source)
	}
	wg.Wait()

	seen := map[string]bool{}
	var candidates []string
	for _, result := range results {
		if result.err != nil {
			fmt.Printf("skip %s: %v\n", result.source, result.err)
			continue
		}
		rand.Shuffle(len(result.links), func(i, j int) {
			result.links[i], result.links[j] = result.links[j], result.links[i]
		})
		added := 0
		for _, link := range result.links {
			if added >= MaxPerSource {
				break
			}
			if key := connectionKey(link); !seen[key] {
				seen[key] = true
				candidates = append(candidates, link)
				added++
			}
		}
		fmt.Printf("%4d candidates from %s\n", added, result.source)
	}
	if len(candidates) == 0 {
		return errors.New("no candidates collected from any source")
	}
	return os.WriteFile(CandidatesFile, []byte(strings.Join(candidates, "\n")+"\n"), 0644)
}

func readResults(path string) ([]rankedNode, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nodes []rankedNode
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		var result testResult
		if err := json.Unmarshal(scanner.Bytes(), &result); err != nil {
			continue
		}
		if node, ok := classify(result); ok {
			nodes = append(nodes, node)
		}
	}
	return nodes, scanner.Err()
}

func classify(result testResult) (rankedNode, bool) {
	endpoints := map[string]endpointResult{}
	for _, endpoint := range result.Endpoints {
		endpoints[endpoint.Label] = endpoint
	}
	base, ok := endpoints[baseTestLabel]
	if !ok || base.Outcome != "ok" || base.Delay <= 0 {
		return rankedNode{}, false
	}
	location := strings.ToUpper(result.Location)
	if location == "" || location == "NULL" {
		location = "XX"
	}
	// Unauthenticated API calls answer 401 from supported regions and 403 from blocked ones.
	aiAllowed := endpoints[openAITestLabel].Code == http.StatusUnauthorized &&
		endpoints[anthropicTestLabel].Code == http.StatusUnauthorized &&
		endpoints[geminiTestLabel].Code == http.StatusOK &&
		!aiRestrictedCountries[location]
	instagram := endpoints[instagramTestLabel].Code
	return rankedNode{
		Link:     result.Link,
		Location: location,
		Delay:    base.Delay,
		AI:       aiAllowed,
		YouTube:  endpoints[youtubeTestLabel].Code == http.StatusNoContent,
		Insta:    (instagram >= 200 && instagram < 400) || instagram == http.StatusTooManyRequests,
		Gaming:   base.Delay <= MaxGamingDelay,
	}, true
}

func nodeTags(node rankedNode) string {
	var tags []string
	if node.AI {
		tags = append(tags, "AI")
	}
	if node.Gaming {
		tags = append(tags, "Gaming")
	}
	if node.YouTube {
		tags = append(tags, "YouTube")
	}
	if node.Insta {
		tags = append(tags, "Instagram")
	}
	if len(tags) == 0 {
		return "Web"
	}
	return strings.Join(tags, " · ")
}

func categoryLinks() string {
	links := ""
	for _, c := range categories {
		links += fmt.Sprintf(`
				<a href="%s/%s" download class="btn btn-outline-light rounded-pill px-3">%s</a>`, CategoryDir, c.File, html.EscapeString(c.Name))
	}
	return links
}

func withRemark(link, remark string) string {
	if strings.HasPrefix(strings.ToLower(link), "vmess://") {
		fields, err := decodeVmess(link)
		if err != nil {
			return link
		}
		encodedRemark, err := json.Marshal(remark)
		if err != nil {
			return link
		}
		fields["ps"] = encodedRemark
		encoded, err := json.Marshal(fields)
		if err != nil {
			return link
		}
		return "vmess://" + base64.StdEncoding.EncodeToString(encoded)
	}
	base, _, _ := strings.Cut(link, "#")
	return base + "#" + url.PathEscape(remark)
}

func renderNodes(nodes []rankedNode) []string {
	links := make([]string, 0, len(nodes))
	for i, node := range nodes {
		remark := fmt.Sprintf("%s-%02d ⚡%dms %s-RankedByMegaCode", node.Location, i+1, node.Delay, nodeTags(node))
		links = append(links, withRemark(node.Link, remark))
	}
	return links
}

func writeList(path string, links []string) error {
	content := ""
	if len(links) > 0 {
		content = strings.Join(links, "\n") + "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func build() error {
	nodes, err := readResults(ResultsFile)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return errors.New("no config passed the proxy test; keeping previous output")
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Delay < nodes[j].Delay })
	seen := map[string]bool{}
	unique := nodes[:0]
	for _, node := range nodes {
		if key := connectionKey(node.Link); !seen[key] {
			seen[key] = true
			unique = append(unique, node)
		}
	}
	nodes = unique

	if err := os.MkdirAll(CategoryDir, 0755); err != nil {
		return err
	}
	for _, c := range categories {
		var matched []rankedNode
		for _, node := range nodes {
			if c.Matches(node) && len(matched) < MaxNodes {
				matched = append(matched, node)
			}
		}
		if err := writeList(filepath.Join(CategoryDir, c.File), renderNodes(matched)); err != nil {
			return err
		}
		fmt.Printf("%4d configs in %s\n", len(matched), c.File)
	}

	if len(nodes) > MaxNodes {
		nodes = nodes[:MaxNodes]
	}
	final := renderNodes(nodes)
	if err := os.WriteFile("index.html", []byte(generateFinalPanel(nodes, final, nodes[0].Delay)), 0644); err != nil {
		return err
	}
	return writeList("cleaned_configs.txt", final)
}

func main() {
	mode := "collect"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	var err error
	switch mode {
	case "collect":
		err = collect()
	case "build":
		err = build()
	default:
		err = fmt.Errorf("unknown mode %q (use collect or build)", mode)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("Build Successful!")
}
