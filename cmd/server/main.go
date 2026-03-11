package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sora960/image-git/internal/gitlogic"
)

func main() {
	// Define flags
	filePath := flag.String("file", "", "Path to the image layer to commit")
	repoName := flag.String("repo", "test-repo", "Name of the target repository")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("❌ Error: You must provide a file path.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("🎨 Image-Git: Committing layer '%s' to '%s'...\n", *filePath, *repoName)

	// Call our core CAS logic
	hash, err := gitlogic.StoreLayer(*repoName, *filePath)
	if err != nil {
		log.Fatalf("❌ Failed to store layer: %v", err)
	}

	fmt.Printf("✅ Success! Layer stored.\n")
	fmt.Printf("Hash: %s\n", hash)
}