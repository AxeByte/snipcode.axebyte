# snipcode

snipcode is a CLI tool that collects your source code, skips files per .grepattern.yaml, and outputs a single file formatted with heredocs for each source file—ideal for LLM ingestion.

## Configuration

### Local: `.grepattern.yaml` in working directory
```yaml
default_name: comp_code.txt
ignore_patterns:
  - "node_modules/**"
  - "*.test.go"
  - "docs/**"
```

### Global: `$XDG_CONFIG_HOME/snipcode/.grepattern.yaml`
Created via `snipcode init-admin`, with a comprehensive set of default ignore patterns.

## Installation
```bash
go install ${MODULE_PATH}@latest
```

## Usage
```bash
snipcode init                # create local .grepattern.yaml
snipcode init-admin          # create global config in ~/.config/snipcode/.grepattern.yaml
snipcode compile             # compile files
snipcode compile -o out.txt  # override output filename
snipcode compile --with-tree # append file tree
```

## Logging
- On `compile`: logs which config files are loaded (global and/or local)
- ➡️ per-file inclusion with char count
- logs skipping the output file if detected
- 🌳 tree listing if requested
- 📥 final summary