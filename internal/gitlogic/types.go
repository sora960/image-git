package gitlogic

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Issue #19: Define Commit Structure
type Commit struct {
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Timestamp int64     `json:"timestamp"`
	Manifest  Manifest  `json:"manifest"`
}

// Issue #20: Implement Content-Addressable Saving
func SaveCommit(repoName string, message string) (string, error) {
	manifest, err := LoadManifest(repoName)
	if err != nil {
		return "", err
	}

	commit := Commit{
		Message:   message,
		Timestamp: time.Now().Unix(),
		Manifest:  manifest, // Removed the '*' because your LoadManifest returns the struct
	}

	// Generate Hash
	data, _ := json.Marshal(commit)
	hashBytes := sha256.Sum256(data)
	commit.Hash = hex.EncodeToString(hashBytes[:])

	// Ensure directory exists
	historyDir := filepath.Join("data", "repositories", repoName, "history")
	os.MkdirAll(historyDir, 0755)

	// Save as [HASH].json
	commitPath := filepath.Join(historyDir, commit.Hash+".json")
	commitData, _ := json.MarshalIndent(commit, "", "  ")
	
	err = os.WriteFile(commitPath, commitData, 0644)
	return commit.Hash, err
}