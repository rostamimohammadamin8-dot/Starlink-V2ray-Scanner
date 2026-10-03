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

`cleaned_configs.txt` stays the general list for web browsing.

### ⚙️ How it works
Every 6 hours the workflow:
1. `go run main.go collect` downloads the sources listed in `main.go`, decodes plain or base64 subscriptions, keeps only valid `vmess`/`vless`/`trojan`/`ss`/`hysteria2`/`tuic` links and writes up to 150 random candidates per source to `candidates.txt`.
2. [xray-knife](https://github.com/lilendian0x00/xray-knife) connects through every candidate, requests each test URL above and records the results in `results.jsonl`.
3. `go run main.go build` keeps the 50 fastest configs that passed, regenerates `cleaned_configs.txt`, `index.html` and the `subs/` lists.

Sources were chosen by testing a random sample from each public collector; only those where at least 40% of the sampled configs worked were kept. Results are measured from GitHub's servers, so availability from your own network may differ.

### 🔐 Personal config without a VPS
To build your own private subscription on a free Cloudflare account (no VPS, no bank card), see the Persian guide: [docs/personal-config-cloudflare-fa.md](docs/personal-config-cloudflare-fa.md).
