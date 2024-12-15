GC=go build
BINPATH=../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/hla_tcp ./cmd/hla/hla_tcp
run:
	./$(BINPATH)/hla_tcp
clean:
	rm -f hla_tcp.yml $(BINPATH)/hla_tcp
