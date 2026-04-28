package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"aichatbot/internal/db"
	"aichatbot/internal/orchestrator"
	"aichatbot/internal/registry"
)

type AIRequest struct {
	ProjectName string `json:"projectName"`
	Query       string `json:"query"`
}

func AIQueryHandler(w http.ResponseWriter, r *http.Request) {

	var req AIRequest
	json.NewDecoder(r.Body).Decode(&req)

	// 1. LOAD PROJECT
	project, exists := registry.GetProject(req.ProjectName)
	if !exists {
		http.Error(w, "Project not found", 404)
		return
	}

	// 2. INTENT
	intent := orchestrator.ParseIntent(req.Query)

	// 3. QUERY BUILDER
	dbQuery := orchestrator.BuildQuery(req.Query)

	// 4. CONNECT MONGO (ONLY FOR MONGODB PROJECTS)
	var result interface{}

	if project.DBType == "mongodb" {

		executor, err := db.NewMongoExecutor(project.DBUri, project.DBName)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if dbQuery.Action == "count" {
			fmt.Println("the dbQuery.Collection is ===: %s", dbQuery.Collection)
			count, err := executor.CountDocuments(dbQuery.Collection)
			if err != nil {
				fmt.Println("Mongo ERROR:", err)
				http.Error(w, err.Error(), 500)
				return
			}
			result = map[string]interface{}{
				"count": count,
			}
		} else {
			data, _ := executor.FindAll(dbQuery.Collection)
			result = data
		}
	}

	// 5. RESPONSE
	response := map[string]interface{}{
		"project": project.ProjectName,
		"intent":  intent,
		"query":   req.Query,
		"dbQuery": dbQuery,
		"result":  result,
		"status":  "REAL DB EXECUTION DONE",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
