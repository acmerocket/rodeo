# rodeo
Format JSON based on markdown templates. Accept input as well-formed chunks of JSON to be parsed and applied to a template.

```
go install github.com/acmerocket/rodeo@latest
goat firehose --ops | rodeo post
```

`rodeo` is designed to accept a stream of JSON, and convert it into fancy terminal graphics by mapping JSON fields into
[Markdown](https://www.markdownguide.org/basic-syntax/)
[templates](https://pkg.go.dev/text/template) using
[Glamour](https://github.com/charmbracelet/glamour).

Designed originally for use with
[goat](https://github.com/bluesky-social/indigo/tree/main/cmd/goat)

## Usage

FIXME: Usage here!

#### Examples

Grab the [ATmosphere](https://atproto.com/guides/glossary#atmosphere) firehose and process it using the standard templates:
```
goat firehose --ops | rodeo
```

Grab the firehose and process only specific message types, using the standard templates:
```
goat firehose --ops | rodeo post like
```

`rodeo` supports mutiple paramaters, each allowing the following forms:
- `app.bsky.feed.post`: Match the full message `$type`.
- `post`: Partial match of message type.
- `post=default`: Use a different template, in this case `default`, instead of the one derived from message type.
- `post=path/to/template.md`: Use the specified template for messages that match `post`.


To list supported [[built-in message templates]](./templates/):
```
rodeo --list-templates
```

To use an existing template in for a different message type, for example if there is no tempate for a `new.message.type` that you wish to display with the `list` template:
```
goat firehose --ops | rodeo new.message.type=list
```

To use an external template located on the filesystem (`path/to/template.md`)
```
goat firehose --ops | post=path/to/template.md like=path/to/other.md
```

## Development
```
git clone https://github.com/acmerocket/rodeo
cd rodeo
make test
```

This project uses [`make`](https://www.gnu.org/software/make/) and provides:
|---          |--- |
| **build**   | Build the project. |
| **test**    | Run the test suite. |
| **cover**   | Run the test suite with coverage. |
| **install** | Install the project with `go install`. |
| **release** | Release the project. |
| **clean**   | Clean up, remove all generates files. |

#### Releases
To release a new version, specify the version:
```
VERSION=v0.4.0 make release
```
