package edit

import (
	"context"
	"os"
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
	"github.com/liao-eli/cc-cli-go/internal/tools"
)

func TestEditTool_Name(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, "Edit", tool.Name())
}

func TestEditTool_Description(t *testing.T) {
	tool := New()
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestEditTool_IsReadOnly(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, false, tool.IsReadOnly(nil))
}

func TestEditTool_IsConcurrencySafe(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, false, tool.IsConcurrencySafe(nil))
}

func TestEditTool_Execute_Success(t *testing.T) {
	tool := New()
	content := "hello world"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path":  filePath,
		"old_string": "world",
		"new_string": "universe",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	data, err := os.ReadFile(filePath)
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, "hello universe", string(data))
}

func TestEditTool_Execute_FileNotFound(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"file_path":  "/nonexistent/file.txt",
		"old_string": "old",
		"new_string": "new",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for nonexistent file")
	}
}

func TestEditTool_Execute_EmptyPath(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"file_path":  "",
		"old_string": "old",
		"new_string": "new",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for empty file path")
	}
}

func TestEditTool_Execute_EmptyOldString(t *testing.T) {
	tool := New()
	filePath := testutil.TempFile(t, "content")

	input := map[string]interface{}{
		"file_path":  filePath,
		"old_string": "",
		"new_string": "new",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for empty old_string")
	}
}

func TestEditTool_Execute_OldStringNotFound(t *testing.T) {
	tool := New()
	content := "hello world"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path":  filePath,
		"old_string": "nonexistent",
		"new_string": "new",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for old_string not found")
	}
}

func TestEditTool_Execute_MultipleMatches(t *testing.T) {
	tool := New()
	content := "hello hello hello"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path":  filePath,
		"old_string": "hello",
		"new_string": "hi",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for multiple matches")
	}
}

func TestEditTool_Execute_ExactMatch(t *testing.T) {
	tool := New()
	content := "foo bar baz\nbar qux\nbar"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path":  filePath,
		"old_string": "bar qux",
		"new_string": "BAR QUX",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	data, err := os.ReadFile(filePath)
	testutil.AssertNoError(t, err)
	testutil.AssertContains(t, string(data), "BAR QUX")
}

func TestEditTool_Execute_ReplaceWithEmpty(t *testing.T) {
	tool := New()
	content := "hello world"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path":  filePath,
		"old_string": " world",
		"new_string": "",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	data, err := os.ReadFile(filePath)
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, "hello", string(data))
}

func TestEditTool_Execute_MultilineReplacement(t *testing.T) {
	tool := New()
	content := "line1\nline2\nline3"
	filePath := testutil.TempFile(t, content)

	input := map[string]interface{}{
		"file_path":  filePath,
		"old_string": "line1\nline2",
		"new_string": "new1\nnew2",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	data, err := os.ReadFile(filePath)
	testutil.AssertNoError(t, err)
	testutil.AssertContains(t, string(data), "new1")
	testutil.AssertContains(t, string(data), "new2")
}

func TestEditTool_UserFacingName(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"file_path": "/path/to/file.txt",
	}

	testutil.AssertEqual(t, "/path/to/file.txt", tool.UserFacingName(input))
	testutil.AssertEqual(t, "Edit", tool.UserFacingName(nil))
}
