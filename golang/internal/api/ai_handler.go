package api

import (
	"encoding/json"
	"net/http"

	"aichatbot/internal/ai"
	"aichatbot/internal/db"
	"aichatbot/internal/orchestrator"
	"aichatbot/internal/registry"
	"aichatbot/internal/security"
)

type AIRequest struct {
	ProjectName string `json:"projectName"`
	Query       string `json:"query"`
}

func AIQueryHandler(w http.ResponseWriter, r *http.Request) {

	var req AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", 400)
		return
	}

	// 1. LOAD PROJECT
	project, exists := registry.GetProject(req.ProjectName)
	if !exists {
		http.Error(w, "project not found", 404)
		return
	}

	// 2. SECURITY
	if err := security.ValidateQuery(req.Query); err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	if err := security.ValidateProjectAccess(project); err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	// 3. LOAD SCHEMA
	schema, err := db.LoadSchema(project.ProjectName)
	if err != nil {
		http.Error(w, "schema not found: "+err.Error(), 500)
		return
	}

	// 4. DETECT COLLECTIONS
	detectedCollections := db.DetectCollections(req.Query, schema)

	// 🚨 STRICT MODE (recommended)
	if len(detectedCollections) == 0 {
		http.Error(w, "no valid collection detected from query", 400)
		return
	}

	// 5. INTENT
	intent := orchestrator.ParseIntent(req.Query)

	// 6. BUILD QUERY
	dbQuery := orchestrator.BuildQuery(
		intent.Operation,
		detectedCollections,
	)

	// 7. VALIDATE ACTION
	if err := security.ValidateAction(dbQuery.Action); err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	// 8. EXECUTE MULTI-COLLECTION
	results := make(map[string]interface{})

	if project.DBType == "mongodb" {

		executor, err := db.NewMongoExecutor(project.DBUri, project.DBName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, col := range dbQuery.Collections {

			// 🔐 SECURITY PER COLLECTION
			if !security.IsCollectionAllowed(project.ProjectName, col) {
				continue
			}

			if dbQuery.Action == "count" {

				count, err := executor.CountDocuments(col)
				if err != nil {
					continue
				}

				results[col] = map[string]interface{}{
					"count": count,
				}

			} else {

				data, err := executor.FindAll(col)
				if err != nil {
					continue
				}

				results[col] = data
			}
		}
	}

	// 9. BUILD PROMPT
	prompt := ai.BuildPrompt(
		project.ProjectName,
		req.Query,
		results,
		dbQuery.Collections,
	)

	// 10. CALL AI
	aiResponse, err := ai.Generate(prompt, project.AIModel)
	if err != nil {
		aiResponse = "AI unavailable. Showing raw data."
	}

	// 11. RESPONSE
	response := map[string]interface{}{
		"project":     project.ProjectName,
		"intent":      intent,
		"query":       req.Query,
		"dbQuery":     dbQuery,
		"data":        results,
		"ai_response": aiResponse,
		"status":      "AI RESPONSE GENERATED",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}