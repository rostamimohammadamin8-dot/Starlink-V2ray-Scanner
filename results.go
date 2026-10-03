package main

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
)

type endpointResult struct {
	Label   string `json:"label"`
	Outcome string `json:"outcome"`
	Code    int    `json:"code"`
	Delay   int64  `json:"delay"`
}

type testResult struct {
	Link     string `json:"link"`
	Location string `json:"location"`
	// ExitIP is the address the test endpoints saw, reported by xray-knife as "ip".
	ExitIP   string `json:"ip"`
	Protocol struct {
		Address string `json:"address"`
	} `json:"protocol"`
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
	ExitIP   string
	StaticIP bool
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
	exitIP := net.ParseIP(result.ExitIP)
	serverIP := net.ParseIP(result.Protocol.Address)
	return rankedNode{
		Link:     result.Link,
		Location: location,
		Delay:    base.Delay,
		AI:       aiAllowed,
		YouTube:  endpoints[youtubeTestLabel].Code == http.StatusNoContent,
		Insta:    (instagram >= 200 && instagram < 400) || instagram == http.StatusTooManyRequests,
		Gaming:   base.Delay <= MaxGamingDelay,
		ExitIP:   ipString(exitIP),
		// Traffic leaves from the server itself rather than a CDN or rotating pool.
		StaticIP: exitIP != nil && serverIP != nil && exitIP.Equal(serverIP),
	}, true
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

type nodeTag struct {
	Key   string
	Label string
}

func nodeTagList(node rankedNode) []nodeTag {
	var tags []nodeTag
	if node.AI {
		tags = append(tags, nodeTag{"ai", "AI"})
	}
	if node.Gaming {
		tags = append(tags, nodeTag{"gaming", "Gaming"})
	}
	if node.YouTube {
		tags = append(tags, nodeTag{"youtube", "YouTube"})
	}
	if node.Insta {
		tags = append(tags, nodeTag{"instagram", "Instagram"})
	}
	if node.StaticIP {
		tags = append(tags, nodeTag{"static", "StaticIP"})
	}
	if len(tags) == 0 {
		tags = append(tags, nodeTag{"web", "Web"})
	}
	return tags
}

func nodeTags(node rankedNode) string {
	labels := make([]string, 0, 5)
	for _, tag := range nodeTagList(node) {
		labels = append(labels, tag.Label)
	}
	return strings.Join(labels, " · ")
}

func ipString(ip net.IP) string {
	if ip == nil {
		return ""
	}
	return ip.String()
}
