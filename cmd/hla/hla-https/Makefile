GC=go build
BINPATH=../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/hla-https ./cmd/hla/hla-https
run:
	./$(BINPATH)/hla-https
clean:
	rm -f hla-https.yml $(BINPATH)/hla-https
