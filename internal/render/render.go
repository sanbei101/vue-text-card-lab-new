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
	Layout           cardengine.TextLayout
	Highlights       []cardengine.HighlightRect
	TextAnchor       string
	TextX            float64
	Particles        []Particle
	HighlightOpacity float64
}

// Particle 对应 useCardRender 中的装饰粒子。
type Particle struct {
	X        float64
	Y        float64
	Size     float64
	Rotation float64
}

// RenderCard 根据模板、标题和关键词生成 SVG 字符串。
func RenderCard(tpl templates.CardTemplate, title, keyword string, measurer cardengine.Measurer) (string, error) {
	title = strings.TrimSpace(title)

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

	layout := cardengine.BuildTextLayout(title, engineTpl, measurer)
	highlights := cardengine.FindHighlightRects(layout, keyword, engineTpl, measurer)

	data := RenderData{
		Template:         tpl,
		Title:            title,
		Keyword:          keyword,
		Layout:           layout,
		Highlights:       highlights,
		TextAnchor:       textAnchor(tpl.TextBox.Align),
		TextX:            textX(tpl.TextBox),
		Particles:        makeParticles(title, tpl.ID),
		HighlightOpacity: defaultHighlightOpacity(tpl.Kind),
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tpl.Kind, data); err != nil {
		return "", fmt.Errorf("render template %s: %w", tpl.Kind, err)
	}
	return buf.String(), nil
}

func defaultHighlightOpacity(kind string) float64 {
	if kind == "night" {
		return 0.9
	}
	return 0
}

func textAnchor(align string) string {
	switch align {
	case "center":
		return "middle"
	case "right":
		return "end"
	default:
		return "start"
	}
}

func textX(box templates.TextBoxSpec) float64 {
	switch box.Align {
	case "center":
		return box.X + box.Width/2
	case "right":
		return box.X + box.Width
	default:
		return box.X
	}
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
