L=1

.PHONY: default run clean
default: run 
run: 
	go run . $(L)
clean: 
	ls | grep -xv "Makefile" | grep -xv "main.go" | xargs rm | true
