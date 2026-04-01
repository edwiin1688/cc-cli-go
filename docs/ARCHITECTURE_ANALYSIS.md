# Architecture Analysis - Claude Code (TS/Bun)

> **目的**: 提供給 Go 版本移植的完整架構參考

---

## 1. 專案概述

**專案類型**: 反編譯/逆向工程版本的 Anthropic Claude Code CLI

**Runtime**: Bun (非 Node.js)

**Build**: 單一 bundle (~25MB), ESM + TSX

**Feature Flags**: 全部關閉 (`feature()` 永遠回傳 `false`)

---

## 2. Entry Points & Bootstrap

### 2.1 Entry Flow

```
cli.tsx (入口)
    ↓ polyfill: feature(), MACRO, BUILD_TARGET
    ↓ fast-path: --version, --dump-system-prompt
    ↓
main.tsx (Commander CLI)
    ↓ init(): telemetry, config, trust dialog
    ↓
REPL.tsx (Interactive TUI)
    ↓
QueryEngine.ts → query.ts
```

### 2.2 Key Entry Files

| File | Purpose | Lines |
|------|---------|-------|
| `src/entrypoints/cli.tsx` | True entrypoint, runtime polyfills | 319 |
| `src/main.tsx` | Commander CLI definition, initialization | ~3000+ |
| `src/entrypoints/init.ts` | One-time initialization (telemetry, config) | - |

### 2.3 Bootstrap Polyfills (cli.tsx)

```typescript
// Feature flag - always false in decompiled version
const feature = (_name: string) => false;

// Build-time macros
globalThis.MACRO = {
    VERSION: "2.1.888",
    BUILD_TIME: new Date().toISOString(),
    ...
};

globalThis.BUILD_TARGET = "external";
globalThis.BUILD_ENV = "production";
globalThis.INTERFACE_TYPE = "stdio";
```

**移植影響**: Go 版本不需要 feature flag polyfill，直接移除相關分支。

---

## 3. Core Query Loop

### 3.1 Query Flow Diagram

```
query.ts (核心 API 對話循環)
    │
    ├─→ 建構 messages + systemPrompt
    │
    ├─→ autocompact / microcompact (context 管理)
    │
    ├─→ callModel() → Anthropic API Streaming
    │       │
    │       ├─→ message_start
    │       ├─→ content_block_start (thinking, tool_use)
    │       ├─→ content_block_delta (text, tool_input)
    │       ├─→ content_block_stop
    │       ├─→ message_delta (stop_reason)
    │       ├─→ message_stop
    │
    ├─→ tool_use 處理
    │       │
    │       ├─→ StreamingToolExecutor (並行執行)
    │       │       ├─→ canUseTool() (權限檢查)
    │       │       ├─→ tool.call()
    │       │       ├─→ yield tool_result
    │       │
    │       ├─→ runTools() (串行執行)
    │
    ├─→ stop hooks 檢查
    │
    └─→ continue 或 return Terminal
```

### 3.2 Query.ts 關鍵結構

```typescript
export type QueryParams = {
  messages: Message[]
  systemPrompt: SystemPrompt
  userContext: { [k: string]: string }
  systemContext: { [k: string]: string }
  canUseTool: CanUseToolFn
  toolUseContext: ToolUseContext
  fallbackModel?: string
  querySource: QuerySource
  maxOutputTokensOverride?: number
  maxTurns?: number
  skipCacheWrite?: boolean
  taskBudget?: { total: number }
  deps?: QueryDeps
}

export async function* query(params: QueryParams): AsyncGenerator<
  | StreamEvent
  | RequestStartEvent
  | Message
  | TombstoneMessage
  | ToolUseSummaryMessage,
  Terminal
>
```

**移植重點**:
- Go 使用 goroutine + channel 實現 streaming
- Generator 改為 channel-based streaming
- State 改用 struct + mutable fields

---

## 4. Tool System

### 4.1 Tool Interface

```typescript
export type Tool<
  Input extends AnyObject = AnyObject,
  Output = unknown,
  P extends ToolProgressData = ToolProgressData,
> = {
  name: string
  aliases?: string[]
  searchHint?: string
  
  // Core methods
  call(
    args: z.infer<Input>,
    context: ToolUseContext,
    canUseTool: CanUseToolFn,
    parentMessage: AssistantMessage,
    onProgress?: ToolCallProgress<P>,
  ): Promise<ToolResult<Output>>
  
  description(input, options): Promise<string>
  prompt(options): Promise<string>
  
  // Validation & Permissions
  validateInput?(input, context): Promise<ValidationResult>
  checkPermissions(input, context): Promise<PermissionResult>
  
  // Schema
  inputSchema: AnyObject // Zod schema
  inputJSONSchema?: ToolInputJSONSchema
  outputSchema?: z.ZodType<unknown>
  
  // Behavior flags
  isEnabled(): boolean
  isConcurrencySafe(input): boolean
  isReadOnly(input): boolean
  isDestructive?(input): boolean
  interruptBehavior?(): 'cancel' | 'block'
  
  // UI rendering (React components)
  renderToolUseMessage(input, options): React.ReactNode
  renderToolResultMessage?(content, progress, options): React.ReactNode
  renderToolUseProgressMessage?(progress, options): React.ReactNode
  
  // API serialization
  mapToolResultToToolResultBlockParam(content, toolUseID): ToolResultBlockParam
  
  // Other
  maxResultSizeChars: number
  userFacingName(input): string
  getActivityDescription?(input): string | null
  toAutoClassifierInput(input): unknown
  ...
}
```

### 4.2 Tool Registry (src/tools.ts)

工具列表（需移植的核心工具）:

| Tool | Purpose | 必要性 |
|------|---------|--------|
| BashTool | Execute shell commands | **必要** |
| FileReadTool | Read file contents | **必要** |
| FileEditTool | Edit files (exact string replace) | **必要** |
| GlobTool | File pattern matching | **必要** |
| GrepTool | Search file contents | **必要** |
| FileWriteTool | Write new files | **必要** |
| AgentTool | Subagent spawning | 可選 |
| TodoWriteTool | Task tracking | 可選 |
| WebFetchTool | Fetch URLs | 可選 |
| WebSearchTool | Web search | 可選 |

**移除的工具**:
- Computer Use tools (stub)
- MCP tools
- Voice tools
- Magic Docs
- Plugin tools

### 4.3 Tool Execution Flow

```
StreamingToolExecutor
    │
    ├─→ addTool(toolBlock, assistantMessage)
    │       └─→ enqueue to pendingTools
    │
    ├─→ getCompletedResults()
    │       │
    │       ├─→ pendingTools queue
    │       │
    │       ├─→ executeTool()
    │       │       ├─→ findToolByName()
    │       │       ├─→ parse input (Zod)
    │       │       ├─→ validateInput()
    │       │       ├─→ canUseTool() (權限)
    │       │       │       ├─→ Permission dialog (REPL)
    │       │       │       └─→ Auto-approve/deny
    │       │       ├─→ tool.call()
    │       │       └─→ mapToolResultToToolResultBlockParam()
    │       │
    │       └─→ yield UserMessage with tool_result
    │
    └─→ getRemainingResults() (abort cleanup)
```

**移植重點**:
- Go 用 goroutine pool 實現並行工具執行
- Zod schema 改用 Go struct + validator
- React UI 改用 Bubble Tea components

---

## 5. API Layer

### 5.1 API Client (src/services/api/claude.ts)

```typescript
// 支援的 Providers
- Anthropic direct (sdk)
- AWS Bedrock (@anthropic-ai/bedrock-sdk)
- Google Vertex (@anthropic-ai/vertex-sdk)
- Azure (future)

// API Request 建構
buildRequestParams({
  systemPrompt,
  messages,
  tools,
  model,
  betas,
  thinkingConfig,
  max_tokens,
  ...
})

// Streaming Events
type BetaRawMessageStreamEvent =
  | { type: 'message_start', message: Message }
  | { type: 'content_block_start', index, content_block }
  | { type: 'content_block_delta', index, delta }
  | { type: 'content_block_stop', index }
  | { type: 'message_delta', delta, usage }
  | { type: 'message_stop' }
  | { type: 'ping' }
  | { type: 'error', error }
```

**移植重點**:
- Go 用 `net/http` + SSE streaming
- 只保留 Anthropic direct provider
- Event parsing 改用 struct unmarshal

### 5.2 Key API Types

```typescript
// Message types
type Message =
  | UserMessage
  | AssistantMessage
  | SystemMessage
  | ToolUseSummaryMessage
  | TombstoneMessage
  | AttachmentMessage
  | ProgressMessage

// Content blocks
type ContentBlock =
  | { type: 'text', text: string }
  | { type: 'thinking', thinking: string }
  | { type: 'redacted_thinking' }
  | { type: 'tool_use', id, name, input }
  | { type: 'tool_result', tool_use_id, content, is_error }
```

---

## 6. State Management

### 6.1 AppState (src/state/AppState.tsx)

```typescript
export type AppState = {
  messages: Message[]
  tools: Tools
  toolPermissionContext: ToolPermissionContext
  mcp: {
    clients: MCPServerConnection[]
    tools: Tools
    resources: Record<string, ServerResource[]>
  }
  fastMode: boolean
  effortValue: number | undefined
  advisorModel: string | undefined
  sessionInfo: SessionInfo
  ... // 30+ fields
}
```

### 6.2 State Store (Zustand-style)

```typescript
// src/state/store.ts
export function createStore<T>(initialState: T): {
  getState(): T
  setState(updater: (prev: T) => T): void
  subscribe(listener: (state: T) => void): () => void
}
```

### 6.3 Bootstrap State (Module-level Singletons)

```typescript
// src/bootstrap/state.ts
let sessionId: string
let cwd: string
let mainLoopModelOverride: string | undefined
let tokenCounts: { input: number, output: number }
...
```

**移植重點**:
- Go 用 struct + methods 實現 state
- Module-level singletons 改用 package-level variables
- Zustand-style store 改用 simple struct + mutex

---

## 7. Context & System Prompt

### 7.1 Context Building (src/context.ts)

```typescript
export async function getSystemContext(): Promise<{ [k: string]: string }> {
  return {
    git_status: await getGitStatus(),
    git_branch: await getBranch(),
    date: new Date().toISOString(),
    cwd: getCwd(),
    ... // 其他環境資訊
  }
}

export async function getUserContext(): Promise<{ [k: string]: string }> {
  // 從 CLAUDE.md, memory files 取得
}
```

### 7.2 CLAUDE.md Discovery (src/utils/claudemd.ts)

```typescript
// 從當前目錄往上搜尋 CLAUDE.md
// 路徑: ./CLAUDE.md, ../CLAUDE.md, ../../CLAUDE.md, ...
// 合併所有找到的 CLAUDE.md 內容
```

**移植重點**: 直接移植，邏輯簡單。

---

## 8. Permission System

### 8.1 Permission Modes

```typescript
type PermissionMode =
  | 'default'  // Ask for permission
  | 'accept'   // Auto-accept all
  | 'plan'     // Plan mode (read-only)
  | 'auto'     // AI decides (with gates)
```

### 8.2 Permission Rules

```typescript
type PermissionRule = {
  toolName: string
  inputPattern?: string  // e.g. "git *", "*.ts"
  behavior: 'allow' | 'deny' | 'ask'
}

// Sources:
// - ~/.claude/settings.json
// - Project .claude/settings.json
// - Policy (enterprise)
// - CLI flags
```

### 8.3 Permission Check Flow

```
canUseTool(tool, input)
    │
    ├─→ Check permission mode
    │       ├─→ 'accept': return true
    │       ├─→ 'plan': check isReadOnly
    │       ├─→ 'auto': run classifier
    │
    ├─→ Check rules by source
    │       ├─→ alwaysAllowRules
    │       ├─→ alwaysDenyRules
    │       ├─→ alwaysAskRules
    │
    ├─→ No matching rule → Show dialog
    │       ├─→ Accept / Deny / Always Accept / Always Deny
    │       └─→ Update rules
    │
    └─→ Return result
```

**移植重點**: 完整移植，邏輯核心。

---

## 9. TUI Layer (Ink)

### 9.1 Ink Framework (React for Terminal)

```
src/ink/
    ├─→ ink.tsx (render wrapper + ThemeProvider)
    ├─→ components/
    │       ├─→ App.tsx
    │       ├─→ Box.tsx
    │       ├─→ Text.tsx
    │       ├─→ ScrollBox.tsx
    │       ├─→ Button.tsx
    │       ├─→ ...
    ├─→ reconciler (custom React reconciler)
    └─→ hooks/
            ├─→ useInput()
            ├─→ useTerminalSize()
            └─→ useSearchHighlight()
```

### 9.2 REPL Screen (src/screens/REPL.tsx)

```typescript
// Main interactive screen
<REPL>
    ├─→ <PromptInput> (user input)
    ├─→ <Messages> (conversation history)
    ├─→ <PermissionDialog> (tool permission)
    ├─→ <Notifications>
    ├─→ <Spinner>
    └─→ <Keybindings>
```

**移植重點**: 改用 Bubble Tea，完全重寫 UI layer。

---

## 10. Session Management

### 10.1 Session Storage

```
~/.claude/
    ├─→ sessions/
    │       ├─→ <uuid>.jsonl (transcript)
    │       └──→ <uuid>.metadata.json
    ├─→ projects/
    │       └─→ <cwd-hash>/
    │               ├─→ sessions/
    │               ├─→ settings.json
    │               └───→ memory/
    ├─→ settings.json (global settings)
    └──→ credentials.json (OAuth tokens)
```

### 10.2 Conversation Recovery

```typescript
// src/utils/conversationRecovery.ts
loadConversationForResume(sessionId)
    │
    ├─→ Read .jsonl transcript
    ├─→ Parse messages
    ├─→ Rebuild state
    └─→ Return messages + metadata
```

**移植重點**: 直接移植，檔案格式相同。

---

## 11. Key Dependencies

### 11.1 Production Dependencies

| Package | Purpose | Go 替代方案 |
|---------|---------|------------|
| `@anthropic-ai/sdk` | API client | 自建 HTTP client |
| `@commander-js/extra-typings` | CLI framework | `cobra` |
| `react` + `react-reconciler` | TUI (Ink) | `bubbletea` |
| `zod` | Schema validation | `go-playground/validator` |
| `chalk` | Terminal colors | `fatih/color` |
| `execa` | Process execution | `os/exec` |
| `chokidar` | File watching | `fsnotify` |
| `ignore` | gitignore parsing | 自建或 `go-git` |
| `diff` | Diff generation | 自建 |
| `marked` | Markdown parsing | `gomarkdown/markdown` |
| `axios` | HTTP requests | `net/http` |
| `yaml` | YAML parsing | `go-yaml/yaml` |
| `jsonc-parser` | JSONC parsing | 自建 |

### 11.2 Stub/Removed Dependencies

**已移除**:
- `@ant/computer-use-*` (stub packages)
- `audio-capture-napi`, `image-processor-napi`, `modifiers-napi`, `url-handler-napi` (napi stubs)
- `@growthbook/growthbook` (analytics - empty impl)
- `@opentelemetry/*` (analytics - empty impl)
- `sharp` (image - stub)

---

## 12. Feature Flags (全部關閉)

```typescript
// 所有 feature() 調用都回傳 false
// 這些分支在 Go 版本中完全移除:

feature('ABLATION_BASELINE')
feature('DUMP_SYSTEM_PROMPT')
feature('DAEMON')
feature('BG_SESSIONS')
feature('TEMPLATES')
feature('BYOC_ENVIRONMENT_RUNNER')
feature('SELF_HOSTED_RUNNER')
feature('BRIDGE_MODE')
feature('SSH_REMOTE')
feature('KAIROS')
feature('COORDINATOR_MODE')
feature('CHICAGO_MCP')
feature('DIRECT_CONNECT')
feature('LODESTONE')
feature('REACTIVE_COMPACT')
feature('CONTEXT_COLLAPSE')
feature('HISTORY_SNIP')
feature('TOKEN_BUDGET')
feature('EXPERIMENTAL_SKILL_SEARCH')
feature('UPLOAD_USER_SETTINGS')
feature('TRANSCRIPT_CLASSIFIER')
feature('CACHED_MICROCOMPACT')
...
```

---

## 13. TypeScript 錯誤 (~1341 個)

**原因**: 反編譯造成的型別推斷問題

**類型**:
- `unknown` / `never` / `{}` assignments
- Missing type annotations
- Circular dependencies

**影響**: 不影響 Bun runtime 執行，可忽略

---

## 14. Build & Scripts

### 14.1 Build Command

```bash
bun build src/entrypoints/cli.tsx --outdir dist --target bun
```

輸出: `dist/cli.js` (~25MB 單一 bundle)

### 14.2 Dev Command

```bash
bun run src/entrypoints/cli.tsx
```

### 14.3 Package.json Scripts

```json
{
  "build": "bun build src/entrypoints/cli.tsx --outdir dist --target bun",
  "dev": "bun run src/entrypoints/cli.tsx",
  "lint": "biome lint src/",
  "lint:fix": "biome lint --fix src/",
  "format": "biome format --write src/",
  "test": "bun test"
}
```

---

## 15. 移植優先級

### Phase 1: Core Loop (必要)

```
1. CLI entry + argument parsing
2. Anthropic API streaming client
3. query() loop
4. Tool system (Bash, Read, Edit, Glob, Grep, Write)
5. Permission system
6. Context building (CLAUDE.md, git)
7. Session storage
```

### Phase 2: TUI (必要)

```
8. Bubble Tea REPL screen
9. Message rendering
10. Permission dialog
11. Input handling
```

### Phase 3: Polish (可選)

```
12. Compaction (context management)
13. Better error handling
14. Performance optimization
15. Testing
```

---

## 總結

**核心架構**: Entry → Init → REPL → QueryEngine → Query → API → Tools

**關鍵移植點**:
1. Streaming generator → Go channel
2. React/Ink → Bubble Tea
3. Zod → Go struct + validator
4. Module singletons → Package variables
5. Tool interface → Go interface

**移除模組**: Analytics, MCP, Computer Use, Voice, Plugins, Bridge, Daemon, Templates

**預估複雜度**: 中等（核心邏輯清晰，主要工作是重寫 TUI 和 streaming）

---

*Document generated for Go porting reference*