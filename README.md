# ai-cli

Launch Claude Code or Codex in YOLO mode from one command.

The installed command is `ai`. It starts Claude Code or Codex with their
full-permission flags, then passes all remaining arguments through unchanged.

Use it only in directories where you are comfortable running the selected agent
with broad filesystem and command permissions.

## Usage

```sh
ai
ai claude --resume
ai codex resume
ai --resume
```

Running `ai` with no provider opens an interactive choice:

```text
Select AI agent:
  1) Claude
  2) Codex
>
```

Injected arguments:

- Claude: `--dangerously-skip-permissions`
- Codex: `--dangerously-bypass-approvals-and-sandbox`

When the first argument is `claude`, `codex`, `c`, or `x`, `ai` skips the
prompt. Otherwise, it prompts first and forwards every argument to the selected
command.

That means you can pass the same arguments you would normally pass to Claude or
Codex:

```sh
ai claude --resume
ai codex resume
ai --resume
ai --continue
```

For example, `ai --resume` prompts for Claude or Codex, then runs the selected
agent with its YOLO flag plus `--resume`.

## Build

```sh
go build -o ai .
```

## Install

With Homebrew:

```sh
brew install owo-network/brew/ai
```

From source:

```sh
mkdir -p ~/.local/bin
go build -o ~/.local/bin/ai .
```

Make sure `~/.local/bin` is in `PATH`.

## Release

The GitHub Actions release workflow builds Linux, macOS, and Windows binaries in
parallel, then uploads them to the GitHub Release for the pushed tag.

Current release:

```sh
git tag v1.0.0
git push origin v1.0.0
```
