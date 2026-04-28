package main

import (
	"aichatbot/internal/api"
	"aichatbot/internal/registry"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// 🔥 LOAD PROJECT REGISTRY (THIS WAS MISSING)
	err := registry.LoadProjects()
	if err != nil {
		log.Fatal("Failed to load project registry:", err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "AI Chat Bot Running")
	})

	http.HandleFunc("/api/v1/ai-query", api.AIQueryHandler)

	fmt.Println("Server running on :" + port)
	http.ListenAndServe(":"+port, nil)
}