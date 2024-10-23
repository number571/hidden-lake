GC=go build
.PHONY: default run clean 
default: clean run 
run:
	./prog_hle &
clean:
	pkill -15 prog_hle || true
	rm -rf prog_hle
