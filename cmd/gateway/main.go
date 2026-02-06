package main

import (
	"log"

	"github.com/mfinikov/llm-gateway/internal/http"
	"github.com/valyala/fasthttp"
)

func main() {
	log.Println("Starting Gateway Server on :8080...")

	if err := fasthttp.ListenAndServe(":8080", router); err != nil {
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
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
}

func handleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.WriteString(`{"status":"healthy"}`)
}

func handleHome(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.WriteString(`{"message":"LLM Gateway API","version":"1.0.0"}`)
}
