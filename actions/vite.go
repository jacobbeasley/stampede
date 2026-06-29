package actions

import (
	"encoding/json"
	"io/fs"
	"os"

	"buffalo-app/public"
)

// viteManifest holds the parsed Vite manifest file (production only).
type viteManifestEntry struct {
	File string   `json:"file"`
	CSS  []string `json:"css"`
}

var (
	manifestData  map[string]viteManifestEntry
	manifestError error
)

func init() {
	// First try to load from the embedded public FS (used in production binary)
	data, err := fs.ReadFile(public.FS(), "assets/.vite/manifest.json")
	if err != nil {
		// Fallback to local filesystem for development
		data, err = os.ReadFile("public/assets/.vite/manifest.json")
	}

	if err != nil {
		manifestError = err
		return
	}
	manifestError = json.Unmarshal(data, &manifestData)
}

// loadViteManifest reads public/assets/.vite/manifest.json once and caches it.
func loadViteManifest() (map[string]viteManifestEntry, error) {
	return manifestData, manifestError
}

// viteAsset returns the hashed filename for a given source entry (e.g. "assets/js/main.js").
// Falls back to the source path if the manifest is unavailable.
func viteAsset(src string) string {
	m, err := loadViteManifest()
	if err != nil || m == nil {
		return "/" + src
	}
	if entry, ok := m[src]; ok {
		return "/assets/" + entry.File
	}
	return "/" + src
}

// viteCSS returns the hashed CSS filename for a given source entry.
func viteCSS(src string) string {
	m, err := loadViteManifest()
	if err != nil || m == nil {
		return ""
	}
	if entry, ok := m[src]; ok && len(entry.CSS) > 0 {
		return "/assets/" + entry.CSS[0]
	}
	return ""
}
