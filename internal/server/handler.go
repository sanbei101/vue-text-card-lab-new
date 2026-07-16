package server

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/phuslu/log"

	"github.com/sanbei101/blue-card-engine/internal/cardengine"
	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/render"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

type TemplateItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type TemplateListResponse struct {
	Templates []TemplateItem `json:"templates"`
}

type CardItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	SVG  string `json:"svg"`
}

type CardListResponse struct {
	Templates []CardItem `json:"templates"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	templates *templates.Registry
	library   *fonts.Library
}

func NewHandler(reg *templates.Registry, lib *fonts.Library) *Handler {
	return &Handler{
		templates: reg,
		library:   lib,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/api/cards", h.cards)
	mux.HandleFunc("/api/templates", h.templateList)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Error().Err(err).Msg("write health response failed")
	}
}

func (h *Handler) templateList(w http.ResponseWriter, _ *http.Request) {
	allTemplates := h.templates.All()
	result := make([]TemplateItem, 0, len(allTemplates))

	for i := range allTemplates {
		tpl := &allTemplates[i]
		result = append(result, TemplateItem{
			ID:   tpl.ID,
			Name: tpl.Name,
			Kind: tpl.Kind,
		})
	}

	writeJSON(w, http.StatusOK, TemplateListResponse{Templates: result})
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

	allTemplates := h.templates.All()
	results := make([]CardItem, 0, len(allTemplates))

	for i := range allTemplates {
		tpl := &allTemplates[i]
		svg, err := render.Card(tpl, title, keyword, h.library)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error: fmt.Sprintf("render card %s failed", tpl.ID),
			})
			return
		}

		results = append(results, CardItem{
			ID:   tpl.ID,
			Name: tpl.Name,
			SVG:  svg,
		})
	}

	writeJSON(w, http.StatusOK, CardListResponse{Templates: results})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.MarshalWrite(w, data); err != nil {
		log.Error().Err(err).Msg("write json response failed")
	}
}
