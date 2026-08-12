package main

import (
	"io"
	"net/http"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/joho/godotenv"
	"github.com/phuslu/log"
	"github.com/purus-dev/aqua"

	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/render"
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

func vipsLogHandler(messageDomain string, messageLevel vips.LogLevel, message string) {
	switch messageLevel {
	case vips.LogLevelError, vips.LogLevelCritical:
		log.Error().Str("domain", messageDomain).Msg(message)
	case vips.LogLevelWarning:
		log.Warn().Str("domain", messageDomain).Msg(message)
	case vips.LogLevelInfo, vips.LogLevelMessage:
		log.Info().Str("domain", messageDomain).Msg(message)
	case vips.LogLevelDebug:
		log.Debug().Str("domain", messageDomain).Msg(message)
	default:
		log.Info().Str("domain", messageDomain).Msg(message)
	}
}

func main() {
	initLog()
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("No .env file found, using environment variables")
	}
	vips.LoggingSettings(vipsLogHandler, vips.LogLevelInfo)
	vips.Startup(nil)
	defer vips.Shutdown()

	r2Config := aqua.Config{
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("R2_ACCESS_KEY_SECRET"),
		Bucket:          os.Getenv("R2_BUCKET"),
		Endpoint:        os.Getenv("R2_ENDPOINT"),
		Region:          os.Getenv("R2_REGION"),
		UsePathStyle:    true,
	}

	if err := r2Config.Validate(); err != nil {
		log.Fatal().Err(err).Msg("invalid R2 configuration")
	}

	r2Storage := render.NewR2Storage(r2Config, "https://image-bed.sanbei.codes")
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

	handler := server.NewHandler(reg, lib, r2Storage)
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
