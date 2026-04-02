package grep

import (
	"context"
	"testing"

	"github.com/user-name/cc-cli-go/internal/testutil"
	"github.com/user-name/cc-cli-go/internal/tools"
)

func TestGrepTool_Name(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, "Grep", tool.Name())
}

func TestGrepTool_Description(t *testing.T) {
	tool := New()
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestGrepTool_IsReadOnly(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, true, tool.IsReadOnly(nil))
}

func TestGrepTool_IsConcurrencySafe(t *testing.T) {
	tool := New()
	testutil.AssertEqual(t, true, tool.IsConcurrencySafe(nil))
}

func TestGrepTool_Execute_Success(t *testing.T) {
	tool := New()
	tmpDir := testutil.TempDirWithFiles(t, map[string]string{
		"test.txt": "hello world\nhello universe\nfoo bar",
	})

	input := map[string]interface{}{
		"pattern": "hello",
		"path":    tmpDir,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "Found 2 matches")
	}
}

func TestGrepTool_Execute_EmptyPattern(t *testing.T) {
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

func TestGrepTool_Execute_InvalidRegex(t *testing.T) {
	tool := New()
	tmpDir := t.TempDir()

	input := map[string]interface{}{
		"pattern": "[invalid",
		"path":    tmpDir,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for invalid regex")
	}
}

func TestGrepTool_Execute_NoMatches(t *testing.T) {
	tool := New()
	tmpDir := testutil.TempDirWithFiles(t, map[string]string{
		"test.txt": "hello world",
	})

	input := map[string]interface{}{
		"pattern": "nonexistent",
		"path":    tmpDir,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success with no matches, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "No matches found")
	}
}

func TestGrepTool_Execute_WithIncludeFilter(t *testing.T) {
	tool := New()
	tmpDir := testutil.TempDirWithFiles(t, map[string]string{
		"test.txt": "hello world",
		"test.go":  "hello universe",
	})

	input := map[string]interface{}{
		"pattern": "hello",
		"path":    tmpDir,
		"include": "*.go",
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "Found 1 matches")
	}
}

func TestGrepTool_Execute_RegexPattern(t *testing.T) {
	tool := New()
	tmpDir := testutil.TempDirWithFiles(t, map[string]string{
		"test.txt": "foo123\nbar456\nbaz789",
	})

	input := map[string]interface{}{
		"pattern": "[0-9]+",
		"path":    tmpDir,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	if content, ok := result.Content.(string); ok {
		testutil.AssertContains(t, content, "Found 3 matches")
	}
}

func TestGrepTool_Execute_DirectoryNotFound(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"pattern": "test",
		"path":    "/nonexistent/directory",
	}
	tc := &tools.ToolContext{WorkingDir: "/tmp"}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if !result.IsError {
		t.Error("expected error for nonexistent directory")
	}
}

func TestGrepTool_Execute_DefaultPath(t *testing.T) {
	tool := New()
	tmpDir := testutil.TempDirWithFiles(t, map[string]string{
		"test.txt": "hello world",
	})

	input := map[string]interface{}{
		"pattern": "hello",
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	result, err := tool.Execute(context.Background(), input, tc)

	testutil.AssertNoError(t, err)
	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}
}

func TestGrepTool_UserFacingName(t *testing.T) {
	tool := New()

	input := map[string]interface{}{
		"pattern": "hello",
	}

	testutil.AssertEqual(t, "hello", tool.UserFacingName(input))
	testutil.AssertEqual(t, "Grep", tool.UserFacingName(nil))
}
