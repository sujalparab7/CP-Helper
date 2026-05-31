package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type GeminiRequest struct {
	Contents          []GeminiContent    `json:"contents"`
	SystemInstruction *GeminiInstruction `json:"systemInstruction,omitempty"`
}

type GeminiInstruction struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []GeminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func GenerateAIFeedback(problem CFProblem, lang string, verdict string, code string, weaknesses []Weakness) AIFeedback {
	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		return GenerateMockFeedback(problem, lang)
	}

	var weakStrs []string
	for _, w := range weaknesses {
		weakStrs = append(weakStrs, w.Tag)
	}
	weakStr := strings.Join(weakStrs, ", ")

	var prompt string
	if code != "" {
		prompt = fmt.Sprintf(`You are an expert Socratic AI Tutor for Competitive Programming.
The user failed the problem "%s" (Tags: %s).
Their verdict was: %s.
Language: %s.
Overall recent contest weaknesses: %s.

Here is their source code:
%s

Analyze their logic. Do not give them the final code. Act strictly as a Socratic guide. Also, formulate a 5-day structured training matrix strictly based heavily on their weaknesses and this code's algorithmic flaws. Provide your response exactly in this JSON format without any markdown wrapper containing the json block:
{
  "problemName": "%s",
  "level1": "Conceptual hint without giving the data structure",
  "level2": "Structural hint about exact data structures or logic gaps",
  "level3": "Diagnostic feedback pointing out specific bugs or errors",
  "level4": "A failing edge case or concrete trace scenario",
  "matrix": [
    {"day": "Monday", "focus": "Specific algorithmic topic", "objective": "Specific core goal and reason", "action": "Exact problem set action to take (be highly specific)"},
    {"day": "Tuesday", "focus": "...", "objective": "...", "action": "..."}
  ] // Generate all 5 days up to Friday
}`, problem.Name, strings.Join(problem.Tags, ", "), verdict, lang, weakStr, code, problem.Name)
	} else {
		prompt = fmt.Sprintf(`You are an expert Socratic AI Tutor for Competitive Programming.
The user failed the problem "%s" (Tags: %s) with verdict: %s in %s.
Overall recent weaknesses: %s.
Due to platform limits, code couldn't be retrieved.

Analyze the common pitfalls for this algorithmic problem. Act strictly as a Socratic guide. Also, formulate a 5-day structured training matrix strictly based heavily upon their exact weaknesses. Provide your response exactly in this JSON format without any markdown wrapper containing the json block:
{
  "problemName": "%s",
  "level1": "Conceptual hint about the optimal algorithms needed for this specific problem",
  "level2": "Structural hint about how to store the data and avoid frequent errors",
  "level3": "Diagnostic feedback pointing out why %s commonly fails",
  "level4": "A failing edge case scenario they must mentally dry-run",
  "matrix": [
    {"day": "Monday", "focus": "Specific algorithmic topic", "objective": "Specific core goal and reason", "action": "Exact problem set action to take (be highly specific)"},
    {"day": "Tuesday", "focus": "...", "objective": "...", "action": "..."}
  ] // Generate all 5 days up to Friday
}`, problem.Name, strings.Join(problem.Tags, ", "), verdict, lang, weakStr, problem.Name, verdict)
	}

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{Role: "user", Parts: []GeminiPart{{Text: prompt}}},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return GenerateMockFeedback(problem, lang)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey
	
	var finalResp *http.Response
	maxRetries := 3
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		
		if err == nil && resp.StatusCode == http.StatusOK {
			finalResp = resp
			break
		}
		
		if err != nil {
			fmt.Printf("Gemini Attempt %d HTTP error: %v\n", attempt, err)
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Printf("Gemini Attempt %d HTTP Error %d: %s\n", attempt, resp.StatusCode, string(bodyBytes))
			resp.Body.Close()
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt*2) * time.Second) // Progressive backoff
		}
	}

	if finalResp == nil {
		return GenerateMockFeedback(problem, lang)
	}
	defer finalResp.Body.Close()

	var gemResp GeminiResponse
	if err := json.NewDecoder(finalResp.Body).Decode(&gemResp); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return GenerateMockFeedback(problem, lang)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		fmt.Printf("Empty candidates in Gemini response\n")
		return GenerateMockFeedback(problem, lang)
	}

	responseText := gemResp.Candidates[0].Content.Parts[0].Text

	responseText = strings.TrimSpace(responseText)
	if strings.HasPrefix(responseText, "```json") {
		responseText = strings.TrimPrefix(responseText, "```json")
	} else if strings.HasPrefix(responseText, "```") {
		responseText = strings.TrimPrefix(responseText, "```")
	}
	
	if strings.HasSuffix(responseText, "```") {
		responseText = strings.TrimSuffix(responseText, "```")
	}
	responseText = strings.TrimSpace(responseText)

	var feedback AIFeedback
	if err := json.Unmarshal([]byte(responseText), &feedback); err != nil {
		fmt.Printf("JSON Unmarshal Error: %v\nText was: %s\n", err, responseText)
		return GenerateMockFeedback(problem, lang)
	}

	return feedback
}

func GenerateMockFeedback(problem CFProblem, lang string) AIFeedback {
	if problem.Name == "" {
		return AIFeedback{
			ProblemName: "None",
			Level1:      "No failed problems found to analyze.",
			Level2:      "-",
			Level3:      "-",
			Level4:      "-",
		}
	}

	primaryTag := "Ad-hoc"
	if len(problem.Tags) > 0 {
		primaryTag = problem.Tags[0]
	}

	return AIFeedback{
		ProblemName: fmt.Sprintf("%s (%s)", problem.Name, primaryTag),
		Level1:      fmt.Sprintf("[FALLBACK] Consider the constraints and paradigmatic approach for %s.", primaryTag),
		Level2:      "[API ERROR OR NO KEY] The API request failed or the key is invalid. Please check the terminal logs.",
		Level3:      fmt.Sprintf("[FALLBACK] Check your %s implementation for standard bounds or overflow errors.", lang),
		Level4:      "[FALLBACK] Consider negative elements or N=1.",
		Matrix: []TrainingObjective{
			{Day: "Monday", Focus: primaryTag, Objective: "Upsolving failures", Action: "Solve 3 dynamically equivalent rating problems."},
			{Day: "Tuesday", Focus: "General Practice", Objective: "Daily Consistency", Action: "1 Virtual Contest participation"},
			{Day: "Wednesday", Focus: "Implementation Speed", Objective: "Ad-hoc typing", Action: "Solve 5 problems -200 Elo under time pressure"},
			{Day: "Thursday", Focus: "Core Weaknesses", Objective: "Targeted Algorithmic Theory", Action: "Read CP-Algorithms theory on weakest detected topic."},
			{Day: "Friday", Focus: "Retrospective", Objective: "Consolidate learning", Action: "Upsolve any missed structural problems this week."},
		},
	}
}
