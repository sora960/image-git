package gitlogic

import (
    "image"
	"image/color"
    "image/draw"
    "image/png"
    "os"
    "sort"
)


// Dynamic Canvas Size Logic
var canvas *image.RGBA

func RenderRepository(repoName string) error {
    manifest, err := LoadManifest(repoName)
    if err != nil {
        return err
    }

    if len(manifest.Layers) == 0 {
        return nil // Nothing to render
    }

    // 1. Sort layers by Z-Index (Back to Front)
    sort.Slice(manifest.Layers, func(i, j int) bool {
        return manifest.Layers[i].ZIndex < manifest.Layers[j].ZIndex
    })

    var canvas *image.RGBA

    for i, layer := range manifest.Layers {
        // 2. Open the layer's blob
        file, err := os.Open("data/repositories/" + repoName + "/objects/" + layer.Hash + ".png")
        if err != nil {
            continue 
        }
        
        img, _, err := image.Decode(file)
        file.Close()
        if err != nil {
            continue
        }

        // 3. Dynamic Initialization
        // Set canvas size based on the first (bottom-most) layer
        if i == 0 {
            canvas = image.NewRGBA(img.Bounds())
            // Optional: Fill with transparent background if needed
            draw.Draw(canvas, img.Bounds(), image.Transparent, image.Point{}, draw.Src)
        }

        // 4. Composite Layer with custom Opacity
        // Create a mask with the specific alpha level from our manifest
        // We multiply by 255 because Go's Alpha uses an 8-bit range (0-255)
        mask := image.NewUniform(color.Alpha{uint8(layer.Opacity * 255)})

        // Use DrawMask instead of Draw
        draw.DrawMask(
            canvas,          // Target
            canvas.Bounds(), // Rect
            img,             // Source image
            image.Point{},   // Source point
            mask,            // THE FIX: Applying the opacity value
            image.Point{},   // Mask point
            draw.Over,       // Porter-Duff operator
        )
    }

    // 5. Save the output
    out, err := os.Create("data/repositories/" + repoName + "/preview.png")
    if err != nil {
        return err
    }
    defer out.Close()

    return png.Encode(out, canvas)
}