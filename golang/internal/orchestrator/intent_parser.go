package orchestrator

import "strings"

type Intent struct {
	Type       string  `json:"type"`
	Entity     string  `json:"entity"`
	Operation  string  `json:"operation"`
	Confidence float64 `json:"confidence"`
}

// SIMPLE RULE-BASED INTENT PARSER (v1)
// later we upgrade to embedding + AI
func ParseIntent(query string) Intent {

	q := strings.ToLower(query)

	intent := Intent{
		Confidence: 0.6,
	}

	// COUNT / LENGTH / TOTAL
	if strings.Contains(q, "total") ||
		strings.Contains(q, "count") ||
		strings.Contains(q, "length") {

		intent.Type = "aggregation"
		intent.Operation = "count"
		intent.Entity = "unknown"
		intent.Confidence = 0.75
		return intent
	}

	// REPORT
	if strings.Contains(q, "report") {
		intent.Type = "report"
		intent.Operation = "fetch"
		intent.Entity = "unknown"
		intent.Confidence = 0.7
		return intent
	}

	// DEFAULT FALLBACK
	intent.Type = "raw"
	intent.Operation = "fetch"
	intent.Entity = "unknown"
	intent.Confidence = 0.4

	return intent
}