# snipcode

snipcode is a CLI tool that collects your source code, skips files per .grepattern.yaml, and outputs a single file formatted with heredocs for each source file—ideal for LLM ingestion.

## Configuration: .grepattern.yaml
```yaml
default_name: comp_code.txt
ignore_patterns:
  - "node_modules/**"
  - "*.test.go"
  - "docs/**"
```

## Installation
```bash
go install ${MODULE_PATH}@latest
```

## Usage
```bash
snipcode init               # create .grepattern.yaml
snipcode compile            # compile files
snipcode compile -o out.txt # override output
snipcode compile --with-tree # append file tree
```

## Logging
- ➡️ per-file inclusion with char count
- 🌳 tree listing if requested
- 📥 final summary
