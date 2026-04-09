# Graph Report - .  (2026-04-10)

## Corpus Check
- Corpus is ~37,102 words - fits in a single context window. You may not need a graph.

## Summary
- 521 nodes · 600 edges · 54 communities detected
- Extraction: 98% EXTRACTED · 2% INFERRED · 0% AMBIGUOUS · INFERRED: 13 edges (avg confidence: 0.89)
- Token cost: 15,000 input · 8,000 output

## God Nodes (most connected - your core abstractions)
1. `Input` - 14 edges
2. `BashTool` - 9 edges
3. `EditTool` - 9 edges
4. `GlobTool` - 9 edges
5. `GrepTool` - 9 edges
6. `ReadTool` - 9 edges
7. `WriteTool` - 9 edges
8. `Settings` - 8 edges
9. `Model` - 8 edges
10. `Compactor` - 7 edges

## Surprising Connections (you probably didn't know these)
- `Testing Strategy with 127 Test Cases` --conceptually_related_to--> `Testing Framework with 127 Test Cases`  [INFERRED]
  TESTING.md → CHANGELOG.md
- `Permission System Implementation` --conceptually_related_to--> `Permission System with Multiple Modes`  [INFERRED]
  docs/CORE_FEATURES.md → CHANGELOG.md
- `27 Completed Features Across 3 Phases` --conceptually_related_to--> `Core Features Overview`  [INFERRED]
  TODO.md → README.md
- `Tool System with Core Tools` --conceptually_related_to--> `Tool System Interface-based Design`  [INFERRED]
  docs/CORE_FEATURES.md → README.md
- `CLI Entry & Argument Parsing` --conceptually_related_to--> `ADR-001: Language Selection (Go)`  [INFERRED]
  docs/CORE_FEATURES.md → docs/TECH_DECISIONS.md

## Hyperedges (group relationships)
- **Phase 1 Core Tools Implementation (6 tools)** — core_features_tool_system, arch_analysis_tool_system, go_module_tools_package, testing_coverage_goals [EXTRACTED 0.95]
- **Query Loop Architecture (API + Tools + Context)** — core_features_query_loop, core_features_api_client, core_features_tool_system, arch_analysis_query_flow [EXTRACTED 0.95]
- **Technical Decision Framework (10 ADRs)** — tech_decisions_adr001, tech_decisions_adr002, tech_decisions_adr003, tech_decisions_adr005, tech_decisions_adr008 [EXTRACTED 0.90]

## Communities

### Community 0 - "Anthropic API Client"
Cohesion: 0.06
Nodes (12): Client, Event, MockClient, StreamEvent, BuildContext(), ContextInfo, getGitBranch(), getGitStatus() (+4 more)

### Community 1 - "Glob Tool Testing"
Cohesion: 0.07
Nodes (2): AssertContains(), contains()

### Community 2 - "Integration Testing"
Cohesion: 0.07
Nodes (2): ReadTool, WriteTool

### Community 3 - "Streaming Message Events"
Cohesion: 0.1
Nodes (18): ContentBlockDeltaEvent, ContentBlockStartEvent, MessageDeltaEvent, MessageStartEvent, Behavior, Checker, Decision, Mode (+10 more)

### Community 4 - "Bash Tool Implementation"
Cohesion: 0.08
Nodes (2): BashTool, EditTool

### Community 5 - "Application Core"
Cohesion: 0.12
Nodes (8): ValidationError, InitialModelWithSessionAndSettings(), InitialModelWithSettings(), Model, PermissionDialog, QueryResultMsg, StreamEventMsg, FormatValidationErrors()

### Community 6 - "Configuration Management"
Cohesion: 0.14
Nodes (14): APISettings, DefaultSettings(), Load(), loadGlobalSettings(), loadProjectSettings(), loadSettingsFile(), mergeSettings(), PermissionRule (+6 more)

### Community 7 - "Config Unit Tests"
Cohesion: 0.1
Nodes (0): 

### Community 8 - "Grep Tool Implementation"
Cohesion: 0.16
Nodes (8): compileIncludePattern(), countFiles(), countMatches(), formatGrepResults(), GrepMatch, GrepTool, joinMatches(), walkDirectory()

### Community 9 - "User Input System"
Cohesion: 0.19
Nodes (3): NewInput(), Input, InputMode

### Community 10 - "Permission System"
Cohesion: 0.12
Nodes (0): 

### Community 11 - "CLI Command Handling"
Cohesion: 0.19
Nodes (5): CompactionResult, Compactor, Option, ManualCompact(), NewCompactor()

### Community 12 - "File Tools (Read/Write)"
Cohesion: 0.13
Nodes (0): 

### Community 13 - "Tool Registry & Dispatch"
Cohesion: 0.13
Nodes (0): 

### Community 14 - "Message Processing Pipeline"
Cohesion: 0.14
Nodes (0): 

### Community 15 - "Component 15"
Cohesion: 0.14
Nodes (0): 

### Community 16 - "Component 16"
Cohesion: 0.15
Nodes (2): Error, ErrorType

### Community 17 - "Component 17"
Cohesion: 0.15
Nodes (0): 

### Community 18 - "Component 18"
Cohesion: 0.18
Nodes (3): formatFileList(), GlobTool, joinFiles()

### Community 19 - "Component 19"
Cohesion: 0.17
Nodes (0): 

### Community 20 - "Component 20"
Cohesion: 0.36
Nodes (9): APIAuthenticationError(), APIConnectionError(), APIErrorFromStatusCode(), APIInvalidResponseError(), APIModelNotFoundError(), APIRateLimitError(), APITimeoutError(), NewAPIError() (+1 more)

### Community 21 - "Component 21"
Cohesion: 0.33
Nodes (10): NewToolError(), ToolCommandError(), ToolExecutionError(), ToolFileNotFoundError(), ToolInputValidationError(), ToolInvalidPathError(), ToolNotFoundError(), ToolPermissionDeniedError() (+2 more)

### Community 22 - "Component 22"
Cohesion: 0.29
Nodes (8): CleanupOldSessions(), generateUUID(), GetLastSession(), getSessionDir(), LoadSession(), Metadata, NewSession(), Session

### Community 23 - "Component 23"
Cohesion: 0.2
Nodes (5): ContentParam, MessageParam, Request, SystemBlock, ToolParam

### Community 24 - "Component 24"
Cohesion: 0.22
Nodes (0): 

### Community 25 - "Component 25"
Cohesion: 0.29
Nodes (0): 

### Community 26 - "Component 26"
Cohesion: 0.47
Nodes (1): Engine

### Community 27 - "Component 27"
Cohesion: 0.33
Nodes (1): Registry

### Community 28 - "Component 28"
Cohesion: 0.4
Nodes (1): ContentBlock

### Community 29 - "Component 29"
Cohesion: 0.4
Nodes (2): Message, MessageType

### Community 30 - "Component 30"
Cohesion: 0.4
Nodes (5): Core Query Flow with Tool Integration, CLI Entry & Argument Parsing, Query Loop Implementation, ADR-001: Language Selection (Go), ADR-005: Channel-based Streaming

### Community 31 - "Component 31"
Cohesion: 0.4
Nodes (5): Tool System with Registry Pattern, Tool System with Core Tools, internal/tools Package with Tool Interface, Tool System Interface-based Design, ADR-008: Goroutine Pool Tool Execution

### Community 32 - "Component 32"
Cohesion: 0.67
Nodes (0): 

### Community 33 - "Component 33"
Cohesion: 0.67
Nodes (0): 

### Community 34 - "Component 34"
Cohesion: 0.67
Nodes (1): Usage

### Community 35 - "Component 35"
Cohesion: 0.67
Nodes (3): Anthropic API Streaming Client, internal/api Package Structure, ADR-003: Anthropic-only API Support

### Community 36 - "Component 36"
Cohesion: 0.67
Nodes (3): Permission System with Multiple Modes, Permission System Implementation, internal/permission Package Structure

### Community 37 - "Component 37"
Cohesion: 1.0
Nodes (0): 

### Community 38 - "Component 38"
Cohesion: 1.0
Nodes (0): 

### Community 39 - "Component 39"
Cohesion: 2.0
Nodes (0): 

### Community 40 - "Component 40"
Cohesion: 1.0
Nodes (2): Core Features Overview, 27 Completed Features Across 3 Phases

### Community 41 - "Component 41"
Cohesion: 1.0
Nodes (2): REPL Screen with Bubble Tea, ADR-002: Bubble Tea TUI Framework

### Community 42 - "Component 42"
Cohesion: 1.0
Nodes (2): Entry Points & Bootstrap Sequence, Context Building with Git Integration

### Community 43 - "Component 43"
Cohesion: 1.0
Nodes (2): Testing Framework with 127 Test Cases, Testing Strategy with 127 Test Cases

### Community 44 - "Component 44"
Cohesion: 1.0
Nodes (0): 

### Community 45 - "Component 45"
Cohesion: 1.0
Nodes (1): Unified Error Types System

### Community 46 - "Component 46"
Cohesion: 1.0
Nodes (1): Context Compaction System

### Community 47 - "Component 47"
Cohesion: 1.0
Nodes (1): Session Storage in JSONL Format

### Community 48 - "Component 48"
Cohesion: 1.0
Nodes (1): Message Rendering Component

### Community 49 - "Component 49"
Cohesion: 1.0
Nodes (1): Multi-line Input with History Navigation

### Community 50 - "Component 50"
Cohesion: 1.0
Nodes (1): Code Coverage Goals by Module

### Community 51 - "Component 51"
Cohesion: 1.0
Nodes (1): 6-Item Future Optimization Roadmap

### Community 52 - "Component 52"
Cohesion: 1.0
Nodes (1): 12 Issues with Status Tracking

### Community 53 - "Component 53"
Cohesion: 1.0
Nodes (1): GitHub Repository Configuration Guide

## Knowledge Gaps
- **61 isolated node(s):** `MessageStartEvent`, `ContentBlockStartEvent`, `ContentBlockDeltaEvent`, `MessageDeltaEvent`, `Event` (+56 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Component 37`** (2 nodes): `main.go`, `main()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 38`** (2 nodes): `claudemd.go`, `findCLAUDEMDFiles()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 39`** (2 nodes): `uuid.go`, `generateUUID()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 40`** (2 nodes): `Core Features Overview`, `27 Completed Features Across 3 Phases`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 41`** (2 nodes): `REPL Screen with Bubble Tea`, `ADR-002: Bubble Tea TUI Framework`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 42`** (2 nodes): `Entry Points & Bootstrap Sequence`, `Context Building with Git Integration`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 43`** (2 nodes): `Testing Framework with 127 Test Cases`, `Testing Strategy with 127 Test Cases`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 44`** (1 nodes): `setup_git_sync.ps1`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 45`** (1 nodes): `Unified Error Types System`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 46`** (1 nodes): `Context Compaction System`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 47`** (1 nodes): `Session Storage in JSONL Format`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 48`** (1 nodes): `Message Rendering Component`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 49`** (1 nodes): `Multi-line Input with History Navigation`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 50`** (1 nodes): `Code Coverage Goals by Module`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 51`** (1 nodes): `6-Item Future Optimization Roadmap`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 52`** (1 nodes): `12 Issues with Status Tracking`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Component 53`** (1 nodes): `GitHub Repository Configuration Guide`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What connects `MessageStartEvent`, `ContentBlockStartEvent`, `ContentBlockDeltaEvent` to the rest of the system?**
  _61 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Anthropic API Client` be split into smaller, more focused modules?**
  _Cohesion score 0.06 - nodes in this community are weakly interconnected._
- **Should `Glob Tool Testing` be split into smaller, more focused modules?**
  _Cohesion score 0.07 - nodes in this community are weakly interconnected._
- **Should `Integration Testing` be split into smaller, more focused modules?**
  _Cohesion score 0.07 - nodes in this community are weakly interconnected._
- **Should `Streaming Message Events` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `Bash Tool Implementation` be split into smaller, more focused modules?**
  _Cohesion score 0.08 - nodes in this community are weakly interconnected._
- **Should `Application Core` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._