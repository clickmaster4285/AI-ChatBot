package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "AI ERP Running")
	})

	fmt.Println("Server running on :" + port)
	http.ListenAndServe(":"+port, nil)
}