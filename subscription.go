package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

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
