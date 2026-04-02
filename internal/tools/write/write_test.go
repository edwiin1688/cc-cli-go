package write

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
	"github.com/liao-eli/cc-cli-go/internal/tools"
)

func TestWriteTool_Name(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, "Write", tool.Name())
}

func TestWriteTool_Description(t *testing.T) {
	tool := New()
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestWriteTool_IsReadOnly(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, false, tool.IsReadOnly(nil))
}

func TestWriteTool_IsConcurrencySafe(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, false, tool.IsConcurrencySafe(nil))
}

func TestWriteTool_Execute_Success(t *testing.T) {
	tool := New()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"

	input := map[string]interface{}{
		"file_path": filePath,
		"content":   content,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	data, err := os.ReadFile(filePath)
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, content, string(data))
}

func TestWriteTool_Execute_FileExists(t *testing.T) {
	tool := New()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "existing.txt")

	if err := os.WriteFile(filePath, []byte("old"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	input := map[string]interface{}{
		"file_path": filePath,
		"content":   "new",
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for existing file")
	}
	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "already exists")
	}
}

func TestWriteTool_Execute_EmptyPath(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"file_path": "",
		"content":   "test",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for empty file path")
	}
}

func TestWriteTool_Execute_CreateParentDirs(t *testing.T) {
	tool := New()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "subdir", "deep", "test.txt")
	content := "nested content"

	input := map[string]interface{}{
		"file_path": filePath,
		"content":   content,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	data, err := os.ReadFile(filePath)
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, content, string(data))
}

func TestWriteTool_UserFacingName(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"file_path": "/path/to/file.txt",
	}

	testutil.AssertEqual(t, "/path/to/file.txt", tool.UserFacingName(input))
	testutil.AssertEqual(t, "Write", tool.UserFacingName(nil))
}
