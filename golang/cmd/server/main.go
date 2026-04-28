package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "AI Chat Bot Running")
	})

	http.HandleFunc("/ai-health", func(w http.ResponseWriter, r *http.Request) {

		payload := []byte(`{
			"model": "llama3:8b",
			"prompt": "Say hello",
			"stream": false
		}`)

		resp, err := http.Post(
			"http://ollama:11434/api/generate",
			"application/json",
			bytes.NewBuffer(payload),
		)

		if err != nil {
			http.Error(w, "Failed to call Ollama: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	fmt.Println("Server running on :" + port)
	http.ListenAndServe(":"+port, nil)
}