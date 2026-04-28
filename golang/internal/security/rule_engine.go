package security

import (
	"errors"
	"strings"

	"aichatbot/internal/registry"
)

var allowedOperations = map[string]bool{
	"count": true,
	"find":  true,
}

var blockedKeywords = []string{
	"drop",
	"delete",
	"update",
	"insert",
	"remove",
	"shutdown",
}

// CHECK IF QUERY IS SAFE
func ValidateQuery(query string) error {

	q := strings.ToLower(query)

	// 1. BLOCK DANGEROUS WORDS
	for _, word := range blockedKeywords {
		if strings.Contains(q, word) {
			return errors.New("forbidden operation detected")
		}
	}

	return nil
}

// VALIDATE DB ACTION
func ValidateAction(action string) error {

	if !allowedOperations[action] {
		return errors.New("operation not allowed: only READ operations are permitted")
	}

	return nil
}

// PROJECT ACCESS CHECK
func ValidateProjectAccess(project registry.ProjectConfig) error {

	if !project.OnlyGetApiAccess {
		return errors.New("project API access disabled")
	}

	return nil
}