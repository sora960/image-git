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
	targetFrame := flag.Int("frame", -1, "The specific frame index to render")
	startFrame := flag.Int("start", 0, "Starting frame for a new layer")
	endFrame := flag.Int("end", 999, "Ending frame for a new layer")
	zIndex := flag.Int("z", 0, "Z-Index for the layer (lower is further back)")
	opacity := flag.Float64("opacity", 1.0, "Opacity of the layer (0.0 to 1.0)")
    deleteLayer := flag.String("delete", "", "Name of the layer to remove from the manifest")
	showStatus := flag.Bool("status", false, "Show the current layer stack for the repo")
	flag.Parse()


// 1. Handle Status Table
if *showStatus {
	m, err := gitlogic.LoadManifest(*repoName)
	if err != nil {
		log.Fatalf("❌ Failed to load manifest: %v", err)
	}
	fmt.Printf("\n📂 Repository: %s\n", *repoName)
	fmt.Printf("%-12s | %-8s | %-3s | %-10s | %-7s\n", "NAME", "HASH", "Z", "RANGE", "ALPHA")
	fmt.Println("------------------------------------------------------------")
	for _, l := range m.Layers {
		fmt.Printf("%-12s | %-8s | %-3d | %d-%-7d | %.2f\n", 
			l.Name, l.Hash[:8], l.ZIndex, l.StartFrame, l.EndFrame, l.Opacity)
	}
	return
}

// 2. Handle Deletion
if *deleteLayer != "" {
	err := gitlogic.RemoveLayer(*repoName, *deleteLayer)
	if err != nil {
		log.Fatalf("❌ Delete failed: %v", err)
	}
	fmt.Printf("🗑️  Layer '%s' removed from %s manifest.\n", *deleteLayer, *repoName)
	return
}


// 2. Handle Compositing
if *doComposite {
    // If user provided a specific frame (not -1), render just that one
    if *targetFrame != -1 {
        fmt.Printf("🎬 Image-Git: Rendering Frame %d for repo '%s'...\n", *targetFrame, *repoName)
        err := gitlogic.CompositeFrame(*repoName, *targetFrame)
        if err != nil {
            log.Fatalf("❌ Single frame render failed: %v", err)
        }
        fmt.Printf("✅ Success! Created frame_%04d.png\n", *targetFrame)
        return
    }

    // Otherwise, if they gave a range (and no specific frame), do the sequence
    if *endFrame > *startFrame {
        err := gitlogic.CompositeSequence(*repoName, *startFrame, *endFrame)
        if err != nil {
            log.Fatalf("❌ Animation render failed: %v", err)
        }
        fmt.Println("✅ Animation sequence rendered successfully.")
        return
    }
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
		Opacity: *opacity,
		StartFrame: *startFrame,
		EndFrame:   *endFrame,
		ZIndex:     *zIndex,
	}
	manifest.Layers = append(manifest.Layers, newLayer)

	err = gitlogic.SaveManifest(*repoName, manifest)
	if err != nil {
		log.Fatalf("❌ Failed to save manifest: %v", err)
	}

	fmt.Printf("✅ Success! Manifest updated with Layer: %s\n", hash[:8])
}