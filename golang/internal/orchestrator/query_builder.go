package orchestrator

type DBQuery struct {
	Action      string   `json:"action"`
	Collections []string `json:"collections"`
}

func BuildQuery(action string, detectedCollections []string) DBQuery {

	return DBQuery{
		Action:      action,
		Collections: detectedCollections,
	}
}