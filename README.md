# rodeo
Format JSON based on markdown templates. Accept input as well-formed chunks of JSON to be parsed and applied to a template.

```
go install github.com/acmerocket/rodeo@latest
goat firehose --ops | rodeo
```

Designed for use with [goat](https://github.com/bluesky-social/indigo/tree/main/cmd/goat)

## Usage
```
TBD
```

## Build
```
git clone git@github.com:acmerocket/rodeo.git
cd rodeo
make

# to release
VERSION=v0.4.0 make release
```
