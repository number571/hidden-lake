N=1
L=1

.PHONY: default run clean
default: run 
run: 
	for i in {1..$(N)}; \
		do echo "http $$i"; \
		go run . $(L); \
	done;
clean: 
	ls | grep -xv "Makefile" | grep -xv "main.go" | xargs rm | true
