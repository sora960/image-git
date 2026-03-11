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
	layerName := flag.String("name", "base-layer", "The descriptive name for this layer")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("❌ Error: You must provide a file path.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("🎨 Image-Git: Processing '%s' into repo '%s'...\n", *layerName, *repoName)

	// 1. Store the physical file (CAS Logic)
	hash, err := gitlogic.StoreLayer(*repoName, *filePath)
	if err != nil {
		log.Fatalf("❌ Failed to store layer: %v", err)
	}

	// 2. Load the existing Manifest (The "Tree")
	manifest, err := gitlogic.LoadManifest(*repoName)
	if err != nil {
		log.Fatalf("❌ Failed to load manifest: %v", err)
	}

	// 3. Append the new Layer metadata
	newLayer := gitlogic.Layer{
		Name:    *layerName,
		Hash:    hash,
		Opacity: 1.0, // Default to full opacity
	}
	manifest.Layers = append(manifest.Layers, newLayer)

	// 4. Save the updated Manifest
	err = gitlogic.SaveManifest(*repoName, manifest)
	if err != nil {
		log.Fatalf("❌ Failed to save manifest: %v", err)
	}

	fmt.Printf("✅ Success! Manifest updated with Layer: %s\n", hash[:8])
}