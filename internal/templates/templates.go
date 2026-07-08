package templates

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed templates.json
var templatesFS embed.FS

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
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Kind       string      `json:"kind"`
	Background string      `json:"background"`
	Foreground string      `json:"foreground"`
	Accent     string      `json:"accent"`
	Muted      string      `json:"muted"`
	FontFamily string      `json:"fontFamily"`
	FontWeight int         `json:"fontWeight"`
	TextBox    TextBoxSpec `json:"textBox"`
	Radius     float64     `json:"radius"`
	Frame      string      `json:"frame,omitempty"`
}

// Registry 保存所有已加载模板。
type Registry struct {
	byID   map[string]CardTemplate
	byKind map[string]CardTemplate
	list   []CardTemplate
}

// Load 从嵌入的 templates.json 加载模板注册表。
func Load() (*Registry, error) {
	data, err := templatesFS.ReadFile("templates.json")
	if err != nil {
		return nil, fmt.Errorf("read templates.json: %w", err)
	}

	var list []CardTemplate
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("parse templates.json: %w", err)
	}

	reg := &Registry{
		byID:   make(map[string]CardTemplate),
		byKind: make(map[string]CardTemplate),
		list:   make([]CardTemplate, 0, len(list)),
	}

	for _, tpl := range list {
		if _, ok := reg.byID[tpl.ID]; ok {
			return nil, fmt.Errorf("duplicate template id: %s", tpl.ID)
		}
		reg.byID[tpl.ID] = tpl
		reg.byKind[tpl.Kind] = tpl
		reg.list = append(reg.list, tpl)
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
