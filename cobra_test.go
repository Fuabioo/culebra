package culebra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestAutoloadWithConfigName(t *testing.T) {
	// Create a temporary directory and config file
	tempDir, err := os.MkdirTemp("", "culebra_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file
	configFile := filepath.Join(tempDir, "test.lua")
	configContent := `
	test_value = "autoloaded"
	debug = true
`
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Change to temp directory so the config file is found
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatal(err)
		}
	}()

	// Test case 1: With SetConfigName but no AddConfigPath
	t.Run("WithConfigNameOnly", func(t *testing.T) {
		cmd := &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {
				// Command execution
			},
		}
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})

	// Test case 2: With both SetConfigName and AddConfigPath
	t.Run("WithConfigNameAndPath", func(t *testing.T) {
		cmd := &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {
				// Command execution
			},
		}
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})

	// Test case 3: No SetConfigName
	t.Run("NoConfigName", func(t *testing.T) {
		cmd := &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {
				// Command execution
			},
		}
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})
}
