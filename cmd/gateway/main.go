package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mfinikov/llm-gateway/internal/http"
	"github.com/valyala/fasthttp"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	server := &fasthttp.Server{
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Name:         "llm-gateway",
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutdown signal received, stopping server...")
		if err := server.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	log.Printf("Starting Gateway Server on %s...", addr)
	if err := server.ListenAndServe(addr); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func router(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	switch path {
	case "/chat":
		http.HandleChatRequest(ctx)
	case "/health":
		handleHealth(ctx)
	case "/":
		handleHome(ctx)
	default:
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		_, _ = ctx.WriteString(`{"error":` + strconv.Quote("Not Found") + `}`)
	}
}

func handleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	_, _ = ctx.WriteString(`{"status":"healthy"}`)
}

func handleHome(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	_, _ = ctx.WriteString(`{"message":"LLM Gateway API","version":"1.0.0"}`)
}
