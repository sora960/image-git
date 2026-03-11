package gitlogic

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// CreateTestAssets generates two real PNGs for compositing tests
func CreateTestAssets() error {
	// 1. Create a Blue Background (500x500)
	bg := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(bg, bg.Bounds(), &image.Uniform{color.RGBA{0, 0, 255, 255}}, image.Point{}, draw.Src)
	
	f1, _ := os.Create("bg_test.png")
	defer f1.Close()
	png.Encode(f1, bg)

	// 2. Create a Red Square "Character" (200x200) on a transparent canvas
	fg := image.NewRGBA(image.Rect(0, 0, 500, 500))
	// Draw a red square in the middle
	redRect := image.Rect(150, 150, 350, 350)
	draw.Draw(fg, redRect, &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src)

	f2, _ := os.Create("fg_test.png")
	defer f2.Close()
	return png.Encode(f2, fg)
}