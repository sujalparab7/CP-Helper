package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	resp, _ := http.Get("https://generativelanguage.googleapis.com/v1beta/models?key=" + key)
	
	var res struct {
		Models []struct {
			Name            string   `json:"name"`
			SupportedMethods []string `json:"supportedGenerationMethods"`
		} `json:"models"`
	}
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
	    fmt.Println(err)
	    return
	}
    
	for _, m := range res.Models {
        for _, method := range m.SupportedMethods {
            if method == "generateContent" {
                fmt.Println(m.Name)
                return
            }
        }
	}
	fmt.Println("No supported models found.")
}
