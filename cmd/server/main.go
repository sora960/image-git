package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sora960/image-git/internal/gitlogic"
)

func main() {
	// 1. Declare the pointers once using := 
	filePath := flag.String("file", "", "Path to the image layer to commit")
	repoName := flag.String("repo", "test-repo", "Name of the target repository")
	layerName := flag.String("name", "base-layer", "The descriptive name for this layer")
	doComposite := flag.Bool("composite", false, "Generate a preview.png by merging all layers")
	
	flag.Parse()

	// 2. Handle Compositing
	if *doComposite {
		fmt.Printf("🎨 Image-Git: Merging layers for repo '%s'...\n", *repoName)
		err := gitlogic.CompositeLayers(*repoName)
		if err != nil {
			log.Fatalf("❌ Compositing failed: %v", err)
		}
		fmt.Printf("✅ Success! Created: data/repositories/%s/preview.png\n", *repoName)
		return
	}

	// 3. Handle Storing
	if *filePath == "" {
		fmt.Println("❌ Error: You must provide a file path or use --composite.")
		flag.Usage()
		os.Exit(1)
	}

	hash, err := gitlogic.StoreLayer(*repoName, *filePath)
	if err != nil {
		log.Fatalf("❌ Failed to store layer: %v", err)
	}

	manifest, err := gitlogic.LoadManifest(*repoName)
	if err != nil {
		log.Fatalf("❌ Failed to load manifest: %v", err)
	}

	newLayer := gitlogic.Layer{
		Name:    *layerName,
		Hash:    hash,
		Opacity: 1.0,
	}
	manifest.Layers = append(manifest.Layers, newLayer)

	err = gitlogic.SaveManifest(*repoName, manifest)
	if err != nil {
		log.Fatalf("❌ Failed to save manifest: %v", err)
	}

	fmt.Printf("✅ Success! Manifest updated with Layer: %s\n", hash[:8])
}