GC=go build
BINPATH=../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/hla-tcp ./cmd/hla/hla-tcp
run:
	./$(BINPATH)/hla-tcp
clean:
	rm -f hla-tcp.yml $(BINPATH)/hla-tcp
