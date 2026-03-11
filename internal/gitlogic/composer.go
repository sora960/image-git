package gitlogic

import (
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