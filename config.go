package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type config struct {
	PointerReceiver *bool   `yaml:"pointer-receiver,omitempty"`
	MaxDepth        *int    `yaml:"maxdepth,omitempty"`
	Method          *string `yaml:"method,omitempty"`

	Types      []string `yaml:"type,omitempty"`
	Skips      []string `yaml:"skip,omitempty"`
	OutputPath *string  `yaml:"output,omitempty"`
	BuildTags  []string `yaml:"build-tags,omitempty"`
}

func loadConfig() error {
	flagsSetOnCLI := make(map[string]struct{})
	flag.Visit(func(f *flag.Flag) { flagsSetOnCLI[f.Name] = struct{}{} })
	return loadConfigFile(strings.TrimSpace(*configFileF), flagsSetOnCLI)
}

// loadConfigFile reads YAML from path and merges into package-level flags.
// flagsSetOnCLI is the set of flag names that appeared on the command line; those are not overwritten from the file.
// An empty path is a no-op.
func loadConfigFile(configPath string, flagsSetOnCLI map[string]struct{}) error {
	if strings.TrimSpace(configPath) == "" {
		return nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()

	var cfg config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("decoding config file: %w", err)
	}

	mergePtr(flagsSetOnCLI, "pointer-receiver", cfg.PointerReceiver, pointerReceiverF)
	mergePtr(flagsSetOnCLI, "maxdepth", cfg.MaxDepth, maxDepthF)
	mergePtr(flagsSetOnCLI, "method", cfg.Method, methodF)

	if len(cfg.Types) > 0 && !flagWasSetOnCLI(flagsSetOnCLI, "type") {
		typesF = typesVal(cfg.Types)
	}
	if len(cfg.Skips) > 0 && !flagWasSetOnCLI(flagsSetOnCLI, "skip") {
		skipsF = skipsVal{}
		for _, skip := range cfg.Skips {
			if err := skipsF.Set(skip); err != nil {
				return fmt.Errorf("parsing skip value: %w", err)
			}
		}
	}
	if cfg.OutputPath != nil && !flagWasSetOnCLI(flagsSetOnCLI, "o") {
		if err := outputF.Set(*cfg.OutputPath); err != nil {
			return fmt.Errorf("setting output: %w", err)
		}
	}
	if len(cfg.BuildTags) > 0 && !flagWasSetOnCLI(flagsSetOnCLI, "tags") {
		buildTagsF = buildTagsVal(cfg.BuildTags)
	}

	return nil
}

func flagWasSetOnCLI(flagsSetOnCLI map[string]struct{}, name string) bool {
	_, ok := flagsSetOnCLI[name]
	return ok
}

func mergePtr[T any](flagsSetOnCLI map[string]struct{}, name string, src, dst *T) {
	if src == nil || flagWasSetOnCLI(flagsSetOnCLI, name) {
		return
	}
	*dst = *src
}
