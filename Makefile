# vim: set ts=4 sw=4 noexpandtab:

all: buckler

buckler: *.go
	go build -v

test: *.go
	go test -v

clean:
	rm -f buckler

.PHONY: test clean
