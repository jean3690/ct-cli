package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Paths resolves and returns all storage paths.
type Paths struct {
	ConfigDir    string
	ConfigFile   string
	TemplatesDir string
	SnippetsDir  string
}

// DefaultDir returns the default config directory ~/.codeTemplate.
func DefaultDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	return filepath.Join(home, ".codeTemplate"), nil
}

// InitPaths initializes the ~/.codeTemplate directory and returns resolved paths.
func InitPaths() (*Paths, error) {
	baseDir, err := DefaultDir()
	if err != nil {
		return nil, err
	}

	p := &Paths{
		ConfigDir:    baseDir,
		ConfigFile:   filepath.Join(baseDir, "config.yaml"),
		TemplatesDir: filepath.Join(baseDir, "templates"),
		SnippetsDir:  filepath.Join(baseDir, "snippets"),
	}

	// Allow override via environment variable.
	if envDir := os.Getenv("CT_HOME"); envDir != "" {
		p.ConfigDir = envDir
		p.ConfigFile = filepath.Join(envDir, "config.yaml")
		p.TemplatesDir = filepath.Join(envDir, "templates")
		p.SnippetsDir = filepath.Join(envDir, "snippets")
	}

	return p, nil
}

// EnsureDirs creates all required directories.
func (p *Paths) EnsureDirs() error {
	for _, dir := range []string{p.ConfigDir, p.TemplatesDir, p.SnippetsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}
	}
	return nil
}

// InitViper sets up viper with config file and environment bindings.
func InitViper(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		defaultDir, err := DefaultDir()
		if err != nil {
			return err
		}
		viper.AddConfigPath(defaultDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("CT")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("defaults.interactive", true)
	viper.SetDefault("defaults.license", "MIT")
		viper.SetDefault("defaults.template", "")
	if runtime.GOOS == "windows" {
		viper.SetDefault("editor", "notepad")
	} else {
		viper.SetDefault("editor", os.Getenv("EDITOR"))
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("reading config: %w", err)
		}
		// Config file not found is OK — defaults will be used.
	}

	return nil
}

// IsInteractive returns whether interactive mode is enabled.
func IsInteractive() bool {
	return viper.GetBool("defaults.interactive") && !viper.GetBool("yes")
}

// ReadConfig reads the raw config file content.
func ReadConfig(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteConfig writes raw content to the config file.
func WriteConfig(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
