package ai

import (
	"encoding/json"
	"fmt"
)

func BuildPrompt(projectName string, userQuery string, data interface{}, collection string) string {

	dataJSON, _ := json.MarshalIndent(data, "", "  ")

	return fmt.Sprintf(`
You are a data assistant.

IMPORTANT DEFINITIONS:
- "project" (system) = %s
- "%s" (collection) = database records

STRICT RULES:
- Use ONLY the provided DATA
- Do NOT confuse system project with database records
- If count is provided, return EXACT number
- Do NOT guess or assume

DATA:
%s

USER QUERY:
%s
`, projectName, collection, string(dataJSON), userQuery)
}