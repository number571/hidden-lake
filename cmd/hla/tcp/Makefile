GC=go build
BINPATH=../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/tcp ./cmd/tcp
run:
	./$(BINPATH)/tcp
clean:
	rm -f hla_tcp.yml $(BINPATH)/tcp
