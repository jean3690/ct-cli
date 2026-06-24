# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

`ct-cli` is a Go CLI tool for managing code templates and snippets. It scaffolds projects from templates and generates files into existing projects using Go's `text/template` engine with interactive prompts (survey.v2).

**Binary:** `ct-cli`  
**Module:** `github.com/jean/codeTemplateCli` (Go 1.26.4)  
**Data directory:** `~/.codeTemplate/` (override with `CT_HOME` env var)

## Build & test

```bash
go build -o ct-cli.exe .      # build
go vet ./...                   # lint
go test ./...                  # run all tests
go test ./pkg/engine/          # single package tests
```

## Architecture

```
main.go → cmd.Execute()
```

### Command tree (cobra)

```
ct-cli
├── new [template] [dir]    — scaffold a new project from a template
├── generate [template]     — generate files into the current directory (no new dir)
├── init                    — create ~/.codeTemplate/ directory structure
├── config show|set         — view/persist config keys
├── template
│   ├── list [-l lang]      — list installed templates
│   ├── info <name>         — show manifest details and variables
│   ├── import <src> [name] — import from local dir (--git for git clone)
│   ├── remove <name>
│   ├── init <name>         — scaffold a new template skeleton for authoring
│   └── edit <name>         — open template dir in $EDITOR
└── snippet
    ├── create <name>       — create interactively or via --code flag
    ├── list [-l lang]
    ├── get <name>
    ├── delete <name>
    └── search <query>
```

### Package responsibilities

| Package | Role |
|---|---|
| `cmd/` | Cobra command definitions; `root.go` wires subcommands and global flags |
| `pkg/engine/` | `TemplateEngine` interface + `textEngine` impl wrapping `text/template` with custom funcs: `upper`, `lower`, `title`, `snake`, `camel`, `kebab`, `pascal`, `now` |
| `pkg/template/` | `TemplateManifest` YAML model, `LoadManifest`/`SaveManifest`/`Validate`, template discovery/find, import from path or git, remove |
| `pkg/scaffold/` | `Scaffold()` — orchestrates manifest loading → variable collection → file rendering → postGenerate hook execution |
| `pkg/prompts/` | `CollectVariables()` — interactive prompting via survey.v2 (string/bool/choice types, regex validation) |
| `pkg/snippet/` | `Snippet` struct + `Store` (YAML file per snippet, CRUD + search with scoring) |
| `pkg/config/` | `Paths` resolution (`~/.codeTemplate/` or `$CT_HOME`), viper init with env prefix `CT`, config read/write |
| `pkg/output/` | Colored terminal output helpers (`Success`/`Error`/`Info`/`Warn`) |

### Template manifest format

Templates live as subdirectories of `~/.codeTemplate/templates/`, each with a `manifest.yaml`:

```yaml
apiVersion: v1
name: my-template
version: "1.0"
description: "..."
variables:                        # prompted interactively (string/bool/choice)
  - name: ProjectName
    type: string
    required: true
files:                            # relative to files/ dir
  - source: main.go.tmpl          # .tmpl = rendered with text/template
    dest: "{{ .ProjectName }}/main.go"
scaffoldDir: files/               # if set, walks entire dir instead of file list
postGenerate: "go mod init {{ .ProjectName }}"
```

- `.tmpl` files are rendered through Go's `text/template`; all other files are copied as-is.
- Both `files[]` and `scaffoldDir` modes support path variables (e.g. `dest: "{{ .ProjectName }}/src/main.go"`).
- `postGenerate` is a shell command run after generation, with variable substitution.

## Key patterns

- **Variable passing:** `--var Name=value` flags on `new` and `generate`, parsed by `parseVarFlags()` in `cmd/new.go`.
- **Global flags:** `--config`/`-c` (config file path), `--verbose`/`-v`, `--yes`/`-y` (skip interactive prompts). All bound to viper.
- **Interactive vs non-interactive:** `config.IsInteractive()` returns `true` unless `--yes` is set or `defaults.interactive` is false.
- **Template fallback:** If the first arg to `ct-cli new` isn't a known template but a default is configured, it's treated as the destination directory.
