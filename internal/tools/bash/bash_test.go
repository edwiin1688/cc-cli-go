package bash

import (
	"context"
	"strings"
	"testing"

	"github.com/user-name/cc-cli-go/internal/testutil"
	"github.com/user-name/cc-cli-go/internal/tools"
)

func TestBashTool_Name(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, "Bash", tool.Name())
}

func TestBashTool_Description(t *testing.T) {
	tool := New()
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestBashTool_IsReadOnly(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, false, tool.IsReadOnly(nil))
}

func TestBashTool_IsConcurrencySafe(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, false, tool.IsConcurrencySafe(nil))
}

func TestBashTool_Execute_Success(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"command": "echo 'hello world'",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "hello world")
	}
}

func TestBashTool_Execute_EmptyCommand(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"command": "",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for empty command")
	}
}

func TestBashTool_Execute_CommandFailure(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"command": "exit 1",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "Error:")
	}
}

func TestBashTool_Execute_CommandNotFound(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"command": "nonexistent_command_12345",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "Error:")
	}
}

func TestBashTool_Execute_WithTimeout(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"command": "sleep 0.1",
		"timeout": 1000,
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}
}

func TestBashTool_Execute_TimeoutExceeded(t *testing.T) {
	tool := New()

	ctx := context.Background()
	input := map[string]interface{}{
		"command": "sleep 10",
		"timeout": 100,
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(ctx, input, tc)

	testutil.AssertNoError(t, err)
	if content, ok := result.Content.(string); ok {
		if content != "" {
			testutil.AssertContains(t, content, "Error:")
		}
	}
}

func TestBashTool_Execute_Stderr(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"command": "echo 'error message' >&2",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "error message")
	}
}

func TestBashTool_Execute_WorkingDirectory(t *testing.T) {
	tool := New()
	tmpDir := t.TempDir()

	input := map[string]interface{}{
		"command": "pwd",
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, tmpDir)
	}
}

func TestBashTool_Execute_PipedCommands(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"command": "echo 'hello world' | grep 'hello'",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "hello world")
	}
}

func TestBashTool_UserFacingName(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"command": "ls -la",
	}

	testutil.AssertEqual(t, "ls -la", tool.UserFacingName(input))

	longCmd := strings.Repeat("x", 60)
	input = map[string]interface{}{
		"command": longCmd,
	}

	expected := longCmd[:50] + "..."
	testutil.AssertEqual(t, expected, tool.UserFacingName(input))

	testutil.AssertEqual(t, "Bash", tool.UserFacingName(nil))
}
