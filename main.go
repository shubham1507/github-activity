package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type GitHubEvent struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload map[string]interface{} `json:"payload"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}

	username := os.Args[1]
	apiURL := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error fetching data from GitHub:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Failed to fetch activity. Status: %s\n", resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading API response:", err)
		return
	}

	var events []GitHubEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			commits, ok := event.Payload["commits"].([]interface{})
			if ok {
				fmt.Printf("- Pushed %d commits to %s\n", len(commits), event.Repo.Name)
			}
		case "IssuesEvent":
			action, _ := event.Payload["action"].(string)
			fmt.Printf("- %s an issue in %s\n", capitalize(action), event.Repo.Name)
		case "WatchEvent":
			fmt.Printf("- Starred %s\n", event.Repo.Name)
		case "CreateEvent":
			refType, _ := event.Payload["ref_type"].(string)
			fmt.Printf("- Created a new %s in %s\n", refType, event.Repo.Name)
		default:
			fmt.Printf("- %s on %s\n", event.Type, event.Repo.Name)
		}
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
