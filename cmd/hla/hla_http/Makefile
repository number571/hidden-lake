GC=go build
BINPATH=../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/hla_http ./cmd/hla/hla_http
run:
	./$(BINPATH)/hla_http
clean:
	rm -f hla_http.yml $(BINPATH)/hla_http
