package render

import (
	"bytes"
	"embed"
	"fmt"
	"html"
	"math"
	"strconv"
	"strings"
	"text/template"

	"github.com/sanbei101/blue-card-engine/internal/cardengine"
	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

//go:embed shared/*.tmpl templates/*.tmpl
var tmplFS embed.FS

var tmpl *template.Template

func init() {
	funcMap := template.FuncMap{
		"f":           formatFloat,
		"svgEscape":   html.EscapeString,
		"add":         func(a, b float64) float64 { return a + b },
		"sub":         func(a, b float64) float64 { return a - b },
		"mul":         func(a, b float64) float64 { return a * b },
		"div":         func(a, b float64) float64 { return a / b },
		"min":         math.Min,
		"max":         math.Max,
		"int":         func(v float64) int { return int(v) },
		"float":       func(v int) float64 { return float64(v) },
		"mod":         func(a, b int) int { return a % b },
		"modf":        func(a, b float64) float64 { return math.Mod(a, b) },
		"seq":         func(n int) []int { s := make([]int, n); for i := range s { s[i] = i + 1 }; return s },
		"slice":       func(s []Particle, start, end int) []Particle { return s[start:end] },
		"defaultZero": func(v float64, d float64) float64 { if v == 0 { return d }; return v },
	}

	var err error
	tmpl, err = template.New("cards").Funcs(funcMap).ParseFS(tmplFS, "shared/*.tmpl", "templates/*.tmpl")
	if err != nil {
		panic(fmt.Sprintf("parse svg templates: %v", err))
	}
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(math.Round(v*10)/10, 'f', -1, 64)
}

// RenderData 是传递给 SVG 模板的数据。
type RenderData struct {
	Template         templates.CardTemplate
	Title            string
	Keyword          string
	Lines            []TextLine
	Highlights       []cardengine.HighlightRect
	Particles        []Particle
	HighlightOpacity float64
	SignaturePath    string
}

// TextLine 保存一行文字的 path 与绘制位置。
type TextLine struct {
	Text  string
	Path  string
	Width float64
	X     float64
	Y     float64
}

// Particle 对应 useCardRender 中的装饰粒子。
type Particle struct {
	X        float64
	Y        float64
	Size     float64
	Rotation float64
}

// RenderCard 根据模板、标题和关键词生成 SVG 字符串。
func RenderCard(tpl templates.CardTemplate, title, keyword string, lib *fonts.Library) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = cardengine.DefaultEmptyText
	}

	engineTpl := cardengine.CardTemplate{
		TextBox: cardengine.TextBox{
			Width:         tpl.TextBox.Width,
			Height:        tpl.TextBox.Height,
			X:             tpl.TextBox.X,
			Y:             tpl.TextBox.Y,
			MaxLines:      tpl.TextBox.MaxLines,
			MinFontSize:   tpl.TextBox.MinFontSize,
			MaxFontSize:   tpl.TextBox.MaxFontSize,
			LineHeight:    tpl.TextBox.LineHeight,
			Align:         tpl.TextBox.Align,
			LetterSpacing: 0,
		},
	}

	family := lib.Resolve(tpl.FontFamily)
	measurer := cardengine.NewFontMeasurer(lib.FaceProvider(family), 4096)
	layout := cardengine.BuildTextLayout(title, engineTpl, measurer)
	highlights := cardengine.FindHighlightRects(layout, keyword, engineTpl, measurer)

	lines, err := buildTextLines(lib, tpl, layout)
	if err != nil {
		return "", fmt.Errorf("build text lines: %w", err)
	}

	signaturePath, _, err := TextToPath(lib.SFNT(tpl.FontFamily), "此刻想说", 24)
	if err != nil {
		return "", fmt.Errorf("build signature path: %w", err)
	}

	data := RenderData{
		Template:         tpl,
		Title:            title,
		Keyword:          keyword,
		Lines:            lines,
		Highlights:       highlights,
		Particles:        makeParticles(title, tpl.ID),
		HighlightOpacity: defaultHighlightOpacity(tpl.Kind),
		SignaturePath:    signaturePath,
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tpl.Kind, data); err != nil {
		return "", fmt.Errorf("render template %s: %w", tpl.Kind, err)
	}
	return buf.String(), nil
}

func buildTextLines(lib *fonts.Library, tpl templates.CardTemplate, layout cardengine.TextLayout) ([]TextLine, error) {
	font := lib.SFNT(tpl.FontFamily)
	if font == nil {
		return nil, fmt.Errorf("sfnt font not found for %s", tpl.FontFamily)
	}

	lines := make([]TextLine, 0, len(layout.Lines))
	for i, text := range layout.Lines {
		path, width, err := TextToPath(font, text, layout.FontSize)
		if err != nil {
			return nil, fmt.Errorf("text to path for line %q: %w", text, err)
		}

		y := tpl.TextBox.Y + layout.FontSize + float64(i)*layout.LineHeightPx
		lines = append(lines, TextLine{
			Text:  text,
			Path:  path,
			Width: width,
			X:     lineStartX(tpl.TextBox, width),
			Y:     y,
		})
	}
	return lines, nil
}

func lineStartX(box templates.TextBoxSpec, lineWidth float64) float64 {
	switch box.Align {
	case "center":
		return box.X + (box.Width-lineWidth)/2
	case "right":
		return box.X + box.Width - lineWidth
	default:
		return box.X
	}
}

func defaultHighlightOpacity(kind string) float64 {
	if kind == "night" {
		return 0.9
	}
	return 0
}

func makeParticles(title, id string) []Particle {
	seed := float64(cardengine.HashText(title + ":" + id))
	particles := make([]Particle, 14)
	for i := range particles {
		particles[i] = Particle{
			X:        60 + cardengine.SeededValue(seed, float64(i))*780,
			Y:        90 + cardengine.SeededValue(seed, float64(i+20))*1020,
			Size:     8 + cardengine.SeededValue(seed, float64(i+40))*22,
			Rotation: cardengine.SeededValue(seed, float64(i+60)) * 180,
		}
	}
	return particles
}
