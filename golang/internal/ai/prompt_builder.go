package ai

import (
	"fmt"
	"strings"
)

func BuildPrompt(project string, query string, data interface{}, collections []string) string {

	return fmt.Sprintf(`
You are an AI assistant for project "%s".

User Query:
%s

Collections:
%s

Data:
%v

Rules:
- Do NOT guess
- Use ONLY the data provided
- If multiple collections exist, mention each clearly

Answer:
`, project, query, strings.Join(collections, ", "), data)
}