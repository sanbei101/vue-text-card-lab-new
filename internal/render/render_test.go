package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

func TestRenderCard_SaveSample(t *testing.T) {
	lib, err := fonts.NewLibrary()
	if err != nil {
		t.Fatalf("load fonts: %v", err)
	}

	reg, err := templates.Load()
	if err != nil {
		t.Fatalf("load templates: %v", err)
	}

	tpl, ok := reg.ByKind("mono")
	if !ok {
		t.Fatal("template mono not found")
	}

	svg, err := RenderCard(tpl, "你好,世界", "世界", lib)
	if err != nil {
		t.Fatalf("render card: %v", err)
	}

	outDir := filepath.Join("..", "..", "testdata")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	outPath := filepath.Join(outDir, "sample-mono.svg")
	if err := os.WriteFile(outPath, []byte(svg), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	t.Logf("saved sample SVG to %s", outPath)
}
