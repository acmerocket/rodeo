# rodeo
Format JSON based on markdown templates.

Designed for use with [goat](https://github.com/bluesky-social/indigo/tree/main/cmd/goat)

## tl;dr

    go install github.com/acmerocket/rodeo

    goat firehose | rodeo

## Build

Accept input as well-formed chunks (lines, seperated by \n) of JSON to be parsed and applied to a template.

## Development Log

### 2025-04-12
Starting work on this to see it I can make an adequite filter/formtter for [goat](https://github.com/bluesky-social/indigo/tree/main/cmd/goat) that prints selected records with fancy terminal formatting based on [lipgloss](https://github.com/charmbracelet/lipgloss).

Today's goals:
1. Fork [`gum`](https://github.com/charmbracelet/gum) as a basis for "terminal tailing commands" that are generating JSON.
2. Get that working with sample data, with "default" format (field names and values in table)
3. Get that working for `goat` output: The rough idea is `goat --some-params | gum json --some-params`

Fork created: https://github.com/philion/gum

Looking over the lipgloss, https://github.com/charmbracelet/lipgloss/blob/master/examples/table/ansi/main.go seems a good place to start:
```go
import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render

	t := table.New()
	t.Row("Bubble Tea", s("Milky"))
	t.Row("Milk Tea", s("Also milky"))
	t.Row("Actual milk", s("Milky as well"))
	fmt.Println(t.Render())
}
```

On the `gum` side, starting with a copy of `table` subcommand, but replacing the CSV parsing with json.

Got a simple rudmentary system working. `gum` too much overhead. working in `tinker.go`, will convert to a fresh project when a good name occurs to me.

Step 3 done, enough: `goat firehose | go run tinker.go` produces expected results.

Moving on to what "template" means. Looking at https://pkg.go.dev/text/template overlaid on a markdown file.

Instead of table:
1. parse json to map
2. look for type
3. if type, look for template
4. if template, apply and print
5. else apply templates/default.md

Markdown?

https://github.com/charmbracelet/glamour

Integrated glamour and templates.

Created new github project: https://github.com/acmerocket/rodeo
