cell: *.go parser.y
	go generate
	go build

test: cell *_test.go
	rm -f *.xlsx
	go test
	./test.sh

.PHONY: install-dev-tools
install-dev-tools:
	go get golang.org/x/tools/cmd/stringer
	go install golang.org/x/tools/cmd/goyacc

.PHONY: clean
clean:
	rm -f *.xlsx
	rm -f cell
	rm -f y.go
	rm -f y.output
