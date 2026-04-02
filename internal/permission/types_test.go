package permission

import (
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
)

func TestChecker_Check_DefaultMode_AllowReadTool(t *testing.T) {
	checker := NewChecker(ModeDefault)

	decision := checker.Check("Read", map[string]interface{}{
		"file_path": "/test.txt",
	})

	testutil.AssertEqual(t, BehaviorAllow, decision.Behavior)
}

func TestChecker_Check_DefaultMode_AskBashTool(t *testing.T) {
	checker := NewChecker(ModeDefault)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "ls",
	})

	testutil.AssertEqual(t, BehaviorAsk, decision.Behavior)
}

func TestChecker_Check_AcceptMode(t *testing.T) {
	checker := NewChecker(ModeAccept)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "rm -rf /",
	})

	testutil.AssertEqual(t, BehaviorAllow, decision.Behavior)
	testutil.AssertContains(t, decision.Reason, "accept mode")
}

func TestChecker_Check_AutoMode_SafeCommand(t *testing.T) {
	checker := NewChecker(ModeAuto)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "ls -la",
	})

	testutil.AssertEqual(t, BehaviorAllow, decision.Behavior)
	testutil.AssertContains(t, decision.Reason, "auto mode")
}

func TestChecker_Check_AutoMode_DangerousCommand(t *testing.T) {
	checker := NewChecker(ModeAuto)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "rm -rf /important",
	})

	testutil.AssertEqual(t, BehaviorAsk, decision.Behavior)
}

func TestChecker_Check_DangerousCommand_RmRf(t *testing.T) {
	checker := NewChecker(ModeDefault)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "rm -rf /data",
	})

	testutil.AssertEqual(t, BehaviorAsk, decision.Behavior)
	testutil.AssertContains(t, decision.Reason, "dangerous")
}

func TestChecker_Check_DangerousCommand_DropTable(t *testing.T) {
	checker := NewChecker(ModeDefault)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "echo 'DROP TABLE users' | mysql",
	})

	testutil.AssertEqual(t, BehaviorAsk, decision.Behavior)
}

func TestChecker_Check_DangerousCommand_GitForcePush(t *testing.T) {
	checker := NewChecker(ModeDefault)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "git push --force origin master",
	})

	testutil.AssertEqual(t, BehaviorAsk, decision.Behavior)
}

func TestChecker_SetRules(t *testing.T) {
	checker := NewChecker(ModeDefault)

	rules := []Rule{
		{ToolName: "Bash", Pattern: "ls", Behavior: BehaviorAllow},
	}
	checker.SetRules(rules)

	decision := checker.Check("Bash", map[string]interface{}{
		"command": "ls -la",
	})

	testutil.AssertEqual(t, BehaviorAllow, decision.Behavior)
}

func TestChecker_Check_NonBashTool_NotDangerous(t *testing.T) {
	checker := NewChecker(ModeDefault)

	decision := checker.Check("Write", map[string]interface{}{
		"file_path": "/test.txt",
		"content":   "test",
	})

	testutil.AssertEqual(t, BehaviorAsk, decision.Behavior)
}

func TestIsDangerousCommand_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"rm -rf", "rm -rf /data", true},
		{"rm -fr", "rm -fr /home", true},
		{"DROP TABLE", "DROP TABLE users;", true},
		{"DROP DATABASE", "DROP DATABASE test;", true},
		{"DELETE FROM", "DELETE FROM users;", true},
		{"git push --force", "git push --force origin master", true},
		{"git push -f", "git push -f", true},
		{"git reset --hard", "git reset --hard HEAD~1", true},
		{"safe command ls", "ls -la", false},
		{"safe command cat", "cat file.txt", false},
		{"safe command echo", "echo 'hello'", false},
		{"safe command git status", "git status", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDangerousCommand("Bash", map[string]interface{}{
				"command": tt.command,
			})
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestGetDangerousReason(t *testing.T) {
	tests := []struct {
		name             string
		command          string
		expectedInReason string
	}{
		{"rm -rf", "rm -rf /data", "deletion"},
		{"DROP TABLE", "DROP TABLE users", "SQL"},
		{"git push --force", "git push --force", "force push"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reason := getDangerousReason(tt.command)
			testutil.AssertContains(t, reason, tt.expectedInReason)
		})
	}
}
