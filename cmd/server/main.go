package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sanbei101/blue-card-engine/internal/fonts"
	"github.com/sanbei101/blue-card-engine/internal/server"
	"github.com/sanbei101/blue-card-engine/internal/templates"
)

func main() {
	lib, err := fonts.NewLibrary()
	if err != nil {
		log.Fatalf("load fonts: %v", err)
	}

	reg, err := templates.Load()
	if err != nil {
		log.Fatalf("load templates: %v", err)
	}

	handler := server.NewHandler(reg, lib)
	mux := http.NewServeMux()
	handler.Register(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5174"
	}

	addr := ":" + port
	log.Printf("Card Engine Server running at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
