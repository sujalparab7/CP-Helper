package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type CFResponse struct {
	Status  string         `json:"status"`
	Result  []CFSubmission `json:"result,omitempty"`
	Comment string         `json:"comment,omitempty"`
}

type CFSubmission struct {
	ID                  int         `json:"id"`
	CreationTimeSeconds int64       `json:"creationTimeSeconds"`
	RelativeTimeSeconds int64       `json:"relativeTimeSeconds"`
	Problem             CFProblem   `json:"problem"`
	Author              CFParty     `json:"author"`
	ProgrammingLanguage string      `json:"programmingLanguage"`
	Verdict             string      `json:"verdict"`
	Testset             string      `json:"testset"`
	PassedTestCount     int         `json:"passedTestCount"`
	TimeConsumedMillis  int         `json:"timeConsumedMillis"`
	MemoryConsumedBytes int         `json:"memoryConsumedBytes"`
}

type CFProblem struct {
	ContestId int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Points    float64  `json:"points,omitempty"`
	Rating    int      `json:"rating,omitempty"`
	Tags      []string `json:"tags"`
}

type CFParty struct {
	ContestId        int        `json:"contestId,omitempty"`
	Members          []CFMember `json:"members"`
	ParticipantType  string     `json:"participantType"`
	TeamId           int        `json:"teamId,omitempty"`
	TeamName         string     `json:"teamName,omitempty"`
	StartTimeSeconds int64      `json:"startTimeSeconds,omitempty"`
}

type CFMember struct {
	Handle string `json:"handle"`
	Name   string `json:"name,omitempty"`
}

func FetchUserSubmissions(handle string) ([]CFSubmission, error) {
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=200", handle)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cfResp CFResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return nil, fmt.Errorf("failed to decode Codeforces API response: %w", err)
	}

	if cfResp.Status != "OK" {
		return nil, fmt.Errorf("Codeforces API error: %s", cfResp.Comment)
	}

	return cfResp.Result, nil
}

func FetchSubmissionCode(contestId int, submissionId int) (string, error) {
	url := fmt.Sprintf("https://codeforces.com/contest/%d/submission/%d", contestId, submissionId)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyStr := string(bodyBytes)

	re := regexp.MustCompile(`(?s)<pre id="program-source-text"[^>]*>(.*?)</pre>`)
	matches := re.FindStringSubmatch(bodyStr)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find source code in response")
	}

	code := matches[1]
	code = strings.ReplaceAll(code, "&lt;", "<")
	code = strings.ReplaceAll(code, "&gt;", ">")
	code = strings.ReplaceAll(code, "&quot;", "\"")
	code = strings.ReplaceAll(code, "&amp;", "&")
	
	if len(code) > 20000 {
		code = code[:20000]
	}
	
	return code, nil
}
