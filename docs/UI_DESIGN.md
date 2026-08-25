# Tercode — UI Design and Customization Specification

**Status:** Draft v1.0  
**Purpose:** Define Tercode's composable terminal UI philosophy and public customization model.

## 1. Design Principle

Tercode is not a fixed TUI with themes.

It is a **terminal UI engine with a strong default environment**.

The guiding principle is:

> **Powerful by default. Limitless by configuration.**

The design philosophy is inspired by the composability and customization culture of Hyprland, Neovim, Kitty, and Alacritty without copying their interfaces.

## 2. Default Experience

`tercode` must open directly into a polished but restrained interface.

The default UI should prioritize:

- readability
- fast navigation
- agent activity visibility
- prompt accessibility
- useful project state
- low visual noise

The default design should remain recognizable as Tercode.

## 3. UI Layers

```text
UI Runtime
├── Components
├── Layout
├── Theme
├── Input
├── State Selectors
├── Actions
├── Commands
└── Event Hooks
```

## 4. Components

Initial component vocabulary:

```text
Text
Box
Panel
List
Table
Tree
Input
Prompt
Chat
Markdown
Code
Diff
Progress
Status
Notification
Tabs
Split
Scroll
Agent
Tool
Session
Workspace
Git
Model
Provider
User
```

Components should be composable and independently testable.

## 5. Layout

Supported concepts should include:

```text
row
column
split
stack
tabs
overlay
panel
scroll
```

Example:

```yaml
layout:
  type: split
  children:
    - type: workspace
      width: 28%
    - type: vertical
      children:
        - type: agent
          flex: 1
        - type: prompt
          height: 3
```

## 6. Theme

Themes control presentation:

```text
colors
symbols
typography
borders
spacing
styles
```

Themes must not alter domain behavior.

## 7. Keybindings

Keybindings map user input to semantic actions.

Example:

```yaml
keybindings:
  ctrl+space: command_palette
  ctrl+b: toggle_sidebar
  ctrl+n: new_session
```

Do not bind directly to internal function names.

## 8. Commands

Commands should have stable IDs.

Examples:

```text
command_palette
new_session
switch_model
toggle_sidebar
compact_context
show_diff
```

Slash commands are a user-facing representation of commands.

## 9. State Binding

Components should consume stable selectors such as:

```text
active_agent
active_task
active_model
active_provider
git_status
workspace_tree
session_list
tool_activity
```

Components should not reach directly into arbitrary service internals.

## 10. Events and Actions

Events describe facts:

```text
task_started
tool_started
file_changed
model_streamed
task_completed
```

Actions request changes:

```text
start_task
cancel_task
approve_tool
switch_session
change_model
```

Do not confuse events with commands/actions.

## 11. Customization Levels

### Level 1
Theme

### Level 2
Layout

### Level 3
Component composition

### Level 4
Behavior and scripting

## 12. Configuration Philosophy

Configuration should describe intent rather than implementation.

Prefer:

```yaml
agent:
  show_tool_output: true
```

over implementation-specific properties.

## 13. Future Scripting

A scripting layer may provide:

- custom commands
- hooks
- workflows
- dynamic behavior

Lua is a candidate.

The scripting API must expose stable Tercode concepts.

## 14. Design Constraints

The UI must:

- remain responsive during agent execution
- handle narrow terminals
- support resizing
- degrade gracefully when optional data is unavailable
- remain usable with keyboard navigation
- avoid unnecessary animations
- preserve terminal accessibility

## 15. Design Goal

A user's Tercode should be able to become as personal as a Neovim configuration or Hyprland setup while retaining the same underlying Tercode core.
