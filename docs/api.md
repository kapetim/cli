# API

Import path: `github.com/kapetim/cli/src/pkg/...`. The module is stdlib-only.

| Package | Surface |
| --- | --- |
| `config` | `Defaults()`, `Load(dir)` — embedded defaults + `.cli.json` override |
| `logx` | `OK/Skip/Info/Fail`, exit codes `0/1/2/3` |
| `process` | `Run(name, args, Opts) (Out, error)`, `Available(name)` |
| `fsx` | `Files(dir, pred)`, `MarkdownFiles(dir)` |
| `table` | `Parse(md)`, `Scan(path) []Located` (marker → `Lxx-Lyy`), `SplitRow`, `IsDelimiter`, `Clean`, `BRL`, `Percent`, `Approx`, `RenderedLength` |
| `manifest` | `Scan(dir)`, `Load(path)`, `Save(path, m)`, `Diff(got, want)` |
| `lint` | `Run(dir, kinds, cfg) []error`, `Issue` |
| `validate` | `Tables(dir, manifestPath) []error` |
| `run` | `Main(Options) int` — umbrella wiring |
| `version` | `Version` |

## Example

```go
package main

import (
	"os"

	"github.com/kapetim/cli/src/pkg/run"
)

func main() {
	os.Exit(run.Main(run.Options{
		RepoDir:  ".",
		Lint:     []string{"markdown", "shell"},
		Validate: []string{"tables"},
	}))
}
```

## Table manifest (`tables.json`)

`cli scan tables --manifest <path>` writes it; `cli validate tables --manifest
<path>` enforces it. Stored in the consumer repo.

```json
{
  "src/health/diet/status.md": [
    { "begin": 5, "end": 15, "lines": "L5-L15", "cols": 2, "rows": 9 }
  ]
}
```
