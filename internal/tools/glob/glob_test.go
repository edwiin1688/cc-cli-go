package glob

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
	"github.com/liao-eli/cc-cli-go/internal/tools"
)

func TestGlobTool_Name(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, "Glob", tool.Name())
}

func TestGlobTool_Description(t *testing.T) {
	tool := New()
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestGlobTool_IsReadOnly(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, true, tool.IsReadOnly(nil))
}

func TestGlobTool_IsConcurrencySafe(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, true, tool.IsConcurrencySafe(nil))
}

func TestGlobTool_Execute_Success(t *testing.T) {
	tool := New()
	tmpDir := testutil.TempDirWithFiles(t, map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		"file3.go":  "content3",
	})

	input := map[string]interface{}{
		"pattern": "*.txt",
		"path":    tmpDir,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "Found 2 files")
		testutil.AssertContains(t, content, "file1.txt")
		testutil.AssertContains(t, content, "file2.txt")
	}
}

func TestGlobTool_Execute_EmptyPattern(t *testing.T) {
	tool := New()
	tmpDir := t.TempDir()

	input := map[string]interface{}{
		"pattern": "",
		"path":    tmpDir,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for empty pattern")
	}
}

func TestGlobTool_Execute_NoMatches(t *testing.T) {
	tool := New()
	tmpDir := testutil.TempDirWithFiles(t, map[string]string{
		"file1.txt": "content",
	})

	input := map[string]interface{}{
		"pattern": "*.go",
		"path":    tmpDir,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success with no matches, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "No files found")
	}
}

func TestGlobTool_Execute_DirectoryNotFound(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"pattern": "*.txt",
		"path":    "/nonexistent/directory",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for nonexistent directory")
	}
}

func TestGlobTool_Execute_DefaultPath(t *testing.T) {
	tool := New()
	tmpDir := testutil.TempDirWithFiles(t, map[string]string{
		"test.txt": "content",
	})

	input := map[string]interface{}{
		"pattern": "*.txt",
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}
}

func TestGlobTool_Execute_NestedFiles(t *testing.T) {
	tool := New()
	tmpDir := t.TempDir()

	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	files := []string{
		filepath.Join(tmpDir, "root.txt"),
		filepath.Join(subdir, "nested.txt"),
	}

	for _, f := range files {
		if err := os.WriteFile(f, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	input := map[string]interface{}{
		"pattern": "*.txt",
		"path":    tmpDir,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}
}

func TestGlobTool_UserFacingName(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"pattern": "*.go",
	}

	testutil.AssertEqual(t, "*.go", tool.UserFacingName(input))
	testutil.AssertEqual(t, "Glob", tool.UserFacingName(nil))
}
