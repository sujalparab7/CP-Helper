package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	file, _ := os.Open(".env")
	if file != nil {
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
		file.Close()
	}

	key := os.Getenv("GEMINI_API_KEY")
	resp, err := http.Get("https://generativelanguage.googleapis.com/v1beta/models?key=" + key)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bytes))
}
