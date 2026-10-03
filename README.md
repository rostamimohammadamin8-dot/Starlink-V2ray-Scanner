# 🚀 MegaCode Starlink & Gaming Scanner

![Starlink Scanner Status](https://github.com/rostamimohammadamin8/Starlink-Scanner/actions/workflows/runner.yml/badge.svg)

This is an automated **MegaCode** written in **Go** that scans, tests, and filters high-speed V2Ray configs with Starlink standards.

### 💎 Subscription Link
Copy this link and paste it into your **v2rayNG** or **v2rayN**:

```
https://raw.githubusercontent.com/rostamimohammadamin8-dot/Starlink-V2ray-Scanner/main/cleaned_configs.txt
```

### 🎯 Use-case subscriptions
Each config is tested through a real proxy connection against several services, and the ones that pass are also published per use case (fastest first, up to 50 each):

| Use case | Link | Test |
|---|---|---|
| AI (ChatGPT, Claude, Gemini) | `https://raw.githubusercontent.com/rostamimohammadamin8-dot/Starlink-V2ray-Scanner/main/subs/ai.txt` | OpenAI and Anthropic APIs answer `401` (not a region block), Gemini answers `200`, and the exit country is not one where these services are restricted |
| Online gaming | `https://raw.githubusercontent.com/rostamimohammadamin8-dot/Starlink-V2ray-Scanner/main/subs/gaming.txt` | Delay through the proxy is at most 200 ms |
| YouTube | `https://raw.githubusercontent.com/rostamimohammadamin8-dot/Starlink-V2ray-Scanner/main/subs/youtube.txt` | `youtube.com/generate_204` answers `204` |
| Instagram | `https://raw.githubusercontent.com/rostamimohammadamin8-dot/Starlink-V2ray-Scanner/main/subs/instagram.txt` | `instagram.com` answers through the proxy |
| Static IP | `https://raw.githubusercontent.com/rostamimohammadamin8-dot/Starlink-V2ray-Scanner/main/subs/static-ip.txt` | The exit IP seen by the test sites equals the server's own IP (no CDN or rotating pool in front) |

### 🌍 Per-country subscriptions
Every exit country found in the latest run gets its own list (up to 50, fastest first):

```
https://raw.githubusercontent.com/rostamimohammadamin8-dot/Starlink-V2ray-Scanner/main/subs/countries/<CODE>.txt
```

`<CODE>` is the two-letter country code, e.g. `DE`, `NL`, `TR`, `US`. The set of countries changes from run to run.

"Static IP" means the IP stays the same while that server stays online. Free public servers come and go, so a truly permanent IP needs your own server.

`cleaned_configs.txt` stays the general list for web browsing.

### ⚙️ How it works
Every 6 hours the workflow:
1. `go run . collect` downloads the sources listed in `sources.go`, decodes plain or base64 subscriptions, keeps only valid `vmess`/`vless`/`trojan`/`ss`/`hysteria2`/`tuic` links and writes up to 150 random candidates per source to `candidates.txt`.
2. [xray-knife](https://github.com/lilendian0x00/xray-knife) connects through every candidate, requests each test URL above and records the results in `results.jsonl`.
3. `go run . build` keeps the 50 fastest configs that passed, regenerates `cleaned_configs.txt`, `index.html` and the `subs/` and `subs/countries/` lists.

### 🖥️ Dashboard
`index.html` has a language selector (20 languages, right-to-left for Persian, Arabic, Urdu and Pashto), a country filter and a "static IP only" switch. When you pick a country, the download button switches to that country's subscription.

Sources were chosen by testing a random sample from each public collector; only those where at least 40% of the sampled configs worked were kept. Results are measured from GitHub's servers, so availability from your own network may differ.

### 🗂️ Project structure
| File | Purpose |
|---|---|
| `main.go` | Entry point: `go run . collect` / `go run . build` |
| `config.go` | Limits, file names, test endpoint labels, AI-restricted countries |
| `sources.go` | Subscription sources and supported link schemes |
| `subscription.go` | Downloading, base64 decoding, link validation, server dedupe |
| `collect.go` | Builds `candidates.txt` from all sources |
| `results.go` | Reads xray-knife results and classifies each config |
| `categories.go` | Use-case lists written to `subs/` |
| `output.go` | Remarks, ranking and writing the output files |
| `panel.go` | Generates the `index.html` dashboard |
| `assets/panel.js` | Dashboard behaviour: language, country and static-IP filters, copy button |
| `assets/i18n.js` | Dashboard translations |
| `docs/` | Guides (Persian) |

### 🔐 Personal config without a VPS
To build your own private subscription on a free Cloudflare account (no VPS, no bank card), see the Persian guide: [docs/personal-config-cloudflare-fa.md](docs/personal-config-cloudflare-fa.md).
