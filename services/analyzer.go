package services

import (
	"fmt"
	"sort"
)

type AnalysisResult struct {
	Profile    Profile             `json:"profile"`
	Weaknesses []Weakness          `json:"weaknesses"`
	Matrix     []TrainingObjective `json:"matrix"`
	AIFeedback AIFeedback          `json:"aiFeedback"`
}

type Profile struct {
	Status  string `json:"status"`
	Cadence string `json:"cadence"`
	Notes   string `json:"notes"`
}

type Weakness struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

type TrainingObjective struct {
	Day       string `json:"day"`
	Focus     string `json:"focus"`
	Objective string `json:"objective"`
	Action    string `json:"action"`
}

type AIFeedback struct {
	ProblemName string              `json:"problemName"`
	Level1      string              `json:"level1"`
	Level2      string              `json:"level2"`
	Level3      string              `json:"level3"`
	Level4      string              `json:"level4"`
	Matrix      []TrainingObjective `json:"matrix"`
}

func AnalyzeUser(handle string) (*AnalysisResult, error) {
	subs, err := FetchUserSubmissions(handle)
	if err != nil {
		return nil, err
	}

	if len(subs) == 0 {
		return nil, fmt.Errorf("no submissions found for user %s", handle)
	}

	// Group submissions by contest chronologically
	var uniqueContests []int
	contestMap := make(map[int][]CFSubmission)
	for _, sub := range subs {
		if len(contestMap[sub.Problem.ContestId]) == 0 {
			uniqueContests = append(uniqueContests, sub.Problem.ContestId)
		}
		contestMap[sub.Problem.ContestId] = append(contestMap[sub.Problem.ContestId], sub)
	}

	var subsToAnalyze []CFSubmission
	
	// Search up to 5 recent contests for one holding algorithmic failures
	for i := 0; i < len(uniqueContests) && i < 5; i++ {
		cid := uniqueContests[i]
		cSubs := contestMap[cid]
		
		hasFailures := false
		for _, sub := range cSubs {
			if sub.Verdict != "OK" {
				hasFailures = true
				break
			}
		}
		
		if hasFailures {
			subsToAnalyze = cSubs
			break
		}
	}

	// Fallback to the most recent contest if literally no failures exist in the previous 5 attempts
	if len(subsToAnalyze) == 0 {
		subsToAnalyze = contestMap[uniqueContests[0]]
	}

	tagFailures := make(map[string]int)
	var latestFailedProblem CFProblem
	var latestFailedLang string
	var latestFailedVerdict string
	var latestFailedSubID int
	var latestFailedContestID int

	var totalTimeDelta int64
	var deltaCount int

	for i, sub := range subsToAnalyze {
		if sub.Verdict != "OK" {
			if latestFailedProblem.Name == "" {
				latestFailedProblem = sub.Problem
				latestFailedLang = sub.ProgrammingLanguage
				latestFailedVerdict = sub.Verdict
				latestFailedSubID = sub.ID
				latestFailedContestID = sub.Problem.ContestId
			}
			for _, tag := range sub.Problem.Tags {
				tagFailures[tag]++
			}
		}

		if i < len(subsToAnalyze)-1 {
			// Only measure cadence between repetitive submissions on the EXACT same problem!
			if subsToAnalyze[i].Problem.Name == subsToAnalyze[i+1].Problem.Name {
				delta := subsToAnalyze[i].CreationTimeSeconds - subsToAnalyze[i+1].CreationTimeSeconds
				if delta > 0 && delta < 3600 {
					totalTimeDelta += delta
					deltaCount++
				}
			}
		}
	}

	var weaknesses []Weakness
	for tag, count := range tagFailures {
		weaknesses = append(weaknesses, Weakness{Tag: tag, Count: count})
	}
	sort.Slice(weaknesses, func(i, j int) bool {
		return weaknesses[i].Count > weaknesses[j].Count
	})

	if len(weaknesses) > 5 {
		weaknesses = weaknesses[:5]
	}

	avgDelta := int64(0)
	if deltaCount > 0 {
		avgDelta = totalTimeDelta / int64(deltaCount)
	}

	profile := Profile{
		Status:  "Active Training",
		Cadence: fmt.Sprintf("%d seconds avg between active attempts", avgDelta),
		Notes:   "Normal cognitive pacing detected. Systematic problem solving.",
	}

	if avgDelta > 0 && avgDelta < 300 {
		profile.Status = "Panic Submitting Detected"
		profile.Notes = "You are rapidly submitting tweaks on the same problem. Slow down and formally dry-run edge cases instead of using the judge to find bugs."
	}

	var code string
	if latestFailedSubID > 0 && latestFailedContestID > 0 {
		// Scrape the actual code that failed from Codeforces
		code, _ = FetchSubmissionCode(latestFailedContestID, latestFailedSubID)
	}

	aiFeedback := GenerateAIFeedback(latestFailedProblem, latestFailedLang, latestFailedVerdict, code, weaknesses)

	result := &AnalysisResult{
		Profile:    profile,
		Weaknesses: weaknesses,
		Matrix:     aiFeedback.Matrix,
		AIFeedback: aiFeedback,
	}

	return result, nil
}
