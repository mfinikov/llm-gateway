package http

import (
	"strconv"
	"strings"

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

const maxBodySizeBytes = 1 << 20 // 1 MiB

func writeJSONError(ctx *fasthttp.RequestCtx, statusCode int, msg string) {
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")
	_, _ = ctx.WriteString(`{"error":` + strconv.Quote(msg) + `}`)
}

func writeJSONSuccess(ctx *fasthttp.RequestCtx, reply, model string) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	_, _ = ctx.WriteString(`{"status":"success","reply":` + strconv.Quote(reply) + `,"model":` + strconv.Quote(model) + `}`)
}

// HandleChatRequest processes incoming chat requests and routes them to the appropriate model
// It must be POST "/chat" and it uses jsonparser for fast JSON parsing
func HandleChatRequest(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != "/chat" {
		writeJSONError(ctx, fasthttp.StatusNotFound, "Not Found")
		return
	}

	if !ctx.IsPost() {
		writeJSONError(ctx, fasthttp.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	contentType := strings.ToLower(string(ctx.Request.Header.ContentType()))
	if !strings.HasPrefix(contentType, "application/json") {
		writeJSONError(ctx, fasthttp.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	body := ctx.PostBody()
	if len(body) == 0 {
		writeJSONError(ctx, fasthttp.StatusBadRequest, "Request body is required")
		return
	}
	if len(body) > maxBodySizeBytes {
		writeJSONError(ctx, fasthttp.StatusRequestEntityTooLarge, "Request body too large")
		return
	}

	model, err := jsonparser.GetString(body, "model")
	if err != nil {
		writeJSONError(ctx, fasthttp.StatusBadRequest, "Missing or invalid model field")
		return
	}

	message, err := jsonparser.GetString(body, "message")
	if err != nil {
		writeJSONError(ctx, fasthttp.StatusBadRequest, "Missing or invalid message field")
		return
	}

	if model == "" || message == "" {
		writeJSONError(ctx, fasthttp.StatusBadRequest, "Model and message are required")
		return
	}

	writeJSONSuccess(ctx, "Echo: "+message, model)
}
