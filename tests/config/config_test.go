package config_test

import (
	"os"
	"testing"

	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

func TestLoadConfig(t *testing.T) {
	// Set the environment variables for testing
	os.Setenv("ENV", "testing")
	os.Setenv("PORT", "4000")
	defer os.Clearenv()

	t.Run("Load and retrieve config values", func(t *testing.T) {
		cfg := LessGo.LoadConfig()

		// Test string retrieval
		if cfg.Get("ENV", "development") != "testing" {
			t.Errorf("Expected ENV to be 'testing', got '%s'", cfg.Get("ENV", "development"))
		}

		// Test default value retrieval
		if cfg.Get("MISSING_KEY", "default") != "default" {
			t.Errorf("Expected default value for MISSING_KEY, got '%s'", cfg.Get("MISSING_KEY", "default"))
		}

		// Test integer retrieval
		if port := cfg.GetInt("PORT", 8080); port != 4000 {
			t.Errorf("Expected PORT to be 4000, got %d", port)
		}

		// Test boolean retrieval
		if debug := cfg.GetBool("DEBUG", false); debug != false {
			t.Errorf("Expected DEBUG to be false, got %t", debug)
		}

		// Test float retrieval
		if pi := cfg.GetFloat64("PI", 3.14); pi != 3.14 {
			t.Errorf("Expected PI to be 3.14, got %f", pi)
		}
	})

	t.Run("Validation of required keys", func(t *testing.T) {
		cfg := LessGo.LoadConfig()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for missing required keys, but got none")
			}
		}()

		cfg.Validate("PORT", "NON_EXISTENT_KEY")
	})

	t.Run("MergeWithDefaults", func(t *testing.T) {
		cfg := LessGo.LoadConfig()
		defaults := LessGo.Config{"HOST": "localhost", "PORT": "8080"}
		merged := cfg.MergeWithDefaults(defaults)

		if merged.Get("HOST", "") != "localhost" {
			t.Errorf("Expected HOST to be 'localhost', got '%s'", merged.Get("HOST", ""))
		}
		if merged.Get("PORT", "") != "4000" {
			t.Errorf("Expected PORT to be '4000', got '%s'", merged.Get("PORT", ""))
		}
	})

	t.Run("FilterByPrefix", func(t *testing.T) {
		cfg := LessGo.LoadConfig()
		filtered := cfg.FilterByPrefix("MYAPP_")

		if _, exists := filtered["PORT"]; exists {
			t.Errorf("Expected PORT to not exist in filtered config")
		}
	})
}

func TestReload(t *testing.T) {
	// Initial setup
	os.Setenv("ENV", "initial")
	defer os.Clearenv()

	cfg := LessGo.LoadConfig()

	if env := cfg.Get("ENV", "default"); env != "initial" {
		t.Errorf("Expected ENV to be 'initial', got '%s'", env)
	}

	// Change environment variable
	os.Setenv("ENV", "reloaded")
	cfg.Reload()

	if env := cfg.Get("ENV", "default"); env != "reloaded" {
		t.Errorf("Expected ENV to be 'reloaded', got '%s'", env)
	}
}
