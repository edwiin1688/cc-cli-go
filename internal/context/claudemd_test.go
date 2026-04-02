package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindCLAUDEMDFiles_CurrentDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	claudemd := filepath.Join(tmpDir, "CLAUDE.md")

	if err := os.WriteFile(claudemd, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	files := findCLAUDEMDFiles(tmpDir)

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}

	if len(files) > 0 && files[0] != claudemd {
		t.Errorf("expected %s, got %s", claudemd, files[0])
	}
}

func TestFindCLAUDEMDFiles_ParentDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	subdir := filepath.Join(tmpDir, "subdir")

	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	claudemd := filepath.Join(tmpDir, "CLAUDE.md")
	if err := os.WriteFile(claudemd, []byte("# Parent"), 0644); err != nil {
		t.Fatal(err)
	}

	files := findCLAUDEMDFiles(subdir)

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}

	if len(files) > 0 && files[0] != claudemd {
		t.Errorf("expected %s, got %s", claudemd, files[0])
	}
}

func TestFindCLAUDEMDFiles_GEMINI_MD(t *testing.T) {
	tmpDir := t.TempDir()
	geminiMd := filepath.Join(tmpDir, "GEMINI.md")

	if err := os.WriteFile(geminiMd, []byte("# Gemini"), 0644); err != nil {
		t.Fatal(err)
	}

	files := findCLAUDEMDFiles(tmpDir)

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}

	if len(files) > 0 && files[0] != geminiMd {
		t.Errorf("expected %s, got %s", geminiMd, files[0])
	}
}

func TestFindCLAUDEMDFiles_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	subdir := filepath.Join(tmpDir, "project")

	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	parentClaudemd := filepath.Join(tmpDir, "CLAUDE.md")
	childClaudemd := filepath.Join(subdir, "CLAUDE.md")

	if err := os.WriteFile(parentClaudemd, []byte("# Parent"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(childClaudemd, []byte("# Child"), 0644); err != nil {
		t.Fatal(err)
	}

	files := findCLAUDEMDFiles(subdir)

	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestFindCLAUDEMDFiles_NoFiles(t *testing.T) {
	tmpDir := t.TempDir()

	files := findCLAUDEMDFiles(tmpDir)

	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestFindCLAUDEMDFiles_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	level1 := filepath.Join(tmpDir, "level1")
	level2 := filepath.Join(level1, "level2")
	level3 := filepath.Join(level2, "level3")

	if err := os.MkdirAll(level3, 0755); err != nil {
		t.Fatal(err)
	}

	rootClaudemd := filepath.Join(tmpDir, "CLAUDE.md")
	level2Claudemd := filepath.Join(level2, "CLAUDE.md")

	if err := os.WriteFile(rootClaudemd, []byte("# Root"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(level2Claudemd, []byte("# Level2"), 0644); err != nil {
		t.Fatal(err)
	}

	files := findCLAUDEMDFiles(level3)

	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}
