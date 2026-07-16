package main

import (
	"net/http"
	"os"

	"github.com/phuslu/log"
	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/server"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

func main() {
	lib, err := fonts.NewLibrary()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create font library")
	}

	reg, err := templates.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load templates")
	}

	handler := server.NewHandler(reg, lib)
	mux := http.NewServeMux()
	handler.Register(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5174"
	}

	addr := ":" + port
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
	log.Info().Msgf("server is running at http://localhost%s", addr)
}
