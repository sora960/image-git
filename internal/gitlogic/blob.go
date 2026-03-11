package gitlogic

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// StoreLayer calculates a SHA-256 hash of a file and moves it to an immutable object store.
func StoreLayer(repoName string, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Hash the content
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashString := fmt.Sprintf("%x", hash.Sum(nil))

	// Define destination: data/repositories/<repo>/objects/<hash>.png
	dstDir := filepath.Join("data", "repositories", repoName, "objects")
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return "", err
	}

	dstPath := filepath.Join(dstDir, hashString+".png")

	// Skip if already exists (Deduplication)
	if _, err := os.Stat(dstPath); err == nil {
		return hashString, nil
	}

	// Copy to object store
	file.Seek(0, 0)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()
	io.Copy(dstFile, file)

	return hashString, nil
}
