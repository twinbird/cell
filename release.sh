#!/bin/bash

PROGNAME="cell"

if [ $# != 1 ]; then
	echo "Usage: $0 [version]"
	exit 1
fi

go generate
GOOS=linux GOARCH=amd64 go build -o ./bin/$1/linux64/$PROGNAME
GOOS=windows GOARCH=386 go build -o ./bin/$1/windows386/$PROGNAME.exe
GOOS=windows GOARCH=amd64 go build -o ./bin/$1/windows64/$PROGNAME.exe
GOOS=darwin GOARCH=amd64 go build -o ./bin/$1/darwin64/$PROGNAME

cd ./bin/$1/
tar cfvz $PROGNAME-$1.linux64.tar.gz ./linux64
zip -r $PROGNAME-$1.windows386.zip ./windows386
zip -r $PROGNAME-$1.windows64.zip ./windows64
tar cfvz $PROGNAME-$1.darwin64.tar.gz ./darwin64
