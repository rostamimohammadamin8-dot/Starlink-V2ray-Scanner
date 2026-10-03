package main

type category struct {
	Key     string
	Name    string
	File    string
	Matches func(rankedNode) bool
}

var categories = []category{
	{"ai", "AI (ChatGPT, Claude, Gemini)", "ai.txt", func(n rankedNode) bool { return n.AI }},
	{"gaming", "Gaming (low ping)", "gaming.txt", func(n rankedNode) bool { return n.Gaming }},
	{"youtube", "YouTube", "youtube.txt", func(n rankedNode) bool { return n.YouTube }},
	{"instagram", "Instagram", "instagram.txt", func(n rankedNode) bool { return n.Insta }},
	{"static", "Static IP", "static-ip.txt", func(n rankedNode) bool { return n.StaticIP }},
}
