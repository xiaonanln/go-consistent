all: consistent.so

consistent.so: $(wildcard *.go)
	go build -buildmode=c-shared -o $@
