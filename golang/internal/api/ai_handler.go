package api

import (
	"encoding/json"
	"net/http"

	"aichatbot/internal/registry"
)

type AIRequest struct {
	ProjectName string `json:"projectName"`
	Query       string `json:"query"`
}

func AIQueryHandler(w http.ResponseWriter, r *http.Request) {

	var req AIRequest
	json.NewDecoder(r.Body).Decode(&req)

	// 1. Find project
	project, exists := registry.GetProject(req.ProjectName)
	if !exists {
		http.Error(w, "Project not found", 404)
		return	
	}

	// 2. Build response (TEMP test)
	response := map[string]interface{}{
		"project": project.ProjectName,
		"db":      project.DBType,
		"query":   req.Query,
		"status":  "project resolved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}