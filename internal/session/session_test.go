package session

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
	"github.com/liao-eli/cc-cli-go/internal/types"
)

func TestNewSession(t *testing.T) {
	projectID := "/test/project"

	session := NewSession(projectID)

	if session.ID == "" {
		t.Error("expected session ID to be generated")
	}
	testutil.AssertEqual(t, projectID, session.ProjectID)
	testutil.AssertEqual(t, 0, len(session.Messages))
}

func TestSession_AddMessage(t *testing.T) {
	session := NewSession("/test")
	msg := types.NewUserMessage("test message")

	session.AddMessage(msg)

	testutil.AssertEqual(t, 1, len(session.Messages))
	testutil.AssertEqual(t, msg, session.Messages[0])
}

func TestSession_SaveAndLoad(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	session := NewSession("/test")
	session.AddMessage(types.NewUserMessage("hello"))
	session.AddMessage(types.NewAssistantMessage())

	err := session.Save()
	testutil.AssertNoError(t, err)

	loaded, err := LoadSession(session.ID)
	testutil.AssertNoError(t, err)

	testutil.AssertEqual(t, session.ID, loaded.ID)
	testutil.AssertEqual(t, session.ProjectID, loaded.ProjectID)
	testutil.AssertEqual(t, 2, len(loaded.Messages))
}

func TestSession_Save_CreatesDirectory(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	session := NewSession("/test")

	err := session.Save()
	testutil.AssertNoError(t, err)

	sessionDir := filepath.Join(tmpDir, ".claude", "sessions")
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		t.Error("expected sessions directory to be created")
	}
}

func TestLoadSession_NotFound(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	_, err := LoadSession("nonexistent")
	testutil.AssertError(t, err)
}

func TestGetLastSession_NoSessions(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	_, err := GetLastSession()
	testutil.AssertError(t, err)
}

func TestGetLastSession_WithSessions(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	session1 := NewSession("/test1")
	session1.Save()

	session2 := NewSession("/test2")
	session2.Save()

	last, err := GetLastSession()
	testutil.AssertNoError(t, err)

	testutil.AssertEqual(t, session2.ID, last.ID)
}

func TestCleanupOldSessions(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	session := NewSession("/test")
	session.Save()

	err := CleanupOldSessions(0)
	testutil.AssertNoError(t, err)

	sessionDir := filepath.Join(tmpDir, ".claude", "sessions")
	files, _ := os.ReadDir(sessionDir)
	if len(files) > 0 {
		t.Error("expected all sessions to be cleaned up")
	}
}

func TestGenerateUUID(t *testing.T) {
	id := generateUUID()

	if id == "" {
		t.Error("expected non-empty UUID")
	}

	if len(id) == 0 {
		t.Error("expected UUID with length > 0")
	}
}

func TestSession_UpdatedAt(t *testing.T) {
	session := NewSession("/test")
	originalTime := session.UpdatedAt

	session.AddMessage(types.NewUserMessage("test"))

	if !session.UpdatedAt.After(originalTime) {
		t.Error("expected UpdatedAt to be updated after adding message")
	}
}
