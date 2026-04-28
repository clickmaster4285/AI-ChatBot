package orchestrator

import "strings"

type DBQuery struct {
	Action     string
	Collection string
}

func BuildQuery(query string) DBQuery {

	q := strings.ToLower(query)

	// COUNT LOGIC
	if strings.Contains(q, "total") ||
		strings.Contains(q, "count") ||
		strings.Contains(q, "length") {

		return DBQuery{
			Action:     "count",
			Collection: "projects",
		}
	}

	// DEFAULT
	return DBQuery{
		Action:     "find",
		Collection: "projects",
	}
}