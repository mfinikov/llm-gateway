package http

import (
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

func newJSONPostCtx(path, body string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetRequestURI(path)
	ctx.Request.SetBodyString(body)
	return ctx
}

func TestHandleChatRequest_Success(t *testing.T) {
	ctx := newJSONPostCtx("/chat", `{"model":"gpt-4","message":"Hello World"}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusOK, ctx.Response.StatusCode())
	}

	expectedBody := `{"status":"success","reply":"Echo: Hello World","model":"gpt-4"}`
	actualBody := string(ctx.Response.Body())
	if actualBody != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
	}

	contentType := string(ctx.Response.Header.ContentType())
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}

func TestHandleChatRequest_WrongPath(t *testing.T) {
	ctx := newJSONPostCtx("/wrong", `{"model":"gpt-4","message":"Hello"}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusNotFound {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusNotFound, ctx.Response.StatusCode())
	}
}

func TestHandleChatRequest_WrongMethod(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.SetRequestURI("/chat")

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusMethodNotAllowed, ctx.Response.StatusCode())
	}
}

func TestHandleChatRequest_MissingModel(t *testing.T) {
	ctx := newJSONPostCtx("/chat", `{"message":"Hello"}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusBadRequest, ctx.Response.StatusCode())
	}

	expectedBody := `{"error":"Missing or invalid model field"}`
	actualBody := string(ctx.Response.Body())
	if actualBody != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
	}
}

func TestHandleChatRequest_MissingMessage(t *testing.T) {
	ctx := newJSONPostCtx("/chat", `{"model":"gpt-4"}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusBadRequest, ctx.Response.StatusCode())
	}

	expectedBody := `{"error":"Missing or invalid message field"}`
	actualBody := string(ctx.Response.Body())
	if actualBody != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
	}
}

func TestHandleChatRequest_EmptyFields(t *testing.T) {
	ctx := newJSONPostCtx("/chat", `{"model":"","message":""}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusBadRequest, ctx.Response.StatusCode())
	}

	expectedBody := `{"error":"Model and message are required"}`
	actualBody := string(ctx.Response.Body())
	if actualBody != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
	}
}

func TestHandleChatRequest_InvalidJSON(t *testing.T) {
	ctx := newJSONPostCtx("/chat", `{invalid json}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusBadRequest, ctx.Response.StatusCode())
	}
}

func TestHandleChatRequest_InvalidContentType(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.Header.SetContentType("text/plain")
	ctx.Request.SetRequestURI("/chat")
	ctx.Request.SetBodyString(`{"model":"gpt-4","message":"Hello"}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusUnsupportedMediaType {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusUnsupportedMediaType, ctx.Response.StatusCode())
	}
}

func TestHandleChatRequest_EmptyBody(t *testing.T) {
	ctx := newJSONPostCtx("/chat", "")

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusBadRequest, ctx.Response.StatusCode())
	}
}

func TestHandleChatRequest_BodyTooLarge(t *testing.T) {
	largeMessage := strings.Repeat("a", maxBodySizeBytes+1)
	ctx := newJSONPostCtx("/chat", `{"model":"gpt-4","message":"`+largeMessage+`"}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusRequestEntityTooLarge {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusRequestEntityTooLarge, ctx.Response.StatusCode())
	}
}

func TestHandleChatRequest_EscapesJSONOutput(t *testing.T) {
	ctx := newJSONPostCtx("/chat", `{"model":"gpt-4","message":"Hello \"quoted\"\nline"}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusOK, ctx.Response.StatusCode())
	}

	expectedBody := `{"status":"success","reply":"Echo: Hello \"quoted\"\nline","model":"gpt-4"}`
	actualBody := string(ctx.Response.Body())
	if actualBody != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
	}
}
