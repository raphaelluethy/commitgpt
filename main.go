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
	unstagedChanges := getCommandOutput("git", "diff")
	if unstagedChanges == "" {
		fmt.Println("No unstaged changes found.")
		for _, arg := range os.Args[1:] {
			if arg == "--help" || arg == "-h" {
				fmt.Println("Usage: commitgpt [options]")
				fmt.Println("Options:")
				fmt.Println("  --help, -h  Show this help message and exit.")
				os.Exit(0)
			}
		}
		return
	}

	changesOverview := getCommandOutput("git", "diff", "--stat")

	content := fmt.Sprintf("Detailed Changes:\n%s\n\nChanges Overview:\n%s", unstagedChanges, changesOverview)

	summary := getAnthropicSummary(content)

	createGitCommit(summary)

	fmt.Printf("Created commit: %s\n", summary)
}

// getCommandOutput executes a command with the given name and arguments,
// and returns its output as a trimmed string. If the command fails,
// it prints the error and exits the program.
func getCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(output))
}

// getAnthropicSummary sends the provided content to Anthropic's API
// to generate a concise one-line summary of git changes. It returns
// the generated summary or an error message if the API call fails.
func getAnthropicSummary(content string) string {
	prompt := fmt.Sprintf("Summarize the following Git changes:\n\n%s\n\nProvide a concise one-line summary of the changes, like the following: `fix: fixed an issue where a memory leak was happening` or `feat: added the abillity to take screenshots`. ONLY RETURN ONE LINE. Here is the content:", content)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "claude-3-5-sonnet-latest",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 4096,
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

// createGitCommit stages all changes and creates a new git commit
// with the provided message. If either operation fails, it prints
// the error and exits the program.
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
