package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func closeOpenOutputIfNotStdout() {
	if outputF.file != nil && outputF.file != os.Stdout {
		_ = outputF.file.Close()
	}
}

type globalsSnapshot struct {
	pointerReceiver bool
	maxDepth        int
	method          string
	types           typesVal
	skips           skipsVal
	buildTags       buildTagsVal
	output          outputVal
}

func captureGlobals() globalsSnapshot {
	return globalsSnapshot{
		pointerReceiver: *pointerReceiverF,
		maxDepth:        *maxDepthF,
		method:          *methodF,
		types:           append(typesVal(nil), typesF...),
		skips:           cloneSkips(skipsF),
		buildTags:       append(buildTagsVal(nil), buildTagsF...),
		output:          outputF,
	}
}

func restoreGlobals(s globalsSnapshot) {
	closeOpenOutputIfNotStdout()
	*pointerReceiverF = s.pointerReceiver
	*maxDepthF = s.maxDepth
	*methodF = s.method
	typesF = append(typesVal(nil), s.types...)
	skipsF = cloneSkips(s.skips)
	buildTagsF = append(buildTagsVal(nil), s.buildTags...)
	outputF = s.output
}

func resetGlobalsForConfigTest() {
	closeOpenOutputIfNotStdout()
	*pointerReceiverF = false
	*maxDepthF = 0
	*methodF = "DeepCopy"
	typesF = nil
	skipsF = nil
	buildTagsF = nil
	outputF = outputVal{}
}

// configTestCLI is the simulated CLI state (package-level flags) before merging the config file.
type configTestCLI struct {
	PointerReceiver *bool
	MaxDepth        *int
	Method          *string
	Types           typesVal
	Skips           skipsVal
	BuildTags       buildTagsVal
	OutputBasename  string // basename under t.TempDir(); used when "o" is in flagsSetOnCLI
}

// configTestWant is the expected flag state after loadConfigFile.
type configTestWant struct {
	Pointer    *bool
	MaxDepth   *int
	Method     *string
	Types      typesVal
	Skips      skipsVal
	BuildTags  buildTagsVal
	OutputName string // empty = stdout
}

func cloneSkips(s skipsVal) skipsVal {
	if len(s) == 0 {
		return nil
	}
	out := make(skipsVal, len(s))
	for i, m := range s {
		out[i] = make(map[string]struct{}, len(m))
		for k := range m {
			out[i][k] = struct{}{}
		}
	}
	return out
}

func applyConfigTestCLI(t *testing.T, flags map[string]struct{}, cli configTestCLI) (wantOutputName string) {
	t.Helper()
	if len(flags) == 0 {
		return ""
	}
	if flagWasSetOnCLI(flags, "pointer-receiver") && cli.PointerReceiver != nil {
		*pointerReceiverF = *cli.PointerReceiver
	}
	if flagWasSetOnCLI(flags, "maxdepth") && cli.MaxDepth != nil {
		*maxDepthF = *cli.MaxDepth
	}
	if flagWasSetOnCLI(flags, "method") && cli.Method != nil {
		*methodF = *cli.Method
	}
	if flagWasSetOnCLI(flags, "type") && len(cli.Types) > 0 {
		typesF = append(typesVal(nil), cli.Types...)
	}
	if flagWasSetOnCLI(flags, "skip") && len(cli.Skips) > 0 {
		skipsF = cloneSkips(cli.Skips)
	}
	if flagWasSetOnCLI(flags, "tags") && len(cli.BuildTags) > 0 {
		buildTagsF = append(buildTagsVal(nil), cli.BuildTags...)
	}
	if flagWasSetOnCLI(flags, "o") && cli.OutputBasename != "" {
		p := filepath.Join(t.TempDir(), cli.OutputBasename)
		if err := outputF.Set(p); err != nil {
			t.Fatal(err)
		}
		return p
	}
	return ""
}

type configTestExtra struct {
	configPath     string
	wantMethod     *string
	wantOutputName string
}

// cliFlagsSet returns a set of flag names as loadConfigFile expects for flagsSetOnCLI.
func cliFlagsSet(names ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(names))
	for _, n := range names {
		m[n] = struct{}{}
	}
	return m
}

func writeTempConfigYAML(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cfg.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func assertGlobalsMatchWant(t *testing.T, want configTestWant, wantOutputName string) {
	t.Helper()
	if want.Pointer != nil && *pointerReceiverF != *want.Pointer {
		t.Errorf("pointerReceiverF = %v, want %v", *pointerReceiverF, *want.Pointer)
	}
	if want.MaxDepth != nil && *maxDepthF != *want.MaxDepth {
		t.Errorf("maxDepthF = %v, want %v", *maxDepthF, *want.MaxDepth)
	}
	if want.Method != nil && *methodF != *want.Method {
		t.Errorf("methodF = %v, want %v", *methodF, *want.Method)
	}
	if want.Types != nil {
		if diff := cmp.Diff(typesF, want.Types); diff != "" {
			t.Errorf("typesF (-got +want):\n%s", diff)
		}
	}
	if want.Skips != nil {
		if diff := cmp.Diff(skipsF, want.Skips); diff != "" {
			t.Errorf("skipsF (-got +want):\n%s", diff)
		}
	}
	if want.BuildTags != nil {
		if diff := cmp.Diff(buildTagsF, want.BuildTags); diff != "" {
			t.Errorf("buildTagsF (-got +want):\n%s", diff)
		}
	}
	if wantOutputName != "" && outputF.name != wantOutputName {
		t.Errorf("outputF.name = %q, want %q", outputF.name, wantOutputName)
	}
}

func ptr[T any](v T) *T {
	return &v
}

func Test_loadConfigFile(t *testing.T) {
	tests := []struct {
		name          string
		configYAML    string
		configPath    string // use as-is when configYAML == "" and non-empty (missing file)
		flagsSetOnCLI map[string]struct{}
		cli           configTestCLI
		prepare       func(t *testing.T) configTestExtra // optional: dynamic config path / want overrides
		wantErr       bool
		want          configTestWant
	}{
		{
			name:    "empty path is no-op",
			wantErr: false,
		},
		{
			name:       "missing file",
			configPath: filepath.Join(t.TempDir(), "nope.yaml"),
			wantErr:    true,
		},
		{
			name:       "invalid YAML",
			configYAML: "invalid: yaml: [",
			wantErr:    true,
		},
		{
			name: "applies all fields when no CLI flags",
			configYAML: `pointer-receiver: true
maxdepth: 5
method: Clone
type:
  - A
  - B
skip:
  - Field1,Field2
  - Field3
build-tags:
  - t1
  - t2`,
			wantErr: false,
			want: configTestWant{
				Pointer:  ptr(true),
				MaxDepth: ptr(5),
				Method:   ptr("Clone"),
				Types:    typesVal{"A", "B"},
				Skips: skipsVal{
					{"Field1": {}, "Field2": {}},
					{"Field3": {}},
				},
				BuildTags: buildTagsVal{"t1", "t2"},
			},
		},
		{
			name: "CLI method flag is not overwritten by config",
			configYAML: `method: FromConfig
maxdepth: 3`,
			flagsSetOnCLI: cliFlagsSet("method"),
			cli:           configTestCLI{Method: ptr("KeepCLI")},
			wantErr:       false,
			want: configTestWant{
				MaxDepth: ptr(3),
				Method:   ptr("KeepCLI"),
			},
		},
		{
			name: "CLI maxdepth flag is not overwritten by config",
			configYAML: `maxdepth: 99
method: M`,
			flagsSetOnCLI: cliFlagsSet("maxdepth"),
			cli:           configTestCLI{MaxDepth: ptr(7)},
			wantErr:       false,
			want: configTestWant{
				Method:   ptr("M"),
				MaxDepth: ptr(7),
			},
		},
		{
			name: "CLI type flag is not overwritten by config",
			configYAML: `type:
  - FromYAML
pointer-receiver: true`,
			flagsSetOnCLI: cliFlagsSet("type"),
			cli:           configTestCLI{Types: typesVal{"FromCLI"}},
			wantErr:       false,
			want: configTestWant{
				Pointer: ptr(true),
				Types:   typesVal{"FromCLI"},
			},
		},
		{
			name: "CLI skip flag is not overwritten by config",
			configYAML: `skip:
  - X
pointer-receiver: true`,
			flagsSetOnCLI: cliFlagsSet("skip"),
			cli:           configTestCLI{Skips: skipsVal{{"Y": {}}}},
			wantErr:       false,
			want: configTestWant{
				Pointer: ptr(true),
				Skips:   skipsVal{{"Y": {}}},
			},
		},
		{
			name: "CLI tags flag is not overwritten by config",
			configYAML: `build-tags:
  - fromyaml
method: Z`,
			flagsSetOnCLI: cliFlagsSet("tags"),
			cli:           configTestCLI{BuildTags: buildTagsVal{"fromcli"}},
			wantErr:       false,
			want: configTestWant{
				Method:    ptr("Z"),
				BuildTags: buildTagsVal{"fromcli"},
			},
		},
		{
			name: "CLI o flag is not overwritten by config",
			configYAML: `output: "/yaml-only-path/should-be-ignored.go"
method: FromYAML`,
			flagsSetOnCLI: cliFlagsSet("o"),
			cli:           configTestCLI{OutputBasename: "from_cli.go"},
			wantErr:       false,
			want: configTestWant{
				Method: ptr("FromYAML"),
			},
		},
		{
			name:       "output path under temp dir",
			configYAML: "",
			prepare: func(t *testing.T) configTestExtra {
				out := filepath.Join(t.TempDir(), "out.go")
				cfg := filepath.Join(t.TempDir(), "cfg.yaml")
				yaml := fmt.Sprintf("output: %q\nmethod: OutTest\n", out)
				if err := os.WriteFile(cfg, []byte(yaml), 0o644); err != nil {
					t.Fatal(err)
				}
				return configTestExtra{
					configPath:     cfg,
					wantMethod:     ptr("OutTest"),
					wantOutputName: out,
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snap := captureGlobals()
			defer restoreGlobals(snap)

			resetGlobalsForConfigTest()

			var path string
			switch {
			case tt.configYAML != "":
				path = writeTempConfigYAML(t, tt.configYAML)
			case tt.configPath != "":
				path = tt.configPath
			default:
				path = ""
			}

			extra := configTestExtra{}
			if tt.prepare != nil {
				extra = tt.prepare(t)
			}
			loadPath := path
			if extra.configPath != "" {
				loadPath = extra.configPath
			}

			want := tt.want
			if extra.wantMethod != nil {
				want.Method = extra.wantMethod
			}
			wantOutputName := want.OutputName
			if extra.wantOutputName != "" {
				wantOutputName = extra.wantOutputName
			}
			if out := applyConfigTestCLI(t, tt.flagsSetOnCLI, tt.cli); out != "" {
				wantOutputName = out
			}

			err := loadConfigFile(loadPath, tt.flagsSetOnCLI)
			if (err != nil) != tt.wantErr {
				t.Fatalf("loadConfigFile() err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			assertGlobalsMatchWant(t, want, wantOutputName)
		})
	}
}
