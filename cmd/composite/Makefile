GC=go build
ROOPATH=../..
BINPATH=$(ROOPATH)/bin

.PHONY: default build run test clean

default: build run
build:
	$(GC) -o $(BINPATH)/hlc ./cmd/hlc

	for arch in amd64 arm64; \
	do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_linux ./cmd/hlc; \
		CGO_ENABLED=0 GOOS=windows GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_windows.exe ./cmd/hlc; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=$${arch} go build -o $(BINPATH)/hlc_$${arch}_darwin ./cmd/hlc; \
	done;
run:
	./$(BINPATH)/hlc
test:
	for i in {1..$(N)}; do go clean -testcache; echo $$i; go test ./...; done
clean:
	rm -f $(BINPATH)/hlc $(BINPATH)/hlc_*
	rm -rf hl?.* priv.key
