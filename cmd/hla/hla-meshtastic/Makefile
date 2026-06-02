GC=go build
BINPATH=../../../bin
.PHONY: default build run clean
default: build run
build:
	$(GC) -o $(BINPATH)/hla-meshtastic ./cmd/hla/hla-meshtastic
run:
	./$(BINPATH)/hla-meshtastic
clean:
	rm -f hla-meshtastic.yml $(BINPATH)/hla-meshtastic
