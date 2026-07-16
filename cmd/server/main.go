package main

import (
	"io"
	"net/http"
	"os"

	"github.com/phuslu/log"

	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/server"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

func initLog() {
	env := os.Getenv("ENV")
	if env == "production" {
		log.DefaultLogger = log.Logger{
			Level: log.InfoLevel,
			Writer: &log.AsyncWriter{
				ChannelSize:   4096,
				DiscardOnFull: false,
				Writer:        &log.IOWriter{Writer: os.Stderr},
			},
		}
	} else {
		log.DefaultLogger = log.Logger{
			Level:      log.DebugLevel,
			TimeFormat: "15:04:05",
			Caller:     1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}
}

func main() {
	initLog()
	if closer, ok := log.DefaultLogger.Writer.(io.Closer); ok {
		defer closer.Close()
	}

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

	addr := "0.0.0.0:" + port
	log.Info().Msgf("server is running at %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
