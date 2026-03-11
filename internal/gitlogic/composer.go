package gitlogic

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

// CompositeLayers merges all layers in a manifest into one preview image
func CompositeLayers(repoName string) error {
	m, err := LoadManifest(repoName)
	if err != nil {
		return err
	}

	// Create a transparent canvas
	canvas := image.NewRGBA(image.Rect(0, 0, 500, 500))

	for _, layer := range m.Layers {
		path := filepath.Join("data", "repositories", repoName, "objects", layer.Hash+".png")
		f, err := os.Open(path)
		if err != nil {
			continue
		}

		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			continue
		}

		// Draw layer onto canvas using the "Over" operator (Standard transparency blending)
		draw.Draw(canvas, canvas.Bounds(), img, image.Point{}, draw.Over)
	}

	// Save the final result
	outPath := filepath.Join("data", "repositories", repoName, "preview.png")
	out, _ := os.Create(outPath)
	defer out.Close()
	return png.Encode(out, canvas)
}

// CompositeFrame merges only the layers active at a specific frame index
func CompositeFrame(repoName string, frame int) error {
	m, err := LoadManifest(repoName)
	if err != nil {
		return err
	}

	canvas := image.NewRGBA(image.Rect(0, 0, 500, 500))

	for _, layer := range m.Layers {
		// Only draw if the current frame is within the layer's lifespan
		if frame >= layer.StartFrame && frame <= layer.EndFrame {
			path := filepath.Join("data", "repositories", repoName, "objects", layer.Hash+".png")
			f, err := os.Open(path)
			if err != nil {
				continue
			}

			img, err := png.Decode(f)
			f.Close()
			if err != nil {
				continue
			}

			draw.Draw(canvas, canvas.Bounds(), img, image.Point{}, draw.Over)
		}
	}

	// Save as a frame-specific preview
	outPath := filepath.Join("data", "repositories", repoName, fmt.Sprintf("frame_%04d.png", frame))
	out, _ := os.Create(outPath)
	defer out.Close()
	return png.Encode(out, canvas)
}