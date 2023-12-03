package params

import (
	"os"
	"testing"
)

func TestParceConfigFileNotExistedConfig(t *testing.T) {
	if _, err := ParceConfigFile("NotExistedConfig"); err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestParceConfigFileProvided(t *testing.T) {
	configFileContent := `
type: docker
repos:
  - my-repo
cleanupPatterns:
  - pattern
excludePatterns:
  - exclude
`
	tmpfile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configFileContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	config, err := ParceConfigFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if config.Type != "docker" {
		t.Errorf("Expected 'Type' to be 'docker', but got '%s'", config.Type)
	}
	if config.Repos[0] != "my-repo" {
		t.Errorf("Expected 'Repo' to be 'my-repo', but got '%s'", config.Repos)
	}
	if len(config.CleanupPatterns) != 1 {
		t.Errorf("Expected 1 cleanup pattern, but got %d", len(config.CleanupPatterns))
	}
	if len(config.ExcludePatterns) != 1 {
		t.Errorf("Expected 1 exclude pattern, but got %d", len(config.ExcludePatterns))
	}
}
