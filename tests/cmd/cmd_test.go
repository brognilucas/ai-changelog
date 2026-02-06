package cmd_test

import (
	"testing"

	"github.com/lucasbrogni/ai-changelog/cmd"
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
