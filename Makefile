VERSION ?= dev
NAME = rodeo
#$(basename $(dir $(abspath "$PWD")))
ORG = acmerocket
HASH = $(shell git describe --always)
RELEASE = $(VERSION)-$(HASH)
PROJECT = github.com/$(ORG)/$(NAME)
GOVARS = -X $(PROJECT)/main.Version=$(VERSION) -X $(PROJECT)/main.VersionHash=$(HASH)
DEBUGVAR = -X -X $(PROJECT)/main.Debug=ON

all: test

tidy:
	go mod tidy
	go fmt ./...

build: tidy
	go build ./...

test: tidy
	go test ./...

install: test
	go install -v ./...

clean:
	go clean ./...
	rm -fr .cover.html .cover.out dist/

cover:
	go test -v -coverprofile .cover.out ./...
	go tool cover -html .cover.out -o .cover.html
	#open .cover.html

release: test
	git commit -m "$(NAME): releasing version $(RELEASE) on $(shell date)"
	git tag "$(RELEASE)"
	git push origin "$(RELEASE)"
	GOPROXY=proxy.golang.org go list -m "$(PROJECT)@$(RELEASE)"
