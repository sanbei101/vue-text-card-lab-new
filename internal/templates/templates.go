package templates

import (
	"fmt"

	"github.com/sanbei101/blue-card-engine/internal/fonts"
)

// TextBoxSpec 对应前端 TextBoxSpec。
type TextBoxSpec struct {
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	MinFontSize float64 `json:"minFontSize"`
	MaxFontSize float64 `json:"maxFontSize"`
	MaxLines    int     `json:"maxLines"`
	LineHeight  float64 `json:"lineHeight"`
	Align       string  `json:"align"`
}

// CardTemplate 对应前端 CardTemplate。
type CardTemplate struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Kind       string           `json:"kind"`
	Background string           `json:"background"`
	Foreground string           `json:"foreground"`
	Accent     string           `json:"accent"`
	Muted      string           `json:"muted"`
	FontFamily fonts.FontFamily `json:"fontFamily"`
	FontWeight int              `json:"fontWeight"`
	TextBox    TextBoxSpec      `json:"textBox"`
	Radius     float64          `json:"radius"`
	Frame      string           `json:"frame,omitempty"`
}

var builtinTemplates = []CardTemplate{
	{
		ID:         "question-blue",
		Name:       "清透问号",
		Kind:       "question",
		Background: "#eef9fd",
		Foreground: "#093c4a",
		Accent:     "#ffe23b",
		Muted:      "#9fd4e6",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 800,
		Radius:     54,
		TextBox: TextBoxSpec{
			X: 96, Y: 290, Width: 708, Height: 620,
			MinFontSize: 44, MaxFontSize: 88, MaxLines: 6,
			LineHeight: 1.26, Align: "left",
		},
	},
	{
		ID:         "yellow-memo",
		Name:       "便签日记",
		Kind:       "memo",
		Background: "#fff4a8",
		Foreground: "#202018",
		Accent:     "#ff6b35",
		Muted:      "#d8b94f",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 750,
		Radius:     28,
		TextBox: TextBoxSpec{
			X: 90, Y: 245, Width: 720, Height: 620,
			MinFontSize: 44, MaxFontSize: 82, MaxLines: 7,
			LineHeight: 1.26, Align: "left",
		},
	},
	{
		ID:         "editorial-black",
		Name:       "黑白杂志",
		Kind:       "editorial",
		Background: "#f7f5ef",
		Foreground: "#101010",
		Accent:     "#ef4136",
		Muted:      "#b6b1a8",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 800,
		Radius:     0,
		TextBox: TextBoxSpec{
			X: 82, Y: 280, Width: 736, Height: 630,
			MinFontSize: 42, MaxFontSize: 84, MaxLines: 7,
			LineHeight: 1.2, Align: "left",
		},
	},
	{
		ID:         "lavender-soft",
		Name:       "柔雾紫",
		Kind:       "lavender",
		Background: "#eee8ff",
		Foreground: "#30245f",
		Accent:     "#8f68ff",
		Muted:      "#c7b9f5",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 760,
		Radius:     72,
		TextBox: TextBoxSpec{
			X: 110, Y: 300, Width: 680, Height: 550,
			MinFontSize: 44, MaxFontSize: 80, MaxLines: 7,
			LineHeight: 1.26, Align: "center",
		},
	},
	{
		ID:         "green-notebook",
		Name:       "绿格手账",
		Kind:       "notebook",
		Background: "#e9f3df",
		Foreground: "#23452e",
		Accent:     "#ff7a59",
		Muted:      "#a8c69c",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 760,
		Radius:     42,
		TextBox: TextBoxSpec{
			X: 115, Y: 280, Width: 660, Height: 590,
			MinFontSize: 40, MaxFontSize: 78, MaxLines: 8,
			LineHeight: 1.26, Align: "left",
		},
	},
	{
		ID:         "red-stamp",
		Name:       "红印宣言",
		Kind:       "stamp",
		Background: "#fff0e7",
		Foreground: "#7d201b",
		Accent:     "#e83a2f",
		Muted:      "#e0a298",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 900,
		Radius:     24,
		TextBox: TextBoxSpec{
			X: 92, Y: 255, Width: 716, Height: 610,
			MinFontSize: 44, MaxFontSize: 86, MaxLines: 7,
			LineHeight: 1.18, Align: "center",
		},
	},
	{
		ID:         "blue-grid",
		Name:       "蓝图网格",
		Kind:       "grid",
		Background: "#eaf2ff",
		Foreground: "#173f7a",
		Accent:     "#3978ff",
		Muted:      "#b8cff5",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 780,
		Radius:     40,
		TextBox: TextBoxSpec{
			X: 100, Y: 300, Width: 700, Height: 580,
			MinFontSize: 44, MaxFontSize: 80, MaxLines: 7,
			LineHeight: 1.26, Align: "left",
		},
	},
	{
		ID:         "midnight-stars",
		Name:       "深夜星光",
		Kind:       "night",
		Background: "#11162a",
		Foreground: "#f6f1d8",
		Accent:     "#ffcc4d",
		Muted:      "#6a7399",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 740,
		Radius:     58,
		TextBox: TextBoxSpec{
			X: 110, Y: 315, Width: 680, Height: 560,
			MinFontSize: 44, MaxFontSize: 78, MaxLines: 7,
			LineHeight: 1.26, Align: "center",
		},
	},
	{
		ID:         "orange-burst",
		Name:       "橙色爆炸",
		Kind:       "burst",
		Background: "#ff8a42",
		Foreground: "#2b1710",
		Accent:     "#fff3cf",
		Muted:      "#ffb07e",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 900,
		Radius:     46,
		TextBox: TextBoxSpec{
			X: 98, Y: 285, Width: 704, Height: 590,
			MinFontSize: 44, MaxFontSize: 86, MaxLines: 7,
			LineHeight: 1.16, Align: "center",
		},
	},
	{
		ID:         "beige-paper",
		Name:       "纸张摘录",
		Kind:       "paper",
		Background: "#f0eadf",
		Foreground: "#3a3128",
		Accent:     "#315d4f",
		Muted:      "#c9bead",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 700,
		Radius:     18,
		TextBox: TextBoxSpec{
			X: 120, Y: 290, Width: 660, Height: 600,
			MinFontSize: 40, MaxFontSize: 76, MaxLines: 7,
			LineHeight: 1.34, Align: "left",
		},
	},
	{
		ID:         "pink-bubble",
		Name:       "粉色气泡",
		Kind:       "bubble",
		Background: "#ffe9f1",
		Foreground: "#61253c",
		Accent:     "#ff5790",
		Muted:      "#efb7cb",
		FontFamily: fonts.FontHarmonyOS,
		FontWeight: 790,
		Radius:     76,
		TextBox: TextBoxSpec{
			X: 104, Y: 315, Width: 692, Height: 540,
			MinFontSize: 44, MaxFontSize: 80, MaxLines: 7,
			LineHeight: 1.26, Align: "center",
		},
	},
	{
		ID:         "mono-frame",
		Name:       "极简边框",
		Kind:       "mono",
		Background: "#ffffff",
		Foreground: "#181818",
		Accent:     "#181818",
		Muted:      "#d9d9d9",
		FontFamily: fonts.FontSong,
		FontWeight: 700,
		Radius:     12,
		Frame:      "#181818",
		TextBox: TextBoxSpec{
			X: 110, Y: 300, Width: 680, Height: 560,
			MinFontSize: 44, MaxFontSize: 78, MaxLines: 7,
			LineHeight: 1.26, Align: "center",
		},
	},
}

// Registry 保存所有已加载模板。
type Registry struct {
	byID   map[string]CardTemplate
	byKind map[string]CardTemplate
	list   []CardTemplate
}

// Load 返回内嵌的模板注册表。
func Load() (*Registry, error) {
	reg := &Registry{
		byID:   make(map[string]CardTemplate),
		byKind: make(map[string]CardTemplate),
		list:   make([]CardTemplate, 0, len(builtinTemplates)),
	}

	for i := range builtinTemplates {
		tpl := &builtinTemplates[i]
		if _, ok := reg.byID[tpl.ID]; ok {
			return nil, fmt.Errorf("duplicate template id: %s", tpl.ID)
		}
		reg.byID[tpl.ID] = *tpl
		reg.byKind[tpl.Kind] = *tpl
		reg.list = append(reg.list, *tpl)
	}

	return reg, nil
}

// All 返回所有模板。
func (r *Registry) All() []CardTemplate { return r.list }

// ByID 按 id 查找模板。
func (r *Registry) ByID(id string) (CardTemplate, bool) {
	t, ok := r.byID[id]
	return t, ok
}

// ByKind 按 kind 查找模板。
func (r *Registry) ByKind(kind string) (CardTemplate, bool) {
	t, ok := r.byKind[kind]
	return t, ok
}
