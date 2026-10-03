package main

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
