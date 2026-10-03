package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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

	byCountry := map[string][]rankedNode{}
	for _, node := range nodes {
		if node.Location != "XX" && len(byCountry[node.Location]) < MaxNodes {
			byCountry[node.Location] = append(byCountry[node.Location], node)
		}
	}
	if err := os.RemoveAll(CountryDir); err != nil {
		return err
	}
	if err := os.MkdirAll(CountryDir, 0755); err != nil {
		return err
	}
	rendered := map[string]string{}
	var countries []countryCount
	for code, list := range byCountry {
		links := renderNodes(list)
		for i, node := range list {
			rendered[node.Link] = links[i]
		}
		if err := writeList(filepath.Join(CountryDir, code+".txt"), links); err != nil {
			return err
		}
		countries = append(countries, countryCount{code, len(list)})
	}
	sort.Slice(countries, func(i, j int) bool { return countries[i].Code < countries[j].Code })
	fmt.Printf("%4d country lists in %s\n", len(countries), CountryDir)

	top := nodes
	if len(top) > MaxNodes {
		top = top[:MaxNodes]
	}
	final := renderNodes(top)
	inTop := map[string]bool{}
	for i, node := range top {
		rendered[node.Link] = final[i]
		inTop[node.Link] = true
	}

	var cards []panelNode
	perCountry := map[string]int{}
	for _, node := range nodes {
		if !inTop[node.Link] {
			if node.Location == "XX" || perCountry[node.Location] >= MaxPanelPerCountry {
				continue
			}
		}
		perCountry[node.Location]++
		cards = append(cards, panelNode{node, rendered[node.Link], inTop[node.Link]})
	}

	if err := os.WriteFile("index.html", []byte(generateFinalPanel(cards, countries, len(top), top[0].Delay)), 0644); err != nil {
		return err
	}
	return writeList("cleaned_configs.txt", final)
}
