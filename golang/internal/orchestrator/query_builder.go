package orchestrator

import "strings"

type DBQuery struct {
	Action     string
	Collection string
}

func BuildQuery(query string, detectedCollection string, action string) DBQuery {

	// fallback safety
	collection := "projects"

	if detectedCollection != "" {
		collection = detectedCollection
	}

	q := strings.ToLower(query)

	// override action if needed
	if strings.Contains(q, "total") ||
		strings.Contains(q, "count") ||
		strings.Contains(q, "length") {
		action = "count"
	}

	return DBQuery{
		Action:     action,
		Collection: collection,
	}
}