package cmd_test

import (
	"testing"

	"github.com/brognilucas/ai-changelog/cmd"
)

func TestRootCommandExists(t *testing.T) {
	rootCmd := cmd.NewRootCommand()

	if rootCmd == nil {
		t.Fatal("expected root command to not be nil")
	}

	if rootCmd.Use != "ai-changelog" {
		t.Errorf("expected root command Use to be 'ai-changelog', got %q", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("expected root command to have a short description")
	}
}

func TestOutputFlag(t *testing.T) {
	rootCmd := cmd.NewRootCommand()

	flag := rootCmd.PersistentFlags().Lookup("output")
	if flag == nil {
		t.Fatal("expected --output flag to exist")
	}

	if flag.DefValue != "" {
		t.Errorf("expected --output default to be empty, got %q", flag.DefValue)
	}

	if flag.Shorthand != "o" {
		t.Errorf("expected --output shorthand to be 'o', got %q", flag.Shorthand)
	}
}

func TestSinceFlag(t *testing.T) {
	rootCmd := cmd.NewRootCommand()

	flag := rootCmd.PersistentFlags().Lookup("since")
	if flag == nil {
		t.Fatal("expected --since flag to exist")
	}

	if flag.DefValue != "" {
		t.Errorf("expected --since default to be empty, got %q", flag.DefValue)
	}

	if flag.Shorthand != "s" {
		t.Errorf("expected --since shorthand to be 's', got %q", flag.Shorthand)
	}
}

func TestModelFlag(t *testing.T) {
	rootCmd := cmd.NewRootCommand()

	flag := rootCmd.PersistentFlags().Lookup("model")
	if flag == nil {
		t.Fatal("expected --model flag to exist")
	}

	if flag.DefValue != "llama3.2" {
		t.Errorf("expected --model default to be 'llama3.2', got %q", flag.DefValue)
	}

	if flag.Shorthand != "m" {
		t.Errorf("expected --model shorthand to be 'm', got %q", flag.Shorthand)
	}
}

func TestFormatFlag(t *testing.T) {
	rootCmd := cmd.NewRootCommand()

	flag := rootCmd.PersistentFlags().Lookup("format")
	if flag == nil {
		t.Fatal("expected --format flag to exist")
	}

	if flag.DefValue != "markdown" {
		t.Errorf("expected --format default to be 'markdown', got %q", flag.DefValue)
	}

	if flag.Shorthand != "f" {
		t.Errorf("expected --format shorthand to be 'f', got %q", flag.Shorthand)
	}
}

func TestVersionFlag(t *testing.T) {
	rootCmd := cmd.NewRootCommand()

	flag := rootCmd.PersistentFlags().Lookup("version")
	if flag == nil {
		t.Fatal("expected --version flag to exist")
	}

	if flag.DefValue != "" {
		t.Errorf("expected --version default to be empty, got %q", flag.DefValue)
	}

	if flag.Shorthand != "V" {
		t.Errorf("expected --version shorthand to be 'V', got %q", flag.Shorthand)
	}
}
