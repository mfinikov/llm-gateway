package http

import (
	"github.com/buger/jsonparser"
	"github.com/valyala/fasthttp"
)

// ChatRequest represents the structure of a chat request payload
type ChatRequest struct {
	Model   string `json:"model"`
	Message string `json:"message"`
}

// ChatResponse represents the structure of a chat response payload
type ChatResponse struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
	Model  string `json:"model"`
}

// HandleChatRequest processes incoming chat requests and routes them to the appropriate model
// It must be POST "/chat" and it uses
func HandleChatRequest(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != "/chat" {
		ctx.Error("Not Found", fasthttp.StatusNotFound)
		return
	}

	if !ctx.IsPost() {
		ctx.Error("Method Not Allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	body := ctx.PostBody()

	model, err := jsonparser.GetString(body, "model")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.WriteString(`{"error":"Missing or invalid model field"}`)
		return
	}

	message, err := jsonparser.GetString(body, "message")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.WriteString(`{"error":"Missing or invalid message field"}`)
		return
	}

	if model == "" || message == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.WriteString(`{"error":"Model and message fields cannot be empty"}`)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	ctx.WriteString(`{"status":"success","reply":"Echo: `)
	ctx.WriteString(message)
	ctx.WriteString(`","model":"`)
	ctx.WriteString(model)
	ctx.WriteString(`"}`)
}
