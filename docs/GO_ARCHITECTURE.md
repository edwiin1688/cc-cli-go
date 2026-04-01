# Go Architecture - Claude Code Go Port

> **目的**: 定義 Go 版本的架構設計

---

## 設計原則

1. **KISS**: 保持簡潔，避免過度設計
2. **Idiomatic Go**: 遵循 Go 慣例，不模仿 TS/JS 模式
3. **Explicit over Implicit**: 明確優於隱式
4. **Composition over Inheritance**: 組合優於繼承
5. **Error as Values**: 錯誤作為值處理

---

## 1. 整體架構

### 1.1 架構圖

```
┌─────────────────────────────────────────────────────────────┐
│                          CLI Entry                           │
│                      (cmd/claude-code/main.go)               │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      CLI Framework                           │
│                      (internal/cli/)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   cobra     │  │   flags     │  │   config    │         │
│  │   commands  │  │   parsing   │  │   loading   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                        TUI Layer                             │
│                    (internal/tui/)                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Bubble Tea  │  │  Messages   │  │   Dialogs   │         │
│  │    Model    │  │  Rendering  │  │   (perms)   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      Query Engine                            │
│                   (internal/query/)                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Query     │  │   Message   │  │   Context   │         │
│  │    Loop     │  │  Manager    │  │   Builder   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Client                              │
│                    (internal/api/)                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Anthropic │  │  Streaming  │  │   Events    │         │
│  │    Client   │  │    SSE      │  │   Parser    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      Tool System                             │
│                    (internal/tools/)                         │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐    │
│  │  Bash  │ │  Read  │ │  Edit  │ │  Glob  │ │  Grep  │    │
│  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘    │
│  ┌────────┐ ┌────────┐ ┌────────┐                          │
│  │ Write  │ │ Tool   │ │ Exec   │                          │
│  │        │ │Registry│ │ Pool   │                          │
│  └────────┘ └────────┘ └────────┘                          │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                    Permission System                         │
│                   (internal/permission/)                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │    Rules    │  │   Matcher   │  │   Dialog    │         │
│  │   Manager   │  │  (pattern)  │  │   Handler   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                    Session Storage                           │
│                    (internal/session/)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Storage   │  │  Transcript │  │   Resume    │         │
│  │   Manager   │  │   Writer    │  │   Handler   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Data Flow

```
User Input
    │
    ▼
┌──────────────┐
│  TUI (Bubble │
│     Tea)     │
└──────────────┘
    │
    │ tea.Msg (UserInputMsg)
    ▼
┌──────────────┐     ┌──────────────┐
│   Query      │────▶│    API       │
│   Engine     │     │   Client     │
└──────────────┘     └──────────────┘
    │                       │
    │                       │ SSE Events
    │                       ▼
    │                 ┌──────────────┐
    │                 │   Event      │
    │                 │   Parser     │
    │                 └──────────────┘
    │                       │
    │ ◀─────────────────────┘
    │
    │ ToolUse detected
    ▼
┌──────────────┐     ┌──────────────┐
│    Tool      │────▶│  Permission  │
│   Executor   │     │   Check      │
└──────────────┘     └──────────────┘
    │                       │
    │                       │ Ask user
    │                       ▼
    │                 ┌──────────────┐
    │                 │   TUI        │
    │                 │   Dialog     │
    │                 └──────────────┘
    │                       │
    │ ◀─────────────────────┘
    │
    │ Execute tool
    ▼
┌──────────────┐
│    Tool      │
│  Implementation
└──────────────┘
    │
    │ ToolResult
    ▼
┌──────────────┐
│   Message    │
│   Manager    │
└──────────────┘
    │
    │ Continue loop
    ▼
[Back to Query Engine]
```

---

## 2. 核心型別定義

### 2.1 Message Types

```go
// internal/types/message.go

package types

type MessageType string

const (
    MessageTypeUser      MessageType = "user"
    MessageTypeAssistant MessageType = "assistant"
    MessageTypeSystem    MessageType = "system"
    MessageTypeToolUse   MessageType = "tool_use"
    MessageTypeToolResult MessageType = "tool_result"
)

type ContentBlock struct {
    Type    string      `json:"type"`
    Text    string      `json:"text,omitempty"`
    Thinking string     `json:"thinking,omitempty"`
    
    // For tool_use
    ID       string      `json:"id,omitempty"`
    Name     string      `json:"name,omitempty"`
    Input    interface{} `json:"input,omitempty"`
    
    // For tool_result
    ToolUseID string     `json:"tool_use_id,omitempty"`
    Content   interface{} `json:"content,omitempty"`
    IsError   bool       `json:"is_error,omitempty"`
}

type Message struct {
    UUID      string         `json:"uuid"`
    Type      MessageType    `json:"type"`
    Role      string         `json:"role,omitempty"`
    Content   []ContentBlock `json:"content,omitempty"`
    Usage     *Usage         `json:"usage,omitempty"`
    Timestamp time.Time      `json:"timestamp"`
    
    // For assistant messages
    StopReason string        `json:"stop_reason,omitempty"`
    
    // Metadata
    IsAPIMessage bool        `json:"is_api_message,omitempty"`
    APIError     string      `json:"api_error,omitempty"`
}

type Usage struct {
    InputTokens       int `json:"input_tokens"`
    OutputTokens      int `json:"output_tokens"`
    CacheCreation     int `json:"cache_creation_input_tokens,omitempty"`
    CacheRead         int `json:"cache_read_input_tokens,omitempty"`
}
```

### 2.2 Tool Types

```go
// internal/tools/tool.go

package tools

type ToolResult struct {
    Content     interface{}
    IsError     bool
    NewMessages []*types.Message
}

type PermissionResult struct {
    Behavior     string //"allow", "deny", "ask"
    UpdatedInput map[string]interface{}
    Reason       string
}

type ToolContext struct {
    WorkingDir      string
    AbortSignal     context.Context
    AppState        *AppState
    PermissionCtx   *PermissionContext
    AddNotification func(Notification)
    
    // For streaming progress
    OnProgress      func(Progress)
}

type Tool interface {
    // Identity
    Name() string
    Aliases() []string
    Description() string
    SearchHint() string
    
    // Schema
    InputSchema() map[string]interface{}
    OutputSchema() map[string]interface{}
    
    // Execution
    Execute(ctx context.Context, input map[string]interface{}, tc *ToolContext) (*ToolResult, error)
    
    // Behavior
    IsEnabled() bool
    IsReadOnly(input map[string]interface{}) bool
    IsConcurrencySafe(input map[string]interface{}) bool
    IsDestructive(input map[string]interface{}) bool
    InterruptBehavior() InterruptBehavior
    
    // Permissions
    CheckPermissions(input map[string]interface{}, tc *ToolContext) PermissionResult
    ValidateInput(input map[string]interface{}, tc *ToolContext) error
    
    // Serialization
    MapResultToBlock(result *ToolResult, toolUseID string) *types.ContentBlock
    
    // Display
    UserFacingName(input map[string]interface{}) string
    ActivityDescription(input map[string]interface{}) string
}

type InterruptBehavior string

const (
    InterruptCancel InterruptBehavior = "cancel"
    InterruptBlock  InterruptBehavior = "block"
)
```

### 2.3 Query Types

```go
// internal/query/types.go

package query

type QueryParams struct {
    Messages         []*types.Message
    SystemPrompt     []string
    UserContext      map[string]string
    SystemContext    map[string]string
    Tools            []tools.Tool
    Model            string
    MaxTokens        int
    MaxTurns         int
    AbortSignal      context.Context
}

type QueryResult struct {
    Terminal   Terminal
    Messages   []*types.Message
    Usage      *types.Usage
}

type Terminal struct {
    Reason string //"completed", "aborted", "error", "max_turns"
    Error  error
}

type StreamEvent struct {
    Type    string
    Message *types.Message
    Delta   *ContentDelta
    Usage   *types.Usage
}

type ContentDelta struct {
    Type  string //"text_delta", "input_json_delta", "thinking_delta"
    Text  string
    Value string //for partial JSON
}
```

---

## 3. Package 結構

### 3.1 目錄結構

```
claude-code-go/
├── cmd/
│   └── claude-code/
│       └── main.go              # Entry point
│
├── internal/
│   ├── api/                     # Anthropic API client
│   │   ├── client.go            # HTTP client setup
│   │   ├── streaming.go         # SSE streaming
│   │   ├── events.go            # Event parsing
│   │   ├── request.go           # Request building
│   │   └── errors.go            # API errors
│   │
│   ├── cli/                     # CLI framework
│   │   ├── root.go              # Root command
│   │   ├── print.go             # -p mode
│   │   ├── resume.go            # --resume handling
│   │   └── flags.go             # Flag definitions
│   │
│   ├── tui/                     # Bubble Tea TUI
│   │   ├── app.go               # Main app model
│   │   ├── messages.go          # Message list component
│   │   ├── input.go             # Input component
│   │   ├── spinner.go           # Loading spinner
│   │   ├── dialog/              # Dialogs
│   │   │   ├── permission.go    # Permission dialog
│   │   │   └── confirm.go       # Confirmation dialog
│   │   └── theme/               # Theme/colors
│   │       └── theme.go
│   │
│   ├── query/                   # Query engine
│   │   ├── engine.go            # Query loop
│   │   ├── context.go           # Context building
│   │   ├── messages.go          # Message management
│   │   └── compact.go           # Context compaction (optional)
│   │
│   ├── tools/                   # Tool system
│   │   ├── tool.go              # Tool interface
│   │   ├── registry.go          # Tool registry
│   │   ├── executor.go          # Tool executor
│   │   ├── bash/                # Bash tool
│   │   │   └── bash.go
│   │   ├── read/                # File read tool
│   │   │   └── read.go
│   │   ├── edit/                # File edit tool
│   │   │   └── edit.go
│   │   ├── write/               # File write tool
│   │   │   └── write.go
│   │   ├── glob/                # Glob tool
│   │   │   └── glob.go
│   │   └── grep/                # Grep tool
│   │       └── grep.go
│   │
│   ├── permission/              # Permission system
│   │   ├── rules.go             # Rule definitions
│   │   ├── matcher.go           # Pattern matching
│   │   ├── modes.go             # Permission modes
│   │   └── manager.go           # Rule manager
│   │
│   ├── session/                 # Session management
│   │   ├── storage.go           # Session storage
│   │   ├── transcript.go        # JSONL writer
│   │   ├── metadata.go          # Session metadata
│   │   └── resume.go            # Resume handler
│   │
│   ├── context/                 # Context building
│   │   ├── git.go               # Git status/branch
│   │   ├── claudemd.go          # CLAUDE.md discovery
│   │   └── system.go            # System context
│   │
│   ├── config/                  # Configuration
│   │   ├── settings.go          # Settings management
│   │   ├── paths.go             # Path resolution
│   │   └── schema.go            # Settings schema
│   │
│   └── types/                   # Shared types
│       ├── message.go           # Message types
│       ├── tool.go              # Tool types
│       ├── permission.go        # Permission types
│       └── session.go           # Session types
│
├── pkg/                         # Public packages (if any)
│   └── ...
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 3.2 Package 依賴關係

```
cmd/claude-code
    │
    ├── internal/cli
    │       │
    │       ├── internal/config
    │       ├── internal/session
    │       └── internal/tui
    │               │
    │               ├── internal/query
    │               │       │
    │               │       ├── internal/api
    │               │       ├── internal/tools
    │               │       ├── internal/permission
    │               │       ├── internal/session
    │               │       ├── internal/context
    │               │       └── internal/types
    │               │
    │               └── internal/types
    │
    └── internal/types
```

---

## 4. 核心流程設計

### 4.1 Query Loop (Go 版本)

```go
// internal/query/engine.go

package query

type Engine struct {
    client     *api.Client
    tools      *tools.Registry
    permission *permission.Manager
    sessions   *session.Storage
    
    // Channels for streaming
    eventChan  chan StreamEvent
    resultChan chan QueryResult
}

func (e *Engine) Query(ctx context.Context, params QueryParams) (<-chan StreamEvent, <-chan QueryResult) {
    events := make(chan StreamEvent, 100)
    results := make(chan QueryResult, 1)
    
    go func() {
        defer close(events)
        defer close(results)
        
        for {
            // 1. Build request
            req := e.buildRequest(params)
            
            // 2. Stream API response
            stream, err := e.client.Stream(ctx, req)
            if err != nil {
                results <- QueryResult{Terminal: Terminal{Reason: "error", Error: err}}
                return
            }
            
            // 3. Process events
            var toolUses []ToolUse
            for event := range stream {
                events <- event
                
                if event.Type == "content_block_start" && event.Message.Content[0].Type == "tool_use" {
                    toolUses = append(toolUses, event.Message.Content[0])
                }
            }
            
            // 4. Execute tools if any
            if len(toolUses) > 0 {
                results := e.executeTools(ctx, toolUses, params)
                params.Messages = append(params.Messages, results...)
                continue // Loop back for next turn
            }
            
            // 5. No tools - we're done
            results <- QueryResult{Terminal: Terminal{Reason: "completed"}}
            return
        }
    }()
    
    return events, results
}

func (e *Engine) executeTools(ctx context.Context, toolUses []ToolUse, params QueryParams) []*types.Message {
    var results []*types.Message
    
    // Parallel tool execution
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    for _, tu := range toolUses {
        wg.Add(1)
        go func(toolUse ToolUse) {
            defer wg.Done()
            
            tool := e.tools.Get(toolUse.Name)
            if tool == nil {
                mu.Lock()
                results = append(results, e.createToolError(toolUse, "tool not found"))
                mu.Unlock()
                return
            }
            
            // Permission check
            perm := tool.CheckPermissions(toolUse.Input, &ToolContext{
                PermissionCtx: params.PermissionContext,
            })
            
            if perm.Behavior == "ask" {
                // Send permission request to UI
                decision := e.permission.AskUser(toolUse, tool)
                perm.Behavior = decision
            }
            
            if perm.Behavior == "deny" {
                mu.Lock()
                results = append(results, e.createToolDeny(toolUse, perm.Reason))
                mu.Unlock()
                return
            }
            
            // Execute
            result, err := tool.Execute(ctx, toolUse.Input, &ToolContext{
                AbortSignal: params.AbortSignal,
            })
            
            if err != nil {
                mu.Lock()
                results = append(results, e.createToolError(toolUse, err.Error()))
                mu.Unlock()
                return
            }
            
            mu.Lock()
            results = append(results, e.createToolResult(toolUse, result))
            mu.Unlock()
        }(tu)
    }
    
    wg.Wait()
    return results
}
```

### 4.2 API Streaming (Go 版本)

```go
// internal/api/streaming.go

package api

func (c *Client) Stream(ctx context.Context, req *Request) (<-chan StreamEvent, error) {
    events := make(chan StreamEvent, 100)
    
    go func() {
        defer close(events)
        
        // Build HTTP request
        httpReq, err := c.buildHTTPRequest(req)
        if err != nil {
            events <- StreamEvent{Type: "error", Error: err}
            return
        }
        
        // Execute
        resp, err := c.httpClient.Do(httpReq)
        if err != nil {
            events <- StreamEvent{Type: "error", Error: err}
            return
        }
        defer resp.Body.Close()
        
        // SSE parsing
        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            line := scanner.Text()
            
            if !strings.HasPrefix(line, "data: ") {
                continue
            }
            
            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" {
                break
            }
            
            // Parse JSON
            var event map[string]interface{}
            if err := json.Unmarshal([]byte(data), &event); err != nil {
                continue
            }
            
            // Convert to StreamEvent
            streamEvent := parseEvent(event)
            events <- streamEvent
        }
    }()
    
    return events, nil
}

func parseEvent(raw map[string]interface{}) StreamEvent {
    eventType, _ := raw["type"].(string)
    
    switch eventType {
    case "message_start":
        msg := parseMessage(raw["message"].(map[string]interface{}))
        return StreamEvent{Type: "message_start", Message: msg}
        
    case "content_block_start":
        block := parseContentBlock(raw["content_block"].(map[string]interface{}))
        return StreamEvent{Type: "content_block_start", Content: block}
        
    case "content_block_delta":
        delta := parseDelta(raw["delta"].(map[string]interface{}))
        return StreamEvent{Type: "content_block_delta", Delta: delta}
        
    case "content_block_stop":
        return StreamEvent{Type: "content_block_stop"}
        
    case "message_delta":
        delta := raw["delta"].(map[string]interface{})
        usage := raw["usage"].(map[string]interface{})
        return StreamEvent{
            Type: "message_delta",
            StopReason: delta["stop_reason"].(string),
            Usage: parseUsage(usage),
        }
        
    case "message_stop":
        return StreamEvent{Type: "message_stop"}
        
    case "ping":
        return StreamEvent{Type: "ping"}
        
    case "error":
        err := raw["error"].(map[string]interface{})
        return StreamEvent{Type: "error", Error: fmt.Errorf("%v", err)}
        
    default:
        return StreamEvent{Type: "unknown"}
    }
}
```

### 4.3 Bubble Tea TUI

```go
// internal/tui/app.go

package tui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/spinner"
)

type Model struct {
    // State
    messages   []*types.Message
    input      textinput.Model
    viewport   viewport.Model
    spinner    spinner.Model
    ready      bool
    loading    bool
    
    // Query engine
    queryEngine *query.Engine
    eventChan   <-chan StreamEvent
    resultChan  <-chan QueryResult
    
    // Current streaming message
    currentMessage *types.Message
    
    // Permission dialog
    permissionDialog *PermissionDialog
    showPermission   bool
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        textinput.Blink,
        spinner.Tick,
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            return m, tea.Quit
        case tea.KeyEnter:
            if m.showPermission {
                // Handle permission dialog
                m, cmd := m.handlePermissionKey(msg)
                return m, cmd
            }
            // Submit user input
            return m.submitInput()
        }
        
    case tea.WindowSizeMsg:
        m.viewport.Width = msg.Width
        m.viewport.Height = msg.Height - 4 // Reserve space for input
        if !m.ready {
            m.ready = true
        }
        
    case StreamEvent:
        // Handle streaming event from query engine
        m, cmd := m.handleStreamEvent(msg)
        cmds = append(cmds, cmd)
        
    case QueryResult:
        // Query completed
        m.loading = false
        return m, tea.Batch(cmds...)
        
    case PermissionRequestMsg:
        // Show permission dialog
        m.permissionDialog = NewPermissionDialog(msg.ToolUse, msg.Tool)
        m.showPermission = true
        return m, nil
        
    case PermissionResponseMsg:
        // User responded to permission dialog
        m.showPermission = false
        // Send response back to query engine
        cmds = append(cmds, func() tea.Msg { return msg })
    }
    
    // Update components
    m.input, cmd := m.input.Update(msg)
    cmds = append(cmds, cmd)
    
    m.viewport, cmd = m.viewport.Update(msg)
    cmds = append(cmds, cmd)
    
    m.spinner, cmd = m.spinner.Update(msg)
    cmds = append(cmds, cmd)
    
    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    if m.showPermission {
        return m.permissionDialog.View()
    }
    
    return fmt.Sprintf(
        "%s\n\n%s\n\n%s",
        m.viewport.View(),
        m.spinner.View(),
        m.input.View(),
    )
}

// Messages
type StreamEvent stream.StreamEvent
type QueryResult stream.QueryResult

type PermissionRequestMsg struct {
    ToolUse tools.ToolUse
    Tool    tools.Tool
}

type PermissionResponseMsg struct {
    Approved bool
    Always   bool
}
```

---

## 5. 並發模型

### 5.1 Goroutines & Channels

```
Main Goroutine (TUI)
    │
    ├── Query Goroutine (per query)
    │       │
    │       ├── API Streaming Goroutine
    │       │
    │       └── Tool Execution Goroutines (parallel)
    │               │
    │               ├── Tool 1 Goroutine
    │               ├── Tool 2 Goroutine
    │               └── Tool N Goroutine
    │
    └── Session Writer Goroutine (background)
```

### 5.2 Channel 通訊

```go
// Main event loop
eventChan := make(chan StreamEvent, 100)
resultChan := make(chan QueryResult, 1)
permissionChan := make(chan PermissionResponse, 1)

// Query engine publishes events
go queryEngine.Run(ctx, params, eventChan, resultChan, permissionChan)

// TUI subscribes to events
for {
    select {
    case event := <-eventChan:
        // Update UI
    case result := <-resultChan:
        // Query completed
        return
    case <-ctx.Done():
        // Interrupted
        return
    }
}
```

---

## 6. 錯誤處理策略

### 6.1 Error Types

```go
// internal/errors/errors.go

type APIError struct {
    Type    string //"invalid_request", "authentication_error", "rate_limit", "overloaded"
    Message string
    Code    int
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error (%s): %s", e.Type, e.Message)
}

type ToolError struct {
    ToolName string
    Message  string
    Cause    error
}

func (e *ToolError) Error() string {
    return fmt.Sprintf("Tool %s error: %s", e.ToolName, e.Message)
}

type PermissionError struct {
    ToolName string
    Reason   string
}

func (e *PermissionError) Error() string {
    return fmt.Sprintf("Permission denied for %s: %s", e.ToolName, e.Reason)
}
```

### 6.2 Error Recovery

```go
// Retry logic for API errors
func (c *Client) doWithRetry(ctx context.Context, req *Request) (*Response, error) {
    maxRetries := 3
    baseDelay := time.Second
    
    for i := 0; i < maxRetries; i++ {
        resp, err := c.do(ctx, req)
        
        if err == nil {
            return resp, nil
        }
        
        // Check if retryable
        if apiErr, ok := err.(*APIError); ok {
            if apiErr.Type == "rate_limit" || apiErr.Type == "overloaded" {
                delay := baseDelay * time.Duration(1<<i) // Exponential backoff
                time.Sleep(delay)
                continue
            }
        }
        
        return nil, err
    }
    
    return nil, fmt.Errorf("max retries exceeded")
}
```

---

## 7. 配置管理

### 7.1 Settings Schema

```go
// internal/config/settings.go

type Settings struct {
    APIKey         string            `json:"api_key,omitempty"`
    Model          string            `json:"model,omitempty"`
    MaxTokens      int               `json:"max_tokens,omitempty"`
    PermissionMode string            `json:"permission_mode,omitempty"`
    Permissions    PermissionRules   `json:"permissions,omitempty"`
    Tools          map[string]ToolSettings `json:"tools,omitempty"`
}

type PermissionRules struct {
    Allow []PermissionRule `json:"allow,omitempty"`
    Deny  []PermissionRule `json:"deny,omitempty"`
    Ask   []PermissionRule `json:"ask,omitempty"`
}

type PermissionRule struct {
    ToolName     string `json:"tool_name"`
    InputPattern string `json:"input_pattern,omitempty"`
}

type ToolSettings struct {
    Enabled bool                   `json:"enabled"`
    Config  map[string]interface{} `json:"config,omitempty"`
}
```

### 7.2 Settings Loading

```go
// Priority: CLI flags > Project settings > Global settings > Defaults

func LoadSettings() (*Settings, error) {
    settings := DefaultSettings()
    
    // 1. Global settings
    global, err := loadGlobalSettings()
    if err == nil {
        settings = mergeSettings(settings, global)
    }
    
    // 2. Project settings
    project, err := loadProjectSettings()
    if err == nil {
        settings = mergeSettings(settings, project)
    }
    
    // 3. CLI flags (handled separately)
    
    // 4. Environment variables
    if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
        settings.APIKey = apiKey
    }
    
    return settings, nil
}
```

---

## 8. 測試策略

### 8.1 Unit Tests

```go
// internal/tools/bash/bash_test.go

func TestBashTool_Execute(t *testing.T) {
    tool := NewBashTool()
    
    tests := []struct {
        name    string
        input   map[string]interface{}
        want    string
        wantErr bool
    }{
        {
            name:  "echo command",
            input: map[string]interface{}{"command": "echo hello"},
            want:  "hello\n",
        },
        {
            name:    "invalid command",
            input:   map[string]interface{}{"command": ""},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := tool.Execute(context.Background(), tt.input, &ToolContext{})
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Contains(t, result.Content, tt.want)
            }
        })
    }
}
```

### 8.2 Integration Tests

```go
// internal/query/engine_test.go

func TestEngine_Query(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    client := api.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
    engine := NewEngine(client, ...)
    
    params := QueryParams{
        Messages: []*types.Message{
            {Type: types.MessageTypeUser, Content: []types.ContentBlock{
                {Type: "text", Text: "Say hello"},
            }},
        },
        Model: "claude-sonnet-4-20250514",
    }
    
    events, results := engine.Query(context.Background(), params)
    
    var gotResponse bool
    for event := range events {
        if event.Type == "message_stop" {
            gotResponse = true
        }
    }
    
    assert.True(t, gotResponse)
    
    result := <-results
    assert.Equal(t, "completed", result.Terminal.Reason)
}
```

---

## 9. 效能考量

### 9.1 記憶體管理

- Message history 使用指標切片 (`[]*Message`)
- 大型 tool result 寫入臨時檔案
- Context compaction 保持上下文在合理大小

### 9.2 並發優化

- Tool execution 使用 goroutine pool
- Channel buffer 大小根據實際測試調整
- 避免在 hot path 進行記憶體分配

### 9.3 快取策略

- CLAUDE.md 內容快取 (帶 TTL)
- Permission rules 快取
- API response 不快取 (即時性)

---

## 10. 部署 & Build

### 10.1 Build Commands

```makefile
# Makefile

.PHONY: build test clean

VERSION := $(shell git describe --tags --always)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/claude-code ./cmd/claude-code

test:
	go test -v ./...

test-integration:
	go test -v -tags=integration ./...

clean:
	rm -rf bin/
	go clean

install: build
	cp bin/claude-code /usr/local/bin/
```

### 10.2 Cross-Compilation

```bash
# Build for multiple platforms
GOOS=darwin GOARCH=arm64 go build -o bin/claude-code-darwin-arm64 ./cmd/claude-code
GOOS=darwin GOARCH=amd64 go build -o bin/claude-code-darwin-amd64 ./cmd/claude-code
GOOS=linux GOARCH=amd64 go build -o bin/claude-code-linux-amd64 ./cmd/claude-code
GOOS=windows GOARCH=amd64 go build -o bin/claude-code-windows-amd64.exe ./cmd/claude-code
```

---

## 總結

**核心設計**:
- Bubble Tea 作為 TUI framework
- Channel-based streaming (不使用 generator)
- Goroutine pool 進行並行 tool 執行
- Interface-based tool system
- JSONL 檔案作為 session storage

**關鍵差異 (vs TS 版本)**:
- 無 React，改用 Bubble Tea
- 無 Zod，改用 struct + manual validation
- 無 generator，改用 channel
- 無 module singletons，改用 package variables
- 顯式的錯誤處理

**預估程式碼量**: ~3000-4000 LOC (不包含測試)

---

*Document generated for Go porting reference*