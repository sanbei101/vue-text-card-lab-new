package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"connectrpc.com/connect"
	"github.com/phuslu/log"

	cardenginev1 "github.com/sanbei101/blue-card-engine/gen/proto/cardengine/v1"
	cardenginev1connect "github.com/sanbei101/blue-card-engine/gen/proto/cardengine/v1/cardenginev1connect"

	"github.com/sanbei101/blue-card-engine/internal/cardengine"
	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/render"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

type Handler struct {
	templates      *templates.Registry
	library        *fonts.Library
	svg2webpClient *render.Client
}

func NewHandler(reg *templates.Registry, lib *fonts.Library) *Handler {
	client := render.NewClient("https://svg2webp.sanbei.codes")
	return &Handler{
		templates:      reg,
		library:        lib,
		svg2webpClient: client,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)

	path, handler := cardenginev1connect.NewCardEngineServiceHandler(h)
	mux.Handle(path, handler)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Error().Err(err).Msg("write health response failed")
	}
}

func (h *Handler) TemplateList(
	ctx context.Context,
	req *cardenginev1.TemplateListRequest,
) (*cardenginev1.TemplateListResponse, error) {
	allTemplates := h.templates.All()
	result := make([]*cardenginev1.TemplateItem, 0, len(allTemplates))

	for i := range allTemplates {
		tpl := &allTemplates[i]
		item := &cardenginev1.TemplateItem{}
		item.SetId(tpl.ID)
		item.SetName(tpl.Name)
		item.SetKind(tpl.Kind)

		result = append(result, item)
	}
	resp := &cardenginev1.TemplateListResponse{}
	resp.SetTemplates(result)
	return resp, nil
}

func (h *Handler) Cards(
	ctx context.Context,
	req *cardenginev1.CardListRequest,
) (*cardenginev1.CardListResponse, error) {
	if utf8.RuneCountInString(req.GetTitle()) > 100 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("title length must be at most 100 characters"),
		)
	}
	title := req.GetTitle()
	keyword := req.GetKeyword()

	if title == "" {
		title = cardengine.DefaultEmptyText
	}
	if keyword == "" {
		keyword = cardengine.InferKeyword(title)
	}

	allTemplates := h.templates.All()
	results := make([]*cardenginev1.CardItem, 0, len(allTemplates))
	for i := range allTemplates {
		tpl := &allTemplates[i]
		svg, err := render.Card(tpl, title, keyword, h.library)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("render template %s (%s): %w", tpl.Kind, tpl.ID, err),
			)
		}
		svgURL, err := h.svg2webpClient.Convert(ctx, svg)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("convert svg to webp for template %s (%s): %w", tpl.Kind, tpl.ID, err),
			)
		}
		item := &cardenginev1.CardItem{}
		item.SetId(tpl.ID)
		item.SetName(tpl.Kind)
		item.SetUrl(svgURL.URL)

		results = append(results, item)
	}

	resp := &cardenginev1.CardListResponse{}
	resp.SetTemplates(results)

	return resp, nil
}
