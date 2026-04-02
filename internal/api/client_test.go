package api

import (
	"context"
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
	"github.com/liao-eli/cc-cli-go/internal/types"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client == nil {
		t.Error("expected non-nil client")
	}
}

func TestNewRequest(t *testing.T) {
	req := NewRequest("claude-3-5-sonnet-20241022", 4096)

	if req.Model != "claude-3-5-sonnet-20241022" {
		t.Errorf("expected model claude-3-5-sonnet-20241022, got %s", req.Model)
	}

	if req.MaxTokens != 4096 {
		t.Errorf("expected max tokens 4096, got %d", req.MaxTokens)
	}

	if !req.Stream {
		t.Error("expected stream to be true")
	}
}

func TestRequest_SetSystem(t *testing.T) {
	req := NewRequest("claude-3-5-sonnet-20241022", 4096)
	prompts := []string{"You are helpful.", "Be concise."}

	req.SetSystem(prompts)

	if len(req.System) != 2 {
		t.Errorf("expected 2 system blocks, got %d", len(req.System))
	}

	if req.System[0].Text != "You are helpful." {
		t.Errorf("expected first system prompt, got %s", req.System[0].Text)
	}
}

func TestRequest_AddMessage(t *testing.T) {
	req := NewRequest("claude-3-5-sonnet-20241022", 4096)

	msg := &types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			{Type: "text", Text: "Hello"},
		},
	}

	req.AddMessage(msg)

	if len(req.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(req.Messages))
	}

	if req.Messages[0].Role != "user" {
		t.Errorf("expected user role, got %s", req.Messages[0].Role)
	}
}

func TestRequest_AddTool(t *testing.T) {
	req := NewRequest("claude-3-5-sonnet-20241022", 4096)

	tool := ToolParam{
		Name:        "Read",
		Description: "Read a file",
		InputSchema: map[string]interface{}{"type": "object"},
	}

	req.AddTool(tool)

	if len(req.Tools) != 1 {
		t.Errorf("expected 1 tool, got %d", len(req.Tools))
	}

	if req.Tools[0].Name != "Read" {
		t.Errorf("expected Read tool, got %s", req.Tools[0].Name)
	}
}

func TestMockClient_Stream_Success(t *testing.T) {
	client := NewMockClient()
	client.AddTextEvent("Hello")
	client.AddTextEvent("World")

	req := NewRequest("claude-3-5-sonnet-20241022", 4096)

	stream, err := client.Stream(context.Background(), req)
	testutil.AssertNoError(t, err)

	events := []StreamEvent{}
	for event := range stream {
		events = append(events, event)
	}

	if len(events) == 0 {
		t.Error("expected at least one event")
	}

	lastEvent := events[len(events)-1]
	if lastEvent.Type != "message_stop" {
		t.Errorf("expected message_stop, got %s", lastEvent.Type)
	}
}

func TestMockClient_Stream_WithError(t *testing.T) {
	client := NewMockClient()
	client.SetError(context.Canceled)

	req := NewRequest("claude-3-5-sonnet-20241022", 4096)

	stream, err := client.Stream(context.Background(), req)

	if err == nil {
		t.Error("expected error")
	}

	if stream == nil {
		t.Error("expected non-nil stream channel")
	}
}

func TestMockClient_Stream_ContextCancellation(t *testing.T) {
	client := NewMockClient()
	client.AddTextEvent("Test")

	req := NewRequest("claude-3-5-sonnet-20241022", 4096)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stream, err := client.Stream(ctx, req)
	testutil.AssertNoError(t, err)

	_, ok := <-stream
	if ok {
		t.Error("expected stream to be closed due to context cancellation")
	}
}

func TestParseMessageStart(t *testing.T) {
	data := []byte(`{"type":"message","id":"msg_123","role":"assistant","content":[]}`)

	msg, err := ParseMessageStart(data)
	testutil.AssertNoError(t, err)

	if msg.UUID != "msg_123" {
		t.Errorf("expected msg_123, got %s", msg.UUID)
	}

	if msg.Role != "assistant" {
		t.Errorf("expected assistant role, got %s", msg.Role)
	}
}

func TestParseContentBlock(t *testing.T) {
	data := []byte(`{"type":"text","text":"Hello world"}`)

	block, err := ParseContentBlock(data)
	testutil.AssertNoError(t, err)

	if block.Type != "text" {
		t.Errorf("expected text type, got %s", block.Type)
	}

	if block.Text != "Hello world" {
		t.Errorf("expected Hello world, got %s", block.Text)
	}
}

func TestParseDelta(t *testing.T) {
	data := []byte(`{"type":"text_delta","text":"Hello"}`)

	deltaType, text, err := ParseDelta(data)
	testutil.AssertNoError(t, err)

	if deltaType != "text_delta" {
		t.Errorf("expected text_delta, got %s", deltaType)
	}

	if text != "Hello" {
		t.Errorf("expected Hello, got %s", text)
	}
}
