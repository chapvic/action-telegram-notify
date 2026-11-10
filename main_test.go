package main

import (
	"strings"
	"testing"
)

func TestPrepareChatID(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		wantType string
		wantVal  any
	}{
		{"numeric", "123456789", "int64", int64(123456789)},
		{"negative supergroup", "-1001234567890", "int64", int64(-1001234567890)},
		{"username", "@mychannel", "string", "@mychannel"},
		{"non numeric", "some-chat", "string", "some-chat"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := prepareChatID(c.input)
			switch c.wantType {
			case "int64":
				v, ok := got.(int64)
				if !ok {
					t.Fatalf("expected int64, got %T", got)
				}
				if v != c.wantVal.(int64) {
					t.Fatalf("expected %v, got %v", c.wantVal, v)
				}
			case "string":
				v, ok := got.(string)
				if !ok {
					t.Fatalf("expected string, got %T", got)
				}
				if v != c.wantVal.(string) {
					t.Fatalf("expected %q, got %q", c.wantVal, v)
				}
			}
		})
	}
}

func TestTruncateMessage(t *testing.T) {
	cases := []struct {
		name  string
		input string
		limit int
		want  string
	}{
		{"short", "hello", 10, "hello"},
		{"exact", "hello", 5, "hello"},
		{"long", "hello world", 5, "hell…"},
		{"unicode", "привет мир", 7, "привет…"},
		{"empty", "", 10, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := truncateMessage(c.input, c.limit)
			if got != c.want {
				t.Fatalf("expected %q, got %q", c.want, got)
			}
		})
	}
}

func TestNormalizeStatus(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		wantStat string
		wantWarn bool
	}{
		{"success", "success", "success", false},
		{"uppercase", "SUCCESS", "success", false},
		{"with spaces", "  Failure  ", "failure", false},
		{"empty", "", "info", true},
		{"invalid", "unknown", "info", true},
		{"mixed case cancelled", "Cancelled", "cancelled", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotStat, gotWarn := normalizeStatus(c.input)
			if gotStat != c.wantStat {
				t.Fatalf("expected status %q, got %q", c.wantStat, gotStat)
			}
			if (gotWarn != "") != c.wantWarn {
				t.Fatalf("expected warning=%v, got %q", c.wantWarn, gotWarn)
			}
		})
	}
}

func TestBuildCommitLine(t *testing.T) {
	cases := []struct {
		name   string
		server string
		repo   string
		commit string
		want   string
	}{
		{
			"complete",
			"https://github.com",
			"user/repo",
			"abc123",
			"💬 Commit: [open link](https://github.com/user/repo/commit/abc123)",
		},
		{"no server", "", "user/repo", "abc123", ""},
		{"no repo", "https://github.com", "", "abc123", ""},
		{"no commit", "https://github.com", "user/repo", "", ""},
		{"all empty", "", "", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := buildCommitLine(c.server, c.repo, c.commit)
			if got != c.want {
				t.Fatalf("expected %q, got %q", c.want, got)
			}
		})
	}
}

func TestBuildTitle(t *testing.T) {
	got := buildTitle("✅", "", "success", "")
	if !strings.Contains(got, "✅") || !strings.Contains(got, "*success*") {
		t.Fatalf("unexpected title: %q", got)
	}

	got = buildTitle("🔴", "Custom Title", "failure", "MyWorkflow")
	if !strings.Contains(got, "Custom Title") {
		t.Fatalf("expected custom title, got %q", got)
	}
	if strings.Contains(got, "MyWorkflow") {
		t.Fatalf("workflow should be ignored when title set, got %q", got)
	}
}
