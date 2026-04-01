# Go Module Structure - Claude Code Go Port

> **目的**: 定義完整的 Go module 目錄結構

---

## 目錄結構總覽

```
claude-code-go/
├── cmd/                              # 應用程式入口點
│   └── claude-code/
│       └── main.go                   # 主程式入口
│
├── internal/                         # 內部套件（不可被外部引用）
│   ├── api/                          # Anthropic API 客戶端
│   │   ├── client.go                 # HTTP 客戶端
│   │   ├── streaming.go              # SSE 串流處理
│   │   ├── events.go                 # 事件解析
│   │   ├── request.go                # 請求建構
│   │   ├── config.go                 # API 配置
│   │   └── errors.go                 # API 錯誤定義
│   │
│   ├── cli/                          # CLI 框架
│   │   ├── root.go                   # 根命令
│   │   ├── run.go                    # 互動模式命令
│   │   ├── print.go                  # -p 模式命令
│   │   ├── resume.go                 # --resume 處理
│   │   └── flags.go                  # 旗標定義
│   │
│   ├── tui/                          # Bubble Tea TUI
│   │   ├── app.go                    # 主應用模型
│   │   ├── model.go                  # 模型定義
│   │   ├── messages.go               # 訊息渲染
│   │   ├── input.go                  # 輸入元件
│   │   ├── spinner.go                # 載入動畫
│   │   ├── dialog/                   # 對話框
│   │   │   ├── permission.go         # 權限對話框
│   │   │   └── confirm.go            # 確認對話框
│   │   └── theme/                    # 主題與顏色
│   │       └── theme.go
│   │
│   ├── query/                        # 查詢引擎
│   │   ├── engine.go                 # 查詢循環
│   │   ├── types.go                  # 類型定義
│   │   ├── context.go                # 上下文建構
│   │   ├── messages.go               # 訊息管理
│   │   └── compact.go                # 上下文壓縮（可選）
│   │
│   ├── tools/                        # 工具系統
│   │   ├── tool.go                   # 工具介面
│   │   ├── registry.go               # 工具註冊表
│   │   ├── executor.go               # 工具執行器
│   │   ├── bash/                     # Bash 工具
│   │   │   └── bash.go
│   │   ├── read/                     # 檔案讀取工具
│   │   │   └── read.go
│   │   ├── edit/                     # 檔案編輯工具
│   │   │   └── edit.go
│   │   ├── write/                    # 檔案寫入工具
│   │   │   └── write.go
│   │   ├── glob/                     # 檔案模式匹配工具
│   │   │   └── glob.go
│   │   └── grep/                     # 檔案內容搜尋工具
│   │       └── grep.go
│   │
│   ├── permission/                   # 權限系統
│   │   ├── rules.go                  # 規則定義
│   │   ├── matcher.go                # 模式匹配
│   │   ├── modes.go                  # 權限模式
│   │   └── manager.go                # 規則管理器
│   │
│   ├── session/                      # 會話管理
│   │   ├── storage.go                # 會話儲存
│   │   ├── transcript.go             # JSONL 寫入器
│   │   ├── metadata.go               # 會話元資料
│   │   └── resume.go                 # 恢復處理器
│   │
│   ├── context/                      # 上下文建構
│   │   ├── git.go                    # Git 狀態/分支
│   │   ├── claudemd.go               # CLAUDE.md 發現
│   │   └── system.go                 # 系統上下文
│   │
│   ├── config/                       # 配置管理
│   │   ├── settings.go               # 設定管理
│   │   ├── paths.go                  # 路徑解析
│   │   └── schema.go                 # 設定 Schema
│   │
│   └── types/                        # 共享類型
│       ├── message.go                # 訊息類型
│       ├── content.go                # 內容區塊類型
│       ├── usage.go                  # Token 使用量
│       ├── permission.go             # 權限類型
│       └── session.go                # 會話類型
│
├── pkg/                              # 公開套件（可被外部引用）
│   └── (目前為空，預留未來擴展)
│
├── scripts/                          # 建構腳本
│   ├── build.sh                      # 建構腳本
│   └── release.sh                    # 發布腳本
│
├── .github/                          # GitHub 配置
│   └── workflows/
│       └── release.yml               # 自動發布
│
├── docs/                             # 文檔
│   ├── README.md                     # 專案說明
│   ├── ARCHITECTURE.md               # 架構說明
│   └── CONTRIBUTING.md               # 貢獻指南
│
├── go.mod                            # Go module 定義
├── go.sum                            # 依賴校驗和
├── Makefile                          # Make 建構
├── .gitignore                        # Git 忽略規則
├── .goreleaser.yml                   # GoReleaser 配置
└── README.md                         # 專案說明
```

---

## 套件依賴關係圖

```
cmd/claude-code
    │
    └── internal/cli
            │
            ├── internal/config
            ├── internal/session
            └── internal/tui
                    │
                    ├── internal/query
                    │       │
                    │       ├── internal/api
                    │       ├── internal/tools
                    │       ├── internal/permission
                    │       ├── internal/context
                    │       └── internal/types
                    │
                    └── internal/types

Key:
  → A depends on B
  internal/* packages are private
  pkg/* packages are public (currently empty)
```

---

## 各套件詳細說明

### cmd/claude-code/

**職責**: 應用程式入口點

**檔案**:
- `main.go`: 初始化 cobra，執行 CLI

**依賴**: `internal/cli`

---

### internal/api/

**職責**: Anthropic API 客戶端

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `client.go` | HTTP 客戶端 | `Client`, `NewClient()` |
| `streaming.go` | SSE 串流 | `Stream()` |
| `events.go` | 事件解析 | `ParseMessageStart()`, `ParseContentBlock()`, `ParseDelta()` |
| `request.go` | 請求建構 | `Request`, `NewRequest()`, `ToolParam` |
| `config.go` | API 配置 | `DefaultModel`, `DefaultMaxTokens` |
| `errors.go` | 錯誤定義 | `APIError` |

**依賴**: `internal/types`

**使用範例**:
```go
client := api.NewClient(apiKey)
req := api.NewRequest("claude-sonnet-4-20250514", 4096)
events, _ := client.Stream(ctx, req)
```

---

### internal/cli/

**職責**: CLI 命令定義

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `root.go` | 根命令 | `Execute()` |
| `run.go` | 互動模式 | `runCmd` |
| `print.go` | 非互動模式 | `printCmd` |
| `resume.go` | 會話恢復 | `resumeCmd` |
| `flags.go` | 旗標定義 | 全域旗標變數 |

**依賴**: `internal/api`, `internal/tui`, `internal/query`, `internal/tools`, `internal/session`

**使用範例**:
```bash
claude-code run              # 互動模式
claude-code -p "hello"       # 非互動模式
claude-code --resume <uuid>  # 恢復會話
```

---

### internal/tui/

**職責**: Bubble Tea TUI

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `app.go` | 主應用 | `Model`, `Update()`, `View()` |
| `model.go` | 模型定義 | `InitialModel()` |
| `messages.go` | 訊息渲染 | `renderMessage()` |
| `input.go` | 輸入元件 | `InputModel` |
| `spinner.go` | 載入動畫 | `SpinnerModel` |
| `dialog/permission.go` | 權限對話框 | `PermissionDialog` |
| `dialog/confirm.go` | 確認對話框 | `ConfirmDialog` |
| `theme/theme.go` | 主題定義 | `Theme`, `DefaultTheme()` |

**依賴**: `internal/query`, `internal/types`, `github.com/charmbracelet/bubbletea`

**使用範例**:
```go
model := tui.InitialModel()
p := tea.NewProgram(model)
p.Run()
```

---

### internal/query/

**職責**: 查詢引擎（核心循環）

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `engine.go` | 查詢循環 | `Engine`, `NewEngine()`, `Query()` |
| `types.go` | 類型定義 | `QueryParams`, `QueryResult`, `StreamEvent` |
| `context.go` | 上下文建構 | `BuildContext()` |
| `messages.go` | 訊息管理 | `MessageManager` |
| `compact.go` | 上下文壓縮 | `Compact()` (可選) |

**依賴**: `internal/api`, `internal/tools`, `internal/types`

**使用範例**:
```go
engine := query.NewEngine(client, toolReg)
events, results := engine.Query(ctx, params)
```

---

### internal/tools/

**職責**: 工具系統

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `tool.go` | 工具介面 | `Tool`, `ToolResult`, `ToolContext` |
| `registry.go` | 工具註冊表 | `Registry`, `NewRegistry()` |
| `executor.go` | 工具執行器 | `Executor` |
| `bash/bash.go` | Bash 工具 | `New()` |
| `read/read.go` | 讀取工具 | `New()` |
| `edit/edit.go` | 編輯工具 | `New()` |
| `write/write.go` | 寫入工具 | `New()` |
| `glob/glob.go` | Glob 工具 | `New()` |
| `grep/grep.go` | Grep 工具 | `New()` |

**依賴**: `internal/types`

**使用範例**:
```go
reg := tools.NewRegistry()
reg.Register(bash.New())
reg.Register(read.New())

tool := reg.Get("Bash")
result, _ := tool.Execute(ctx, input, tc)
```

---

### internal/permission/

**職責**: 權限系統

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `rules.go` | 規則定義 | `Rule`, `Rules` |
| `matcher.go` | 模式匹配 | `Matcher`, `Match()` |
| `modes.go` | 權限模式 | `PermissionMode` |
| `manager.go` | 規則管理器 | `Manager`, `NewManager()` |

**依賴**: `internal/types`

**使用範例**:
```go
manager := permission.NewManager()
result := manager.Check(tool, input)
if result.Behavior == "ask" {
    decision := AskUser(tool, input)
}
```

---

### internal/session/

**職責**: 會話管理

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `storage.go` | 會話儲存 | `Storage`, `NewStorage()` |
| `transcript.go` | JSONL 寫入 | `TranscriptWriter` |
| `metadata.go` | 會話元資料 | `Metadata` |
| `resume.go` | 恢復處理 | `Resume()` |

**依賴**: `internal/types`

**使用範例**:
```go
storage := session.NewStorage()
storage.Save(sessionID, messages)
messages := storage.Load(sessionID)
```

---

### internal/context/

**職責**: 上下文建構

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `git.go` | Git 狀態 | `GetGitStatus()`, `GetBranch()` |
| `claudemd.go` | CLAUDE.md 發現 | `DiscoverClaudeMD()`, `ReadClaudeMD()` |
| `system.go` | 系統上下文 | `GetSystemContext()` |

**依賴**: 無

**使用範例**:
```go
gitStatus, _ := context.GetGitStatus()
claudemd := context.DiscoverClaudeMD()
```

---

### internal/config/

**職責**: 配置管理

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `settings.go` | 設定管理 | `Settings`, `LoadSettings()` |
| `paths.go` | 路徑解析 | `GetConfigPath()`, `GetSessionPath()` |
| `schema.go` | 設定 Schema | `SettingsSchema` |

**依賴**: 無

**使用範例**:
```go
settings, _ := config.LoadSettings()
model := settings.Model
```

---

### internal/types/

**職責**: 共享類型定義

**檔案**:

| 檔案 | 功能 | 匯出的類型/函數 |
|------|------|----------------|
| `message.go` | 訊息類型 | `Message`, `MessageType` |
| `content.go` | 內容區塊 | `ContentBlock` |
| `usage.go` | Token 使用量 | `Usage` |
| `permission.go` | 權限類型 | `PermissionResult` |
| `session.go` | 會話類型 | `SessionInfo` |

**依賴**: 無

---

## go.mod 範例

```go
module github.com/yourusername/claude-code-go

go 1.21

require (
    github.com/spf13/cobra v1.8.0
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.17.1
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/google/uuid v1.6.0
)
```

---

## Makefile 範例

```makefile
.PHONY: build test clean install

VERSION := $(shell git describe --tags --always)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/claude-code ./cmd/claude-code

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/
	go clean

install: build
	cp bin/claude-code /usr/local/bin/

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

dev:
	go run ./cmd/claude-code run

cross-build:
	GOOS=darwin GOARCH=arm64 go build -o bin/claude-code-darwin-arm64 ./cmd/claude-code
	GOOS=darwin GOARCH=amd64 go build -o bin/claude-code-darwin-amd64 ./cmd/claude-code
	GOOS=linux GOARCH=amd64 go build -o bin/claude-code-linux-amd64 ./cmd/claude-code
	GOOS=windows GOARCH=amd64 go build -o bin/claude-code-windows-amd64.exe ./cmd/claude-code
```

---

## 測試結構

```
claude-code-go/
├── internal/
│   ├── api/
│   │   ├── client_test.go
│   │   ├── streaming_test.go
│   │   └── events_test.go
│   │
│   ├── tools/
│   │   ├── bash/
│   │   │   └── bash_test.go
│   │   ├── read/
│   │   │   └── read_test.go
│   │   └── edit/
│   │       └── edit_test.go
│   │
│   └── query/
│       └── engine_test.go
│
└── testdata/
    ├── test_file.txt
    └── session.jsonl
```

---

## 文檔結構

```
docs/
├── README.md              # 專案說明
├── ARCHITECTURE.md        # 架構說明
├── CONTRIBUTING.md        # 貢獻指南
├── CODE_OF_CONDUCT.md     # 行為準則
└── API.md                 # API 文檔
```

---

## 配置檔案位置

```
~/.claude/
├── settings.json          # 全域設定
├── credentials.json       # API 金鑰（可選）
├── sessions/              # 會話歷史
│   ├── <uuid>.jsonl
│   └── <uuid>.metadata.json
└── projects/
    └── <cwd-hash>/
        ├── settings.json  # 專案設定
        └── sessions/

<project>/.claude/
└── settings.json          # 專案設定（可提交至 git）
```

---

## 預期檔案數量

| 類型 | 數量 |
|------|------|
| Go 原始碼 | ~30-40 個 |
| 測試檔案 | ~10-15 個 |
| 文檔 | ~5 個 |
| 配置 | ~5 個 |
| **總計** | **~50-65 個檔案** |

---

## 預期程式碼量

| 套件 | 預估 LOC |
|------|---------|
| `internal/api` | ~500-700 |
| `internal/cli` | ~300-400 |
| `internal/tui` | ~800-1000 |
| `internal/query` | ~500-600 |
| `internal/tools` | ~800-1000 |
| `internal/permission` | ~300-400 |
| `internal/session` | ~200-300 |
| `internal/context` | ~150-200 |
| `internal/config` | ~150-200 |
| `internal/types` | ~200-300 |
| `cmd/claude-code` | ~50-100 |
| **總計** | **~4000-5200 LOC** |

---

*Document generated for Go porting reference*