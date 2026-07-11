package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var version = "1.0.0"

type provider struct {
	name      string
	command   string
	extraArgs []string
}

var providerOrder = []string{"claude", "codex", "grok"}

var providers = map[string]provider{
	"claude": {
		name:      "Claude",
		command:   "claude",
		extraArgs: []string{"--dangerously-skip-permissions"},
	},
	"codex": {
		name:      "Codex",
		command:   "codex",
		extraArgs: []string{"--dangerously-bypass-approvals-and-sandbox"},
	},
	"grok": {
		name:      "Grok",
		command:   "grok",
		extraArgs: []string{"--always-approve"},
	},
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "ai: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, in *os.File, errOut *os.File) error {
	if shouldShowHelp(args) {
		printHelp(errOut)
		return nil
	}

	if shouldShowVersion(args) {
		fmt.Fprintf(errOut, "ai %s\n", version)
		return nil
	}

	providerName, passThroughArgs := parseInvocation(args)
	if providerName == "" {
		selected, err := promptProvider(in, errOut)
		if err != nil {
			return err
		}
		providerName = selected
	}

	p, ok := providers[providerName]
	if !ok {
		return fmt.Errorf("unknown provider %q", providerName)
	}

	return execProvider(p, passThroughArgs)
}

func parseInvocation(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}

	name := normalizeProvider(args[0])
	if name == "" {
		return "", args
	}

	return name, args[1:]
}

func normalizeProvider(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "claude", "c":
		return "claude"
	case "codex", "x":
		return "codex"
	case "grok", "g":
		return "grok"
	default:
		return ""
	}
}

func normalizeSelection(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1":
		return "claude"
	case "2":
		return "codex"
	case "3":
		return "grok"
	default:
		return normalizeProvider(value)
	}
}

func promptProvider(in *os.File, errOut *os.File) (string, error) {
	if !hasInstalledProvider() {
		return "", missingAllProvidersError()
	}

	reader := bufio.NewReader(in)

	for {
		fmt.Fprintln(errOut, "Select AI agent:")
		for idx, name := range providerOrder {
			p := providers[name]
			status := ""
			if !isProviderInstalled(p) {
				status = " (not installed)"
			}
			fmt.Fprintf(errOut, "  %d) %s%s\n", idx+1, p.name, status)
		}
		fmt.Fprint(errOut, "> ")

		raw, err := reader.ReadString('\n')
		if err != nil && len(raw) == 0 {
			return "", errors.New("failed to read selection")
		}

		name := normalizeSelection(raw)
		if name == "" {
			fmt.Fprintln(errOut, "Please choose 1/claude, 2/codex, or 3/grok.")
			continue
		}

		p := providers[name]
		if !isProviderInstalled(p) {
			return "", missingProviderError(p)
		}

		return name, nil
	}
}

func execProvider(p provider, passThroughArgs []string) error {
	path, err := exec.LookPath(p.command)
	if err != nil {
		return missingProviderError(p)
	}

	argv := buildArgv(p, passThroughArgs)
	return runProvider(path, argv)
}

func hasInstalledProvider() bool {
	for _, name := range providerOrder {
		if isProviderInstalled(providers[name]) {
			return true
		}
	}
	return false
}

func isProviderInstalled(p provider) bool {
	_, err := exec.LookPath(p.command)
	return err == nil
}

func missingProviderError(p provider) error {
	return fmt.Errorf("%s command %q was not found in PATH. Install %s and make sure %q is available, or choose another installed agent", p.name, p.command, p.name, p.command)
}

func missingAllProvidersError() error {
	return errors.New("no supported AI agent was found in PATH. Install Claude Code (`claude`), Codex (`codex`), or Grok (`grok`) and try again")
}

func buildArgv(p provider, passThroughArgs []string) []string {
	argv := make([]string, 0, 1+len(p.extraArgs)+len(passThroughArgs))
	argv = append(argv, p.command)
	argv = append(argv, p.extraArgs...)
	argv = append(argv, passThroughArgs...)
	return argv
}

func shouldShowHelp(args []string) bool {
	if len(args) != 1 {
		return false
	}

	switch args[0] {
	case "-h", "--help", "help":
		return true
	default:
		return false
	}
}

func shouldShowVersion(args []string) bool {
	if len(args) != 1 {
		return false
	}

	switch args[0] {
	case "-v", "--version", "version":
		return true
	default:
		return false
	}
}

func printHelp(errOut *os.File) {
	fmt.Fprintln(errOut, `Usage:
  ai
  ai claude [args...]
  ai codex [args...]
  ai grok [args...]
  ai [args...]

Examples:
  ai
  ai claude --resume
  ai codex resume
  ai grok -c
  ai --resume

When no provider is supplied, ai prompts for Claude, Codex, or Grok, then
passes every remaining argument to the selected command.

Injected arguments:
  Claude: --dangerously-skip-permissions
  Codex:  --dangerously-bypass-approvals-and-sandbox
  Grok:   --always-approve`)
}
