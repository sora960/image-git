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
	targetFrame := flag.Int("frame", 0, "The specific frame index to render")
	startFrame := flag.Int("start", 0, "Starting frame for a new layer")
	endFrame := flag.Int("end", 999, "Ending frame for a new layer")
	flag.Parse()

// 2. Handle Compositing
    if *doComposite {
        // Scenario A: Render a range of frames (Sequence)
        if *endFrame > *startFrame && *targetFrame == 0 { 
            err := gitlogic.CompositeSequence(*repoName, *startFrame, *endFrame)
            if err != nil {
                log.Fatalf("❌ Animation render failed: %v", err)
            }
            fmt.Println("✅ Animation sequence rendered successfully.")
            return
        }

        // Scenario B: Render a single specific frame
        fmt.Printf("🎬 Image-Git: Rendering Frame %d for repo '%s'...\n", *targetFrame, *repoName)
        err := gitlogic.CompositeFrame(*repoName, *targetFrame)
        if err != nil {
            log.Fatalf("❌ Single frame render failed: %v", err)
        }
        fmt.Printf("✅ Success! Created frame_%04d.png\n", *targetFrame)
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
		StartFrame: *startFrame,
		EndFrame:   *endFrame,
	}
	manifest.Layers = append(manifest.Layers, newLayer)

	err = gitlogic.SaveManifest(*repoName, manifest)
	if err != nil {
		log.Fatalf("❌ Failed to save manifest: %v", err)
	}

	fmt.Printf("✅ Success! Manifest updated with Layer: %s\n", hash[:8])
}