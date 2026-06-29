package actions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViteHelpers(t *testing.T) {
	fmt.Printf("--- TestViteHelpers debug output ---\n")
	manifest, err := loadViteManifest()
	fmt.Printf("loadViteManifest err: %v\n", err)
	fmt.Printf("manifest keys: ")
	if manifest != nil {
		for k := range manifest {
			fmt.Printf("%s ", k)
		}
	}
	fmt.Printf("\n")

	assetResult := viteAsset("assets/js/main.js")
	cssResult := viteCSS("assets/js/main.js")

	fmt.Printf("viteAsset(\"assets/js/main.js\"): %q\n", assetResult)
	fmt.Printf("viteCSS(\"assets/js/main.js\"): %q\n", cssResult)
	fmt.Printf("------------------------------------\n")

	assert.NotEmpty(t, assetResult)
}
