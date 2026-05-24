package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseInvocationWithProvider(t *testing.T) {
	providerName, passThroughArgs := parseInvocation([]string{"claude", "--resume", "abc"})

	if providerName != "claude" {
		t.Fatalf("provider = %q, want claude", providerName)
	}

	want := []string{"--resume", "abc"}
	if !reflect.DeepEqual(passThroughArgs, want) {
		t.Fatalf("passThroughArgs = %#v, want %#v", passThroughArgs, want)
	}
}

func TestParseInvocationWithoutProvider(t *testing.T) {
	providerName, passThroughArgs := parseInvocation([]string{"--resume", "abc"})

	if providerName != "" {
		t.Fatalf("provider = %q, want empty", providerName)
	}

	want := []string{"--resume", "abc"}
	if !reflect.DeepEqual(passThroughArgs, want) {
		t.Fatalf("passThroughArgs = %#v, want %#v", passThroughArgs, want)
	}
}

func TestBuildArgv(t *testing.T) {
	got := buildArgv(providers["codex"], []string{"resume", "abc"})
	want := []string{"codex", "--dangerously-bypass-approvals-and-sandbox", "resume", "abc"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("argv = %#v, want %#v", got, want)
	}
}

func TestProviderAliases(t *testing.T) {
	cases := map[string]string{
		"c":      "claude",
		"Claude": "claude",
		"x":      "codex",
		"Codex":  "codex",
	}

	for input, want := range cases {
		if got := normalizeProvider(input); got != want {
			t.Fatalf("normalizeProvider(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestSelectionAliases(t *testing.T) {
	cases := map[string]string{
		"1":      "claude",
		"2":      "codex",
		" c \n":  "claude",
		" Codex": "codex",
	}

	for input, want := range cases {
		if got := normalizeSelection(input); got != want {
			t.Fatalf("normalizeSelection(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestMissingProviderErrorMentionsInstall(t *testing.T) {
	err := missingProviderError(providers["claude"])
	if err == nil {
		t.Fatal("missingProviderError returned nil")
	}

	message := err.Error()
	for _, want := range []string{"Claude", "claude", "Install"} {
		if !strings.Contains(message, want) {
			t.Fatalf("error %q does not contain %q", message, want)
		}
	}
}
