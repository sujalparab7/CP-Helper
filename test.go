package main

import (
	"bufio"
	"cp-advisor/services"
	"fmt"
	"os"
	"strings"
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
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
}

func main() {
	loadEnv()
	fmt.Println("API KEY:", os.Getenv("GEMINI_API_KEY"))

	res, err := services.AnalyzeUser("tourist")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("AIFeedback Level 2:", res.AIFeedback.Level2)
}
