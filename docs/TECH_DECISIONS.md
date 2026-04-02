# Technical Decisions (ADR) - CC-CLI-Go

> **Purpose**: Record key technical decisions for Go port

---

## 名詞解釋 / Glossary

| 名詞 | 說明 | 本專案實作方式 |
|------|------|----------------|
| ADR (Architecture Decision Record) | 架構決策紀錄，用來記錄「為什麼這樣設計」 | 每個 ADR 含 Context / Decision / Rationale / Consequences |
| KISS | Keep It Simple, Stupid；優先簡單可維護方案 | 優先內建能力與少依賴，不先導入高複雜度架構 |
| Bubble Tea | Go 的 TUI 框架（Elm-style） | 以 `Model -> Update -> View` 管理互動式 CLI |
| Elm-style Architecture | 狀態與事件驅動 UI 架構 | 用 `tea.Msg` 驅動狀態變更，避免隱式副作用 |
| Streaming | API 回應分段即時輸出 | 使用 SSE 讀取事件並轉為內部 stream event |
| Channel-based Streaming | 用 Go channel 傳遞串流資料 | `Query()` 回傳 event channel + result channel |
| Goroutine | Go 輕量執行緒 | 工具並行執行、串流處理均以 goroutine 實作 |
| WaitGroup | 等待多個 goroutine 完成的同步機制 | 收斂多工具執行結果，確保流程在完成後再繼續 |
| Mutex | 共享資料互斥鎖 | 並行工具寫入共享結果切片時保護資料一致性 |
| Concurrency-safe Tool | 可安全並行執行的工具 | 標記可並行工具，否則退回串行執行 |
| JSONL (JSON Lines) | 每行一筆 JSON 的儲存格式 | session transcript 採 append-only JSONL 寫入 |
| Append-only | 只追加、不覆寫既有資料 | 對話記錄逐行追加，降低寫入複雜度 |
| Schema Validation | 輸入資料格式與欄位驗證 | `struct tags` + `Validate()` 手動檢查必填與約束 |
| Provider | LLM API 提供者（Anthropic/Bedrock/Vertex） | 目前只支援 Anthropic，避免多供應商複雜度 |
| MCP | Model Context Protocol，外部工具/資源整合協議 | 目前不實作，保留未來擴展空間 |
| System Prompt | 注入模型的系統層指令 | 組合環境資訊與規則後放入請求 |
| Context Compaction | 壓縮對話上下文，降低 token 成本 | 超過閾值時摘要舊訊息，保留最近互動 |
| Permission Mode | 工具執行權限策略模式 | `default/accept/plan/auto` 決定是否詢問使用者 |
| errors.Is / errors.As | Go 錯誤鏈檢查與型別斷言 | 用於判斷取消、API 錯誤等可恢復情境 |

---

## ADR-001: Language Selection

### Status

Accepted

### Context

Need to choose a programming language for reimplementing CC-CLI-Go as a learning project.

Options considered:
- TypeScript (original)
- Go
- Rust
- Python

### Decision

**Go**

### Rationale

| Criteria | Go | Rust | TypeScript | Python |
|----------|-----|------|------------|--------|
| Learning curve | ★★★★★ | ★★☆☆☆ | ★★★☆☆ | ★★★★★ |
| Build simplicity | ★★★★★ | ★★★☆☆ | ★★★☆☆ | ★★★☆☆ |
| Performance | ★★★★☆ | ★★★★★ | ★★★☆☆ | ★★☆☆☆ |
| CLI ecosystem | ★★★★★ | ★★★★☆ | ★★★★★ | ★★★☆☆ |
| TUI ecosystem | ★★★★☆ | ★★★☆☆ | ★★★★★ | ★★☆☆☆ |
| Single binary | ★★★★★ | ★★★★★ | ★★☆☆☆ | ★★☆☆☆ |
| Learning value | ★★★★★ | ★★★★★ | ★★★★☆ | ★★★☆☆ |

**Key factors**:
1. **Learning curve**: Go is easiest to learn for this project
2. **Build simplicity**: Single binary, no runtime dependencies
3. **TUI ecosystem**: Bubble Tea is mature and well-documented
4. **User preference**: Explicitly chose Go for learning purposes

### Consequences

- **Positive**: Fast development, easy deployment, good learning experience
- **Negative**: Not as performant as Rust, less rich TUI than TS/Ink
- **Neutral**: Different patterns from original (no React, no generators)

---

## ADR-002: TUI Framework

### Status

Accepted

### Context

Need a TUI framework for the interactive REPL interface.

Options considered:
- Bubble Tea (Elm-style)
- tcell + tview (widget-based)
- termdash (dashboard-oriented)
- Custom ncurses wrapper

### Decision

**Bubble Tea**

### Rationale

| Framework | Ecosystem | Architecture | Learning Curve | Community |
|-----------|-----------|--------------|----------------|-----------|
| Bubble Tea | ★★★★★ | Elm-style (functional) | ★★★★☆ | ★★★★★ |
| tcell + tview | ★★★★☆ | Widget-based | ★★★☆☆ | ★★★☆☆ |
| termdash | ★★★☆☆ | Dashboard | ★★★☆☆ | ★★☆☆☆ |
| Custom | ★☆☆☆☆ | N/A | ★☆☆☆☆ | ★☆☆☆☆ |

**Key factors**:
1. **Most popular**: Bubble Tea is the de-facto standard for Go TUIs
2. **Elm architecture**: Clean separation of model/update/view
3. **Rich examples**: Many real-world projects to learn from
4. **Active maintenance**: Charm.sh actively maintains

### Consequences

- **Positive**: Clean architecture, good examples, active community
- **Negative**: Learning curve for Elm-style thinking
- **Neutral**: Different paradigm from React (original)

### Examples

```go
// Elm-style architecture
type Model struct { ... }

func (m Model) Init() tea.Cmd { ... }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
func (m Model) View() string { ... }
```

---

## ADR-003: API Provider Support

### Status

Accepted

### Context

Original codebase supports multiple providers:
- Anthropic (direct)
- AWS Bedrock
- Google Vertex
- Azure (future)

### Decision

**Anthropic API only**

### Rationale

**Learning focus**: This is a learning project, not production software.

| Provider | Complexity | Learning Value | Maintenance |
|----------|------------|----------------|-------------|
| Anthropic only | Low | High | Low |
| Multi-provider | High | Medium | High |

**Simplification benefits**:
1. **Less code**: ~30% reduction in API layer
2. **Easier testing**: Single API to mock
3. **Clearer errors**: Direct Anthropic error messages
4. **Faster development**: Focus on core features

### Consequences

- **Positive**: Simpler codebase, faster development
- **Negative**: Users must use Anthropic API directly
- **Neutral**: Can add multi-provider support later if needed

### Future Consideration

If multi-provider support becomes necessary:
1. Define Provider interface
2. Implement Anthropic as first provider
3. Add Bedrock/Vertex as additional implementations

---

## ADR-004: MCP Support

### Status

Accepted

### Context

MCP (Model Context Protocol) allows external tools/servers to be integrated.

Original codebase has:
- MCP server connection management
- MCP tool discovery
- MCP resource handling
- OAuth for MCP servers

### Decision

**No MCP support**

### Rationale

**Complexity analysis**:

| Feature | Complexity | Value for Learning |
|---------|------------|-------------------|
| MCP server connections | High | Low |
| Tool discovery | Medium | Low |
| Resource handling | Medium | Low |
| OAuth flow | High | Low |
| **Total** | **Very High** | **Low** |

**Reasons to exclude**:
1. **Learning focus**: MCP adds complexity without teaching core concepts
2. **Core features**: All essential tools are built-in
3. **Maintenance**: MCP spec changes require updates
4. **Testing**: MCP servers are external dependencies

### Consequences

- **Positive**: 40% less code, simpler architecture, faster development
- **Negative**: Cannot extend with external tools
- **Neutral**: Most users won't miss MCP

### Alternative

If tool extensibility is needed later:
1. Plugin system with simpler interface
2. Direct tool registration
3. Shell-based integration

---

## ADR-005: Streaming Architecture

### Status

Accepted

### Context

Original uses TypeScript AsyncGenerator for streaming:

```typescript
async function* query(params): AsyncGenerator<StreamEvent, Terminal>
```

### Decision

**Channel-based streaming**

### Rationale

**Go idioms**:
- No generators in Go
- Channels are the idiomatic way to stream data
- Goroutines for concurrent production

**Pattern**:

```go
func (e *Engine) Query(ctx context.Context, params QueryParams) (<-chan StreamEvent, <-chan QueryResult) {
    events := make(chan StreamEvent, 100)
    results := make(chan QueryResult, 1)
    
    go func() {
        defer close(events)
        defer close(results)
        // Stream events to channel
    }()
    
    return events, results
}
```

**Benefits**:
1. **Idiomatic Go**: Channels are built for this
2. **Composable**: Easy to chain/transform streams
3. **Cancellable**: Context integration is natural
4. **Buffered**: Control backpressure with buffer size

### Consequences

- **Positive**: Idiomatic, composable, cancellable
- **Negative**: Different mental model from generators
- **Neutral**: Buffer size tuning needed

---

## ADR-006: Schema Validation

### Status

Accepted

### Context

Original uses Zod for runtime schema validation:

```typescript
const InputSchema = z.object({
    command: z.string(),
    timeout: z.number().optional(),
})
```

### Decision

**Go struct tags + manual validation**

### Rationale

**Options considered**:

| Approach | Type Safety | Runtime Validation | Learning Curve |
|----------|-------------|-------------------|----------------|
| struct tags + manual | ★★★★★ | ★★★☆☆ | ★★★★★ |
| go-playground/validator | ★★★★☆ | ★★★★☆ | ★★★★☆ |
| Custom DSL | ★★☆☆☆ | ★★★★★ | ★★☆☆☆ |

**Decision**: Keep it simple with struct tags

```go
type BashInput struct {
    Command string `json:"command" validate:"required"`
    Timeout int    `json:"timeout,omitempty"`
}

func (i *BashInput) Validate() error {
    if i.Command == "" {
        return errors.New("command is required")
    }
    return nil
}
```

**Reasons**:
1. **Learning**: Manual validation teaches the domain
2. **Transparency**: No magic, clear what's happening
3. **Simplicity**: One less dependency
4. **Go idiomatic**: struct tags are standard

### Consequences

- **Positive**: Simple, explicit, no dependencies
- **Negative**: More boilerplate than Zod
- **Neutral**: Can add validator library later if needed

---

## ADR-007: Session Storage Format

### Status

Accepted

### Context

Original stores sessions as JSONL (JSON Lines):

```
{"uuid":"...","type":"user","content":[...]}
{"uuid":"...","type":"assistant","content":[...]}
```

### Decision

**JSONL (same as original)**

### Rationale

**Benefits of JSONL**:
1. **Append-only**: Efficient for streaming writes
2. **Line-delimited**: Easy to parse incrementally
3. **Human-readable**: Can inspect with any text editor
4. **Interoperability**: Compatible with original if needed

**Alternatives considered**:

| Format | Append | Human-Readable | Parse Speed | Size |
|--------|--------|----------------|-------------|------|
| JSONL | ★★★★★ | ★★★★★ | ★★★★☆ | ★★★☆☆ |
| SQLite | ★★★★★ | ★★☆☆☆ | ★★★★★ | ★★★★★ |
| Binary | ★★★☆☆ | ★☆☆☆☆ | ★★★★★ | ★★★★★ |
| JSON array | ★★☆☆☆ | ★★★★★ | ★★★☆☆ | ★★★☆☆ |

### Consequences

- **Positive**: Same format as original, simple implementation
- **Negative**: Larger than binary formats
- **Neutral**: Can migrate to SQLite later if performance issues

---

## ADR-008: Tool Execution Concurrency

### Status

Accepted

### Context

Original can execute multiple tools in parallel when they're concurrency-safe.

### Decision

**Goroutine pool with WaitGroup**

### Rationale

**Pattern**:

```go
func (e *Engine) executeTools(ctx context.Context, tools []ToolUse) []*ToolResult {
    var wg sync.WaitGroup
    var mu sync.Mutex
    results := make([]*ToolResult, len(tools))
    
    for i, t := range tools {
        if !t.IsConcurrencySafe() {
            // Execute serially
            results[i] = e.executeTool(ctx, t)
            continue
        }
        
        wg.Add(1)
        go func(idx int, tool ToolUse) {
            defer wg.Done()
            result := e.executeTool(ctx, tool)
            mu.Lock()
            results[idx] = result
            mu.Unlock()
        }(i, t)
    }
    
    wg.Wait()
    return results
}
```

**Reasons**:
1. **Idiomatic Go**: Goroutines + WaitGroup is standard
2. **Simple**: No external dependencies
3. **Flexible**: Easy to add limits/rate control later
4. **Safe**: Mutex protects shared state

### Consequences

- **Positive**: Simple, idiomatic, no dependencies
- **Negative**: No sophisticated scheduling/rate limiting
- **Neutral**: Can evolve to worker pool if needed

---

## ADR-009: Configuration Location

### Status

Accepted

### Context

Original stores config at:
- Global: `~/.claude/settings.json`
- Project: `.claude/settings.json`

### Decision

**Same locations, same format**

### Rationale

**Benefits**:
1. **Compatibility**: Can share settings with original
2. **Discoverability**: Standard XDG-like location
3. **Simplicity**: Same code, different audience

**File structure**:

```
~/.claude/
├── settings.json          # Global settings
├── credentials.json       # API keys (optional, env var preferred)
├── sessions/              # Session history
│   ├── <uuid>.jsonl
│   └── <uuid>.metadata.json
└── projects/
    └── <cwd-hash>/
        ├── settings.json  # Project-specific settings
        └── sessions/

Project root:
.claude/
└── settings.json          # Project settings (checked into git)
```

### Consequences

- **Positive**: Compatible with original, familiar to users
- **Negative**: Same limitations as original (no XDG support)
- **Neutral**: Implementation is straightforward

---

## ADR-010: Error Handling

### Status

Accepted

### Context

Go uses explicit error handling, different from TypeScript's try/catch.

### Decision

**Explicit errors with custom types**

### Rationale

**Pattern**:

```go
// Define error types
type APIError struct {
    Type    string
    Message string
    Code    int
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error (%s): %s", e.Type, e.Message)
}

// Use errors.Is/As for checking
if errors.Is(err, context.Canceled) {
    // Handle cancellation
}

var apiErr *APIError
if errors.As(err, &apiErr) {
    // Handle API error
}
```

**Reasons**:
1. **Idiomatic Go**: Explicit error handling is standard
2. **Type-safe**: errors.As for type checking
3. **Composable**: errors.Is for wrapping
4. **Clear**: No hidden control flow

### Consequences

- **Positive**: Explicit, type-safe, idiomatic
- **Negative**: More verbose than try/catch
- **Neutral**: Requires discipline to check all errors

---

## Decision Summary

| ADR | Decision | Primary Reason |
|-----|----------|----------------|
| 001 | Go | Learning curve + simplicity |
| 002 | Bubble Tea | Ecosystem + Elm architecture |
| 003 | Anthropic only | Simplicity for learning |
| 004 | No MCP | Complexity reduction |
| 005 | Channel streaming | Idiomatic Go |
| 006 | Struct tags + manual | Simplicity + learning |
| 007 | JSONL | Compatibility + simplicity |
| 008 | Goroutine + WaitGroup | Idiomatic Go |
| 009 | Same config location | Compatibility |
| 010 | Explicit errors | Idiomatic Go |

---

*Document generated for Go porting reference*
