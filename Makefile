export majorVersion=1
export minorVersion=0
export patchVersion=$(shell git log --format='%h' | wc -l)
export ver=$(majorVersion).$(minorVersion).$(patchVersion)

start:
	go run example/main.go

fix:
	go fix ./...

# https://go.dev/blog/govulncheck
# install it by `go install golang.org/x/vuln/cmd/govulncheck@latest`
vuln:
	which govulncheck
	govulncheck ./...

lint:
	gofmt -w=true -s=true -l=true ./
	golint ./...
	go vet ./...
	staticcheck ./...

tools:
	which go
	which golint # go install golang.org/x/lint/golint@latest
	which govulncheck # go install honnef.co/go/tools/cmd/staticcheck@latest

deps:
	go mod verify
	go mod download
	go mod tidy

tag:
	git tag "v$(ver)"


include make/*.mk
