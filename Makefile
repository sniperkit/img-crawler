.PHONY: build
build: spider

.PHONY: spider
spider:
	go build -v -o bin/spider cmd/spider.go

.PHONY: debug
debug:
	go build -gcflags "-N -l" -o bin/spider  cmd/spider.go


.PHONY: clean
clean:
	rm bin/*

fmt:
	goimports -w ./src
	goimports -w ./cmd

format:
	go fmt ./src/...
	go fmt ./cmd/...


install:
	dep ensure

update:
	dep ensure -update
