package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"cp-advisor/handlers"
)

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			val := strings.TrimSpace(parts[1])
			val = strings.TrimSuffix(val, "\r")
			os.Setenv(strings.TrimSpace(parts[0]), val)
		}
	}
}

func main() {
	loadEnv()

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/api/analyze", handlers.HandleAnalyze)

	log.Println("Server listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
