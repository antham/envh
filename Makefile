fmt:
	find ! -path "./vendor/*" -name "*.go" -exec go fmt {} \;

gometalinter:
	gometalinter -D gotype --vendor --deadline=240s --dupl-threshold=200 -e '_string' -j 5 ./...

run-tests:
	go test -v -race -coverprofile=profile.out -covermode=atomic

test-all: gometalinter run-tests
