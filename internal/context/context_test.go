package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
)

func TestBuildContext_Success(t *testing.T) {
	info, err := BuildContext()

	testutil.AssertNoError(t, err)

	if info == nil {
		t.Fatal("expected non-nil context info")
	}

	if info.WorkingDir == "" {
		t.Error("expected non-empty working directory")
	}

	if info.DateTime == "" {
		t.Error("expected non-empty date time")
	}
}

func TestBuildContext_WorkingDirectory(t *testing.T) {
	info, err := BuildContext()

	testutil.AssertNoError(t, err)

	expectedDir, _ := os.Getwd()
	testutil.AssertEqual(t, expectedDir, info.WorkingDir)
}

func TestContextInfo_ToSystemPrompt(t *testing.T) {
	info := &ContextInfo{
		WorkingDir: "/test/project",
		GitBranch:  "main",
		GitStatus:  "M file.txt",
		DateTime:   "2026-04-02 10:00:00",
	}

	prompt := info.ToSystemPrompt()

	testutil.AssertContains(t, prompt, "Environment Information")
	testutil.AssertContains(t, prompt, "/test/project")
	testutil.AssertContains(t, prompt, "main")
	testutil.AssertContains(t, prompt, "M file.txt")
	testutil.AssertContains(t, prompt, "2026-04-02 10:00:00")
}

func TestContextInfo_ToSystemPrompt_WithCLAUDEMD(t *testing.T) {
	tmpDir := t.TempDir()
	claudemd := filepath.Join(tmpDir, "CLAUDE.md")
	content := "# Test Project\nThis is a test project."

	if err := os.WriteFile(claudemd, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	info := &ContextInfo{
		WorkingDir:    tmpDir,
		DateTime:      "2026-04-02 10:00:00",
		CLAUDEMDFiles: []string{claudemd},
	}

	prompt := info.ToSystemPrompt()

	testutil.AssertContains(t, prompt, "CLAUDE.md files found")
	testutil.AssertContains(t, prompt, "Test Project")
}

func TestContextInfo_ToSystemPrompt_NoGit(t *testing.T) {
	info := &ContextInfo{
		WorkingDir: "/test/project",
		DateTime:   "2026-04-02 10:00:00",
	}

	prompt := info.ToSystemPrompt()

	testutil.AssertContains(t, prompt, "Working Directory: /test/project")
	testutil.AssertContains(t, prompt, "Date/Time: 2026-04-02 10:00:00")

	if strings.Contains(prompt, "Git Branch") {
		t.Error("should not contain Git Branch when empty")
	}
}

func TestContextInfo_ToSystemPrompt_EmptyGitStatus(t *testing.T) {
	info := &ContextInfo{
		WorkingDir: "/test/project",
		GitBranch:  "main",
		GitStatus:  "",
		DateTime:   "2026-04-02 10:00:00",
	}

	prompt := info.ToSystemPrompt()

	testutil.AssertContains(t, prompt, "Git Branch: main")

	if strings.Contains(prompt, "Git Status") {
		t.Error("should not contain Git Status when empty")
	}
}

func TestGetGitBranch_InGitRepo(t *testing.T) {
	branch := getGitBranch()

	if branch == "" {
		t.Log("Not in a git repository or no branch detected")
	}
}

func TestGetGitStatus_InGitRepo(t *testing.T) {
	status := getGitStatus()

	t.Logf("Git status: %s", status)
}
