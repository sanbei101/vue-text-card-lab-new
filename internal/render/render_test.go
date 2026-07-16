package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

func setup(tb testing.TB) (*fonts.Library, *templates.Registry) {
	tb.Helper()
	lib, err := fonts.NewLibrary()
	if err != nil {
		tb.Fatalf("load fonts: %v", err)
	}

	reg, err := templates.Load()
	if err != nil {
		tb.Fatalf("load templates: %v", err)
	}

	return lib, reg
}

func TestRenderCard_SaveSample(t *testing.T) {
	lib, reg := setup(t)

	tpl, ok := reg.ByKind("mono")
	if !ok {
		t.Fatal("template mono not found")
	}

	svg, err := Card(&tpl, "遥看抚仙湖吃鱼,美汤底和铜锅洋芋饭绝!", "世界", lib)
	if err != nil {
		t.Fatalf("render card: %v", err)
	}

	outDir := filepath.Join("..", "..", "testdata")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	outPath := filepath.Join(outDir, "sample-mono.svg")
	if err := os.WriteFile(outPath, []byte(svg), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	t.Logf("saved sample SVG to %s", outPath)
}

func BenchmarkRenderCard(b *testing.B) {
	lib, reg := setup(b)
	tpl, ok := reg.ByKind("mono")
	if !ok {
		b.Fatal("template mono not found")
	}
	for b.Loop() {
		Card(&tpl, "你好,世界", "世界", lib)
	}
}

func TestRenderCard_AllTemplates(t *testing.T) {
	lib, reg := setup(t)

	outDir := filepath.Join("..", "..", "testdata", "all-templates")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	titles := map[string]string{
		"question":  "为什么我们总是怀念夏天？",
		"memo":      "今日灵感：把复杂的问题拆成小块。",
		"editorial": "设计是一种沟通方式。",
		"lavender":  "保持温柔，也保持锋利。",
		"notebook":  "生活不是等待暴风雨过去，而是学会在雨中跳舞。",
		"stamp":     "每天都要认真生活",
		"grid":      "清晰的目标胜过盲目的努力。",
		"night":     "星星发亮是为了让每个人找到属于自己的光。",
		"burst":     "充满能量的早晨！",
		"paper":     "读书是为了遇见更好的自己。",
		"bubble":    "今天也要开心鸭",
		"mono":      "少即是多。",
	}

	for i := range reg.All() {
		tpl := &reg.All()[i]
		title := titles[tpl.Kind]
		if title == "" {
			title = "你好，世界"
		}
		keyword := ""
		svg, err := Card(tpl, title, keyword, lib)
		if err != nil {
			t.Fatalf("render template %s (%s): %v", tpl.Kind, tpl.ID, err)
		}

		outPath := filepath.Join(outDir, tpl.Kind+".svg")
		if err := os.WriteFile(outPath, []byte(svg), 0o644); err != nil {
			t.Fatalf("write file %s: %v", outPath, err)
		}
		t.Logf("saved %s", outPath)
	}
}
