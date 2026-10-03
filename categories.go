package main

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
