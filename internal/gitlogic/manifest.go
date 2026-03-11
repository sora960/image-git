package gitlogic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Layer represents a single visual element in the stack
type Layer struct {
	Name    string  `json:"name"`
	Hash    string  `json:"hash"`
	Opacity float64 `json:"opacity"`
	StartFrame int     `json:"start_frame"` // When it appears
	EndFrame   int     `json:"end_frame"`   // When it disappears
	ZIndex     int     `json:"z_index"`
}

// Manifest represents the state of an entire artwork project
type Manifest struct {
	Layers []Layer `json:"layers"`
}

// LoadManifest reads the JSON from the repo or returns a fresh one if missing
func LoadManifest(repoName string) (Manifest, error) {
	path := filepath.Join("data", "repositories", repoName, "manifest.json")
	var m Manifest

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty manifest if it's a new project
			return Manifest{Layers: []Layer{}}, nil
		}
		return m, err
	}

	err = json.Unmarshal(data, &m)
	return m, err
}

// SaveManifest writes the current state back to manifest.json
func SaveManifest(repoName string, m Manifest) error {
	path := filepath.Join("data", "repositories", repoName, "manifest.json")
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// RemoveLayer searches for a layer by name and removes it from the manifest
func RemoveLayer(repoName string, layerName string) error {
    m, err := LoadManifest(repoName)
    if err != nil {
        return err
    }

    var updatedLayers []Layer
    found := false

    for _, l := range m.Layers {
        if l.Name == layerName {
            found = true
            continue // Skip the layer we want to delete
        }
        updatedLayers = append(updatedLayers, l)
    }

    if !found {
        return fmt.Errorf("layer '%s' not found", layerName)
    }

    m.Layers = updatedLayers
    return SaveManifest(repoName, m)
}

// PruneObjects removes files from the objects folder that are not referenced in the manifest
func PruneObjects(repoName string) (int, error) {
    m, err := LoadManifest(repoName)
    if err != nil {
        return 0, err
    }

    // Build a map of referenced hashes
    referenced := make(map[string]bool)
    for _, l := range m.Layers {
        referenced[l.Hash+".png"] = true
    }

    objDir := filepath.Join("data", "repositories", repoName, "objects")
    files, err := os.ReadDir(objDir)
    if err != nil {
        return 0, err
    }

    count := 0
    for _, f := range files {
        if !referenced[f.Name()] {
            err := os.Remove(filepath.Join(objDir, f.Name()))
            if err != nil {
                fmt.Printf("⚠️  Failed to delete %s: %v\n", f.Name(), err)
                continue
            }
            count++
        }
    }
    return count, nil
}

