package db

import (
	"encoding/json"
	"os"
	"strings"
)

type Schema struct {
	Collections map[string]CollectionSchema `json:"collections"`
}

type CollectionSchema struct {
	Type    string   `json:"type"`
	Aliases []string `json:"aliases"`
}

// In-memory cache
var SchemaRegistry = make(map[string]Schema)

// LOAD SCHEMA PER PROJECT
func LoadSchema(projectName string) (Schema, error) {

	if schema, exists := SchemaRegistry[projectName]; exists {
		return schema, nil
	}

	path := "config/schemas/" + projectName + ".json"

	file, err := os.ReadFile(path)
	if err != nil {
		return Schema{}, err
	}

	var schema Schema
	if err := json.Unmarshal(file, &schema); err != nil {
		return Schema{}, err
	}

	SchemaRegistry[projectName] = schema
	return schema, nil
}

// CHECK IF COLLECTION EXISTS
func IsValidCollection(schema Schema, collection string) bool {
	_, exists := schema.Collections[collection]
	return exists
}

// SIMPLE ENTITY → COLLECTION DETECTION
func DetectCollections(query string, schema Schema) []string {

	q := strings.ToLower(query)
	var matches []string

	for colName, colSchema := range schema.Collections {

		// Check aliases
		for _, alias := range colSchema.Aliases {
			if strings.Contains(q, strings.ToLower(alias)) {
				matches = append(matches, colName)
				break
			}
		}
	}

	return matches
}