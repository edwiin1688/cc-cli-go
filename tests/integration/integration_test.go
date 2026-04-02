package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/permission"
	"github.com/liao-eli/cc-cli-go/internal/testutil"
	"github.com/liao-eli/cc-cli-go/internal/tools"
	"github.com/liao-eli/cc-cli-go/internal/tools/read"
	"github.com/liao-eli/cc-cli-go/internal/tools/write"
)

func TestWriteAndReadWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "Hello, Integration Test!"

	writeTool := write.New()
	readTool := read.New()

	writeInput := map[string]interface{}{
		"file_path": filePath,
		"content":   content,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	writeResult, err := writeTool.Execute(context.Background(), writeInput, tc)
	testutil.AssertNoError(t, err)
	if writeResult.IsError {
		t.Fatalf("write failed: %s", writeResult.Content)
	}

	readInput := map[string]interface{}{
		"file_path": filePath,
	}
	readResult, err := readTool.Execute(context.Background(), readInput, tc)
	testutil.AssertNoError(t, err)
	if readResult.IsError {
		t.Fatalf("read failed: %s", readResult.Content)
	}

	if contentStr, ok := readResult.Content.(string); ok {
		testutil.AssertContains(t, contentStr, "Hello, Integration Test!")
	}
}

func TestPermissionWithToolExecution(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	checker := permission.NewChecker(permission.ModeDefault)

	readDecision := checker.Check("Read", map[string]interface{}{
		"file_path": filePath,
	})
	testutil.AssertEqual(t, permission.BehaviorAllow, readDecision.Behavior)

	writeDecision := checker.Check("Write", map[string]interface{}{
		"file_path": filePath,
		"content":   "new content",
	})
	testutil.AssertEqual(t, permission.BehaviorAsk, writeDecision.Behavior)
}

func TestDangerousCommandWithPermission(t *testing.T) {
	checker := permission.NewChecker(permission.ModeDefault)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "rm -rf /important/data",
	})

	testutil.AssertEqual(t, permission.BehaviorAsk, decision.Behavior)
	testutil.AssertContains(t, decision.Reason, "dangerous")
}

func TestToolRegistryIntegration(t *testing.T) {
	registry := tools.NewRegistry()

	registry.Register(write.New())
	registry.Register(read.New())

	writeTool := registry.Get("Write")
	if writeTool == nil {
		t.Fatal("expected Write tool to be registered")
	}

	readTool := registry.Get("Read")
	if readTool == nil {
		t.Fatal("expected Read tool to be registered")
	}

	allTools := registry.All()
	if len(allTools) < 2 {
		t.Errorf("expected at least 2 tools, got %d", len(allTools))
	}
}

func TestMultipleToolExecution(t *testing.T) {
	tmpDir := t.TempDir()

	writeTool := write.New()
	readTool := read.New()
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	files := map[string]string{
		"file1.txt": "Content 1",
		"file2.txt": "Content 2",
		"file3.txt": "Content 3",
	}

	for name, content := range files {
		filePath := filepath.Join(tmpDir, name)
		input := map[string]interface{}{
			"file_path": filePath,
			"content":   content,
		}

		result, err := writeTool.Execute(context.Background(), input, tc)
		testutil.AssertNoError(t, err)
		if result.IsError {
			t.Errorf("write failed for %s: %s", name, result.Content)
		}
	}

	for name, expectedContent := range files {
		filePath := filepath.Join(tmpDir, name)
		input := map[string]interface{}{
			"file_path": filePath,
		}

		result, err := readTool.Execute(context.Background(), input, tc)
		testutil.AssertNoError(t, err)
		if result.IsError {
			t.Errorf("read failed for %s: %s", name, result.Content)
		}

		if contentStr, ok := result.Content.(string); ok {
			testutil.AssertContains(t, contentStr, expectedContent)
		}
	}
}

func TestPermissionModes(t *testing.T) {
	tests := []struct {
		name             string
		mode             permission.Mode
		toolName         string
		input            map[string]interface{}
		expectedBehavior permission.Behavior
	}{
		{
			name:             "Accept mode allows everything",
			mode:             permission.ModeAccept,
			toolName:         "Bash",
			input:            map[string]interface{}{"command": "rm -rf /"},
			expectedBehavior: permission.BehaviorAllow,
		},
		{
			name:             "Default mode asks for Bash",
			mode:             permission.ModeDefault,
			toolName:         "Bash",
			input:            map[string]interface{}{"command": "ls"},
			expectedBehavior: permission.BehaviorAsk,
		},
		{
			name:             "Default mode allows Read",
			mode:             permission.ModeDefault,
			toolName:         "Read",
			input:            map[string]interface{}{"file_path": "/test.txt"},
			expectedBehavior: permission.BehaviorAllow,
		},
		{
			name:             "Auto mode allows safe commands",
			mode:             permission.ModeAuto,
			toolName:         "Bash",
			input:            map[string]interface{}{"command": "ls -la"},
			expectedBehavior: permission.BehaviorAllow,
		},
		{
			name:             "Auto mode asks for dangerous commands",
			mode:             permission.ModeAuto,
			toolName:         "Bash",
			input:            map[string]interface{}{"command": "rm -rf /data"},
			expectedBehavior: permission.BehaviorAsk,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := permission.NewChecker(tt.mode)
			decision := checker.Check(tt.toolName, tt.input)
			testutil.AssertEqual(t, tt.expectedBehavior, decision.Behavior)
		})
	}
}
