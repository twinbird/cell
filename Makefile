cell: *.go parser.y
	go generate
	go build

.PHONY: test
test: cell *_test.go
	rm -f *.xlsx
	go test
	./test.sh

.PHONY: wintest
wintest: cell *_test.go
	del /Q *.xlsx
	go test
	test.bat

.PHONY: clean
clean:
	rm -f *.xlsx
	rm -f cell
	rm -f y.go
	rm -f y.output
	rm -rf bin

.PHONY: winclean
winclean:
	del /Q *.xlsx
	del /Q cell
	del /Q y.go
	del /Q y.output
	rd /s /q bin
