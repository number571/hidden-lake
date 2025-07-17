GC=go build
BINPATH=../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/hla-http ./cmd/hla/hla-http
run:
	./$(BINPATH)/hla-http
clean:
	rm -f hla-http.yml $(BINPATH)/hla-http
