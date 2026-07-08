package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sanbei101/blue-card-engine/internal/cardengine"
	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/render"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

// Handler 持有渲染所需依赖。
type Handler struct {
	templates  *templates.Registry
	library    *fonts.Library
	measurers  map[string]*cardengine.FontMeasurer
	cacheSize  int
}

// NewHandler 创建路由处理器。
func NewHandler(reg *templates.Registry, lib *fonts.Library) *Handler {
	return &Handler{
		templates: reg,
		library:   lib,
		measurers: make(map[string]*cardengine.FontMeasurer),
		cacheSize: 4096,
	}
}

// Register 注册 HTTP 路由。
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/api/cards", h.cards)
	mux.HandleFunc("/api/templates", h.templateList)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) templateList(w http.ResponseWriter, r *http.Request) {
	var result []map[string]any
	for _, tpl := range h.templates.All() {
		result = append(result, map[string]any{
			"id":   tpl.ID,
			"name": tpl.Name,
			"kind": tpl.Kind,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"templates": result})
}

func (h *Handler) cards(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	keyword := r.URL.Query().Get("keyword")
	if title == "" {
		title = cardengine.DefaultEmptyText
	}
	if keyword == "" {
		keyword = cardengine.InferKeyword(title)
	}

	var results []map[string]any
	for _, tpl := range h.templates.All() {
		measurer := h.measurerFor(tpl.FontFamily)
		svg, err := render.RenderCard(tpl, title, keyword, measurer)
		if err != nil {
			log.Printf("render card %s failed: %v", tpl.ID, err)
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"error": fmt.Sprintf("render card %s failed", tpl.ID),
			})
			return
		}

		results = append(results, map[string]any{
			"id":   tpl.ID,
			"name": tpl.Name,
			"svg":  svg,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{"templates": results})
}

func (h *Handler) measurerFor(fontFamily string) *cardengine.FontMeasurer {
	family := h.library.Resolve(fontFamily)
	if m, ok := h.measurers[family]; ok {
		return m
	}
	m := cardengine.NewFontMeasurer(h.library.FaceProvider(family), h.cacheSize)
	h.measurers[family] = m
	return m
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("encode json failed: %v", err)
	}
}
