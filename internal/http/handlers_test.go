package http

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestHandleChatRequest_Success(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/chat")
	ctx.Request.SetBodyString(`{"model":"gpt-4","message":"Hello World"}`)

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
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/wrong")
	ctx.Request.SetBodyString(`{"model":"gpt-4","message":"Hello"}`)

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
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/chat")
	ctx.Request.SetBodyString(`{"message":"Hello"}`)

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
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/chat")
	ctx.Request.SetBodyString(`{"model":"gpt-4"}`)

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
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/chat")
	ctx.Request.SetBodyString(`{"model":"","message":""}`)

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
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetRequestURI("/chat")
	ctx.Request.SetBodyString(`{invalid json}`)

	HandleChatRequest(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fasthttp.StatusBadRequest, ctx.Response.StatusCode())
	}
}
