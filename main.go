package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var anthropicAPIKey = os.Getenv("ANTHROPIC_API_KEY")

const anthropicAPIURL = "https://api.anthropic.com/v1/messages"

func main() {
	// Check for unstaged changes
	unstagedChanges := getCommandOutput("git", "diff")
	if unstagedChanges == "" {
		fmt.Println("No unstaged changes found.")
		return
	}

	// Get changes overview
	changesOverview := getCommandOutput("git", "diff", "--stat")

	// Prepare content for summarization
	content := fmt.Sprintf("Detailed Changes:\n%s\n\nChanges Overview:\n%s", unstagedChanges, changesOverview)

	// Get summary from Anthropic API
	summary := getAnthropicSummary(content)

	// Determine commit type
	commitType := determineCommitType(summary)

	// Create commit message
	commitMessage := fmt.Sprintf("%s: %s", commitType, summary)

	// Create git commit
	createGitCommit(commitMessage)

	fmt.Printf("Created commit: %s\n", commitMessage)
}

func getCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(output))
}

func getAnthropicSummary(content string) string {
	prompt := fmt.Sprintf("Summarize the following Git changes:\n\n%s\n\nProvide a concise one-line summary of the changes, like the following: `fix: fixed an issue where a memory leak was happening` or `feat: added the abillity to take screenshots`. Here is the content:", content)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      "claude-3-5-sonnet-20240620",
		"max_tokens": 8096,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, _ := http.NewRequest("POST", anthropicAPIURL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", anthropicAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error calling Anthropic API: %v\n", err)
		return "Unable to generate summary"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
		if text, ok := content[0].(map[string]interface{})["text"].(string); ok {
			return text
		}
	}

	return "Unable to generate summary"
}

func determineCommitType(summary string) string {
	summary = strings.ToLower(summary)
	if strings.Contains(summary, "fix") || strings.Contains(summary, "bug") {
		return "fix"
	} else if strings.Contains(summary, "feat") || strings.Contains(summary, "feature") {
		return "feat"
	} else if strings.Contains(summary, "docs") || strings.Contains(summary, "documentation") {
		return "docs"
	} else if strings.Contains(summary, "style") {
		return "style"
	} else if strings.Contains(summary, "refactor") {
		return "refactor"
	} else if strings.Contains(summary, "test") {
		return "test"
	} else if strings.Contains(summary, "chore") {
		return "chore"
	}
	return "chore" // default to chore if no specific type is detected
}

func createGitCommit(message string) {
	cmd := exec.Command("git", "add", ".")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error staging changes: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("git", "commit", "-m", message)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error creating commit: %v\n", err)
		os.Exit(1)
	}
}
