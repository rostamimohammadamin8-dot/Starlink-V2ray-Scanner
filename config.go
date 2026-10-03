package main

const (
	MaxNodes        = 50
	MaxPerSource    = 150
	MaxGamingDelay  = 200
	MaxResponseSize = 32 << 20
	CandidatesFile  = "candidates.txt"
	ResultsFile     = "results.jsonl"
	CategoryDir     = "subs"
	CountryDir      = "subs/countries"
	// Extra configs per country shown on the dashboard beyond the overall top list.
	MaxPanelPerCountry = 10
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
