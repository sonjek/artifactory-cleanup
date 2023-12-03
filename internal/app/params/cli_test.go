package params

import (
	"flag"
	"os"
	"testing"
)

func TestInitCLIParamsNoConfigFile(t *testing.T) {
	if _, err := InitCLIParams(); err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestInitCLIParamsConfigFile(t *testing.T) {
	oldCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldCommandLine }()

	newCommandLine := flag.NewFlagSet("test", flag.ExitOnError)
	flag.CommandLine = newCommandLine

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"test", "-config", "config.json"}

	if _, err := InitCLIParams(); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestInitCLIParamsProtectedDestroy(t *testing.T) {
	oldCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldCommandLine }()

	newCommandLine := flag.NewFlagSet("test", flag.ExitOnError)
	flag.CommandLine = newCommandLine

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"test", "-config", "config.json", "-before", "1", "-destroy"}

	if _, err := InitCLIParams(); err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestInitCLIParamsDestroy(t *testing.T) {
	oldCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldCommandLine }()

	newCommandLine := flag.NewFlagSet("test", flag.ExitOnError)
	flag.CommandLine = newCommandLine

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"test", "-config", "config.json", "-before", "2", "-destroy"}

	if _, err := InitCLIParams(); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestInitCLIParamsDryModeAndBefore1(t *testing.T) {
	oldCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldCommandLine }()

	newCommandLine := flag.NewFlagSet("test", flag.ExitOnError)
	flag.CommandLine = newCommandLine

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"test", "-config", "config.json", "-before", "1"}

	params, err := InitCLIParams()

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if params.DeleteBefore != "1w" {
		t.Errorf("Expected Before to be '1w', but got '%s'", params.DeleteBefore)
	}

}
