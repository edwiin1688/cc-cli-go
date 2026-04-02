package read

import (
	"context"
	"testing"

	"github.com/user-name/cc-cli-go/internal/testutil"
	"github.com/user-name/cc-cli-go/internal/tools"
)

func TestReadTool_Name(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, "Read", tool.Name())
}

func TestReadTool_Description(t *testing.T) {
	tool := New()
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestReadTool_IsReadOnly(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, true, tool.IsReadOnly(nil))
}

func TestReadTool_IsConcurrencySafe(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, true, tool.IsConcurrencySafe(nil))
}

func TestReadTool_Execute_Success(t *testing.T) {
	tool := New()
	content := "line1\nline2\nline3"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path": filePath,
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "line1")
		testutil.AssertContains(t, content, "line2")
		testutil.AssertContains(t, content, "line3")
	}
}

func TestReadTool_Execute_FileNotFound(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"file_path": "/nonexistent/file.txt",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for nonexistent file")
	}
}

func TestReadTool_Execute_EmptyPath(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"file_path": "",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for empty file path")
	}
}

func TestReadTool_Execute_WithLimit(t *testing.T) {
	tool := New()
	content := "line1\nline2\nline3\nline4\nline5"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path": filePath,
		"limit":     2,
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "line1")
		testutil.AssertContains(t, content, "line2")
	}
}

func TestReadTool_Execute_WithOffset(t *testing.T) {
	tool := New()
	content := "line1\nline2\nline3\nline4\nline5"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path": filePath,
		"offset":    2,
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "line2")
		testutil.AssertContains(t, content, "line3")
	}
}

func TestReadTool_Execute_WithLimitAndOffset(t *testing.T) {
	tool := New()
	content := "line1\nline2\nline3\nline4\nline5"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path": filePath,
		"limit":     2,
		"offset":    2,
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "line2")
		testutil.AssertContains(t, content, "line3")
	}
}

func TestReadTool_Execute_EmptyFile(t *testing.T) {
	tool := New()
	filePath := testutil.TempFile(t, "")

	input := map[string]interface{}{
		"file_path": filePath,
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success for empty file, got error: %s", result.Content)
	}
}

func TestReadTool_Execute_LineNumbers(t *testing.T) {
	tool := New()
	content := "line1\nline2\nline3"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path": filePath,
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "1:")
		testutil.AssertContains(t, content, "2:")
		testutil.AssertContains(t, content, "3:")
	}
}

func TestReadTool_UserFacingName(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"file_path": "/path/to/file.txt",
	}

	testutil.AssertEqual(t, "/path/to/file.txt", tool.UserFacingName(input))
	testutil.AssertEqual(t, "Read", tool.UserFacingName(nil))
}
