package api

import (
	"encoding/json"
	"net/http"

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

	// 2. SECURITY CHECK
	if err := security.ValidateQuery(req.Query); err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	if err := security.ValidateProjectAccess(project); err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	// 3. LOAD SCHEMA (🔥 NEW)
	schema, err := db.LoadSchema(project.ProjectName)
	if err != nil {
		http.Error(w, "schema not found: "+err.Error(), 500)
		return
	}

	// 4. DETECT COLLECTION FROM QUERY (🔥 NEW)
	detectedCollection, err := db.DetectCollection(req.Query, schema)
	if err != nil {
		// fallback safe default
		detectedCollection = "projects"
	}

	// 5. INTENT
	intent := orchestrator.ParseIntent(req.Query)

	// 6. BUILD DB QUERY (FIXED)
	dbQuery := orchestrator.BuildQuery(
		req.Query,
		detectedCollection,
		intent.Operation,
	)

	// 7. SECURITY VALIDATION
	if err := security.ValidateAction(dbQuery.Action); err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	if !security.IsCollectionAllowed(project.ProjectName, dbQuery.Collection) {
		http.Error(w, "collection not allowed", 403)
		return
	}

	// 8. EXECUTION
	var result interface{}

	if project.DBType == "mongodb" {

		executor, err := db.NewMongoExecutor(project.DBUri, project.DBName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if dbQuery.Action == "count" {
			count, err := executor.CountDocuments(dbQuery.Collection)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			result = map[string]interface{}{
				"count": count,
			}
		} else {
			data, err := executor.FindAll(dbQuery.Collection)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			result = data
		}
	}

	// 9. RESPONSE
	response := map[string]interface{}{
		"project":  project.ProjectName,
		"intent":   intent,
		"query":    req.Query,
		"dbQuery":  dbQuery,
		"result":   result,
		"status":   "SCHEMA-DRIVEN EXECUTION DONE",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}