# Development Log

## 2025-04-14
Starting with some clean up of the project (break out main), adding a minimal test.

Goal is to:
- [x] make "post" param match "app.bsky.feed.post"
- [x] Update readme
- [ ] add initial template externalization: type | type=othertype | type=path/to/template.md

Next, use "README-driven development" to write the usage section in the README to describe how I want to system to work. Later, once the code is implemented and "--help" contains full and accurate usage, I'll copy that to useage as well.

See [README.md#Usage](README.md#Usage)

Updated! Next, add the general type matching above. And start with a test! `rodeo_test.go`

Added test for `match()`, implemented required behavior. DONE!

Next: template parsing from params.




Time to release!
```
VERSION=v0.6.0 make release
```


## 2025-04-13
Moved developer log to seperate file: devlog.md (this file)

Thinking about next steps for #rodeo
- [ ] external templates: ./templates (local search)
- [ ] user templates: ~/.config/rodeo/templates (`rodeo config` creates)
- [ ] user config: ~/.config/rodeo/settings.json (create with default)
- [ ] research filters on `goat`
- [x] Makefile for test, install, release, deploy

Basic ideas being able to easily set the

extra templates:
- snoop.md - everything in the whole doc.
- record.md - everything in the record.
- known lexicons
- collection stats: message count, type count

on exit, print report. (add flag to toggle)

commands:
* none: default piping behavior with standard template search/resolution
* config: init and set config values
    * rodeo config - setup defaults, report on non-standard settings
    * rodeo config [name] [value]
    * rodeo config unset [name]
* -t, --template: add specific template overrides
    * -t, --template template.md - all types are processed
    * -t name=template.md - set a specific template for type name
* -q, --quiet:
    * -q - quiet everything by templates specified on the CLI
    * -q [type] - silence specific types
* -v, --debug - show debug info
* [types]: list of glob types filters to process

Type matching: `app.bsky.feed.like`  would be matched by:
* `app.bsky.feed.like` - full name
* `app.`,  `app.bsky.feed.` - everything at the start
* `.feed.like`,  `.like` - everything at the env

`rodeo post` vs `rodeo .post`?

For now, just substring match: "post", "repost", etc.

string match against type. match, eval. otherwise skip.

checking out docs about releases:
- https://go.dev/doc/modules/release-workflow
- https://go.dev/doc/modules/version-numbers
- https://go.dev/doc/modules/publishing

#### more on version numbering
https://stackoverflow.com/questions/11354518/application-auto-build-versioning
> -X importpath.name=value
> Set the value of the string variable in importpath named name to

```go
package main

import "fmt"

var xyz string

func main() {
    fmt.Println(xyz)
}
```

Then:

```go
$ go run -ldflags "-X main.xyz=abc" main.go
abc
```

In order to set `main.minversion` to the build date and time when building:

```go
go build -ldflags "-X main.minversion=`date -u +.%Y%m%d.%H%M%S`" service.go
```

```go
go build -o mybinary \
  -ldflags "-X main.version=1.0.0 -X 'main.build=$(date)'" \
  main.go
```

> https://go.dev/wiki/GcToolchainTricks

The gc toolchain linker, [cmd/link](https://pkg.go.dev/cmd/link), provides a `-X` option that may be used to record arbitrary information in a Go string variable at link time. The format is `-X importpath.name=val`. Here `importpath` is the name used in an import statement for the package (or `main` for the main package), `name` is the name of the string variable defined in the package, and `val` is the string you want to set that variable to. When using the go tool, use its `-ldflags` option to pass the `-X` option to the linker.

Let’s suppose this file is part of the package `company/buildinfo`:

```go
package buildinfo

var BuildTime string
```

You can build the program using this package using `go build -ldflags="-X 'company/buildinfo.BuildTime=$(date)'"` to record the build time in the string. (The use of `$(date)` assumes you are using a Unix-style shell.)

The string variable must exist, it must be a variable, not a constant, and its value must not be initialized by a function call. There is no warning for using the wrong name in the `-X` option. You can often find the name to use by running `go tool nm` on the program, but that will fail if the package name has any non-ASCII characters, or a `"` or `%` character.

^^^ That's how to build it into the build. Setting up makefile to capture.

But.... Still need to way to roll the version:
1. find current version
2. **update** version. vMajor.minor.patch
    * patch updated by merge-to-main process - "release patch" -> "release"
    * minor version updated by "version" process ("release minor", "release bump")
    * Major version updated by hand/param ("release v3.0.0")
3. store version. how?
4. run release process (below)

make publish/release:
```
go mod tidy
go test ./...

# version = "v0.1.0"
# name = "rodeo"
# owner = "acmerocket"
# project = "github.com/${owner}/${name}"
# project_vers = "${project}@${version}"
git commit -m "$name: releasing version $version, build $buildnum, on $date"
git tag $version
git push origin $version
GOPROXY=proxy.golang.org go list -m "${project}@${version}"

```

generating a build hash:
```
HASH=$(shell git describe --always)
LDFLAGS=-ldflags "-s -w -X main.Version=${HASH}"
```

looking over github project:
- https://github.com/gomatic/go-vbuild
- https://github.com/thatInfrastructureGuy/git-describe -> https://github.com/golang/go/issues/50603
- examples/tests:
    - https://github.com/kgolding/xldflags - really just an example
    - https://github.com/spellgen/ldflagsx
    - https://github.com/JetSetIlly/ldflags_X_test

do another gh search fo `gomatic/go-vbuild` - Only used by gomatic. plenty of examples.

Noting: https://github.com/golang/go/issues/50603 describes a new go versioning feature, but I still haven't figured out how it works. Still reading that long doc.

> `runtime/debug.BuildInfo.Main.Version`

How is that runtime info set?

> Whether we load version control information is controlled by the `-buildvcs` flag. In the default `-buildvcs=auto` mode we only load verscion control info when we're producing an artifact (`go build`, `go install`) or explicitly providing information to the user (`go list`). See also [#52338 (comment)](https://github.com/golang/go/issues/52338#issuecomment-1104144397) for why we don't get version control information in go run.
> You could set `-buildvcs=true` to force version control information to be populated.

parsing server strings: https://pkg.go.dev/golang.org/x/mod/semver

best guess, so far: https://pkg.go.dev/github.com/Masterminds/semver/v3
more semver management/compare: https://github.com/Masterminds/semver Any way to increment? Provides increment functions for major, minor, build: https://github.com/Masterminds/semver/blob/1558ca3488226e3490894a145e831ad58a5ff958/version.go#L326

older, with inc functions: https://github.com/blang/semver/blob/4487282d78122a245e413d7515e7c516b70c33fd/v4/range.go#L280

Who uses Masterminds/semver? https://github.com/search?q=Masterminds%2Fsemver%2Fv3++language%3AGo&type=code tells me there are over 11k projects using v3 of this library, so that seems valid.

releasing go with git actions: https://github.com/go-semantic-release/semantic-release?tab=readme-ov-file#releasing-a-go-application-with-github-actions

looks like this provide a "release major" and "release minor" commands: https://github.com/go-task/task/ "task yaml" can't find docs for release.

https://github.com/goreleaser/goreleaser, https://goreleaser.com/quick-start/

goreleaser has `incpatch "v1.2.4"`, `incminor== "v1.2.4"` and `incmajor "v1.2.4"`

https://github.com/search?q=path%3A**%2F.goreleaser.yaml+incminor&type=code
`path:**/.goreleaser.yaml incminor` - gh search for all files named `.goreleaser.yaml` containing `incminor`

Total 31. But all for some sort of snapshot: https://goreleaser.com/customization/snapshots


working in a clean version of `rodeo`:
```
sudo snap install --classic goreleaser

brew install goreleaser
goreleaser init
goreleaser release --snapshot --clean
goreleaser check
goreleaser healthcheck

goreleaser build
```

mac: `brew install goreleaser`

Still a lot of overhead, and nothing clear for my need.

I guess: `VERSION=v0.4.0; make release`

### golang makefiles

```
all: build test

build:
        go build -v ./...

test:
        go test -v ./...

install:
        go install -v ./...

clean:
        go clean -v ./...
        rm -f .cover.html .cover.out

cover:
        go test -v -coverprofile .cover.out ./...
        go tool cover -html .cover.out -o .cover.html
        #open .cover.html
```

see:
- ./github/rwtxt/Makefile
- ./github/micro/Makefile


notes:
```
release: test
	# version = "v0.1.0"
	# name = "rodeo"
	# owner = "acmerocket"
	# project = "github.com/${owner}/${name}"
	# project_vers = "${project}@${version}"
	# HASH=$(shell git describe --always)
	# LDFLAGS=-ldflags "-s -w -X main.Version=${HASH}"
	# FIXME increment build number
	name := "rodeo"
	buildhash := $(shell git describe --always)
	release_version := $(VERSION)-$(buildhash)
	git commit -m "$(name): releasing version $(release_version) on $(date)"
	git tag "$(release_version)"
	git push origin "$(release_version)"
	GOPROXY=proxy.golang.org go list -m "${project}@${release_version}"
	# requires export GITHUB_TOKEN="YOUR_GH_TOKEN", see https://github.com/settings/tokens/new?scopes=repo,write:packages
	# meh... still no bump major/minor
	# goreleaser release
```

to release:
```
VERSION=v0.4.0 make release
```

next step:

### type filting
Add simple param handling, first is an array of filters string: type.
```
rodeo post like
```
will ignore evernthing put post and like types.

golang "standard" arg handling, using https://pkg.go.dev/flag

https://gobyexample.com/command-line-flags:
```
wordPtr := flag.String("word", "foo", "a string")
numbPtr := flag.Int("numb", 42, "an int")
forkPtr := flag.Bool("fork", false, "a bool")

var svar string
flag.StringVar(&svar, "svar", "bar", "a string var")

flag.Parse()

fmt.Println("word:", *wordPtr)
fmt.Println("numb:", *numbPtr)
fmt.Println("fork:", *forkPtr)
fmt.Println("svar:", svar)
fmt.Println("tail:", flag.Args())
```

Added simple logging, with slog.

The post like stuff seems to work, but I'm only seeing one type.

So.... add type summary in a final table at exit.
so... global hashmap of type and counts of that type. int

Collecting the data, but can't print summary as the only way is to catch ctrl-C.

https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i

```
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time" // or "runtime"
)

func cleanup() {
    fmt.Println("cleanup")
}

func main() {
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        cleanup()
        os.Exit(1)
    }()

    for {
        fmt.Println("sleeping...")
        time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
    }
}
```

interrupt handling is working. type counts printing

```
{{ range $key, $value := . }}
* **{{ $key }}**: {{ $value }}
{{ end }}
```

time to cut a new release!

```
VERSION=0.5.0 make release
```

Adding a summary template.


## 2025-04-12
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

todo
- [ ] wrap width (use terminal)
- [x] non-record records

Everything committed and uploaded.
