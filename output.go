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

	if len(nodes) > MaxNodes {
		nodes = nodes[:MaxNodes]
	}
	final := renderNodes(nodes)
	if err := os.WriteFile("index.html", []byte(generateFinalPanel(nodes, final, nodes[0].Delay)), 0644); err != nil {
		return err
	}
	return writeList("cleaned_configs.txt", final)
}
