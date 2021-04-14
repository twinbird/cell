cell: *.go
	go build

test: *.go
	go test

clean:
	rm -f *.xlsx
