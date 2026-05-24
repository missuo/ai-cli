# ai-cli

Tiny command-line launcher for local AI coding agents.

The installed command is `ai`. It prompts for Claude or Codex when no provider
is supplied, injects the full-permission argument for the selected provider, and
passes all remaining arguments through unchanged.

## Usage

```sh
ai
ai claude --resume
ai codex resume
ai --resume
```

Injected arguments:

- Claude: `--dangerously-skip-permissions`
- Codex: `--dangerously-bypass-approvals-and-sandbox`

When the first argument is `claude`, `codex`, `c`, or `x`, `ai` skips the prompt.
Otherwise, it prompts first and forwards every argument to the selected command.

## Build

```sh
go build -o ai .
```

## Install

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
