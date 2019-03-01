fmt:
	find ! -path "./vendor/*" -name "*.go" -exec gofmt -s -w {} \;

doc-hunt:
	doc-hunt check -e

run-tests:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

test-all: run-tests doc-hunt
