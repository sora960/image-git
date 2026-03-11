package gitlogic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func StoreLayer(repoName string, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file err: %w", err)
	}
	defer file.Close()
	return StoreStream(repoName, file)
}

func StoreStream(repoName string, reader io.Reader) (string, error) {

	// 1. Create temporary file in our LOCAL project temp folder
    // Change "" to "./temp" here:
    tempFile, err := os.CreateTemp("./temp", "image-git-upload-*") 
    if err != nil {
        return "", fmt.Errorf("temp file err: %w", err)
    }

	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	// 2. MultiWriter: Hash and Write at once
	hasher := sha256.New()
	mw := io.MultiWriter(tempFile, hasher)

	if _, err := io.Copy(mw, reader); err != nil {
		tempFile.Close()
		return "", fmt.Errorf("streaming copy err: %w", err)
	}
	tempFile.Close() 

	hash := hex.EncodeToString(hasher.Sum(nil))

	// 3. Ensure the destination directory exists
	objDir := filepath.Join("data", "repositories", repoName, "objects")
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir err: %w", err)
	}

	finalPath := filepath.Join(objDir, hash+".png")
	
	// 4. Move to permanent storage
	if err := os.Rename(tempPath, finalPath); err != nil {
		return "", fmt.Errorf("rename/move err: %w", err)
	}

	return hash, nil
}