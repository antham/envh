fmt:
	find ! -path "./vendor/*" -name "*.go" -exec go fmt {} \;

gometalinter:
	gometalinter -D gotype --vendor --deadline=240s --dupl-threshold=200 -e '_string' -j 5 ./...

run-tests:
	go test -v $(shell glide nv)

test-all: gometalinter run-tests

test-package:
	go test -race -cover -coverprofile=/tmp/gommit github.com/antham/gommit/$(pkg)
	go tool cover -html=/tmp/gommit
