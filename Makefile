cell: *.go y.go
	go build

y.go: parser.y
	goyacc parser.y

test: cell *_test.go
	go test
	./test.sh

clean:
	rm -f *.xlsx
	rm -f cell
	rm -f y.go
	rm -f y.output
