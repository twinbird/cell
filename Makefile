cell: *.go parser.y
	go generate
	go build

test: cell *_test.go
	rm -f *.xlsx
	go test
	./test.sh

clean:
	rm -f *.xlsx
	rm -f cell
	rm -f y.go
	rm -f y.output
