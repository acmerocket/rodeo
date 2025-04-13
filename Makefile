all: test

build:
	   go build -v ./...

tidy:
	   go mod tidy

test: tidy
    	go test ./...

install: test
        go install -v ./...

clean:
        go clean ./...
        rm -f .cover.html .cover.out

cover:
        go test -v -coverprofile .cover.out ./...
        go tool cover -html .cover.out -o .cover.html
        #open .cover.html

release: test
        # version = "v0.1.0"
        # name = "rodeo"
        # owner = "acmerocket"
        # project = "github.com/${owner}/${name}"
        # project_vers = "${project}@${version}"
        FIXME increment build number
        FIXME generate build id
        git commit -m "$name: releasing version $version, build $buildnum, on $date"
        git tag $version
        git push origin $version
        GOPROXY=proxy.golang.org go list -m "${project}@${version}"
