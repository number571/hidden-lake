N=1

.PHONY: default test clean
default: test 

test: 
	for i in {1..$(N)}; \
		do echo $$i; \
		go run ./tcp; \
	done;
clean: 
	ls | grep -xv ".gitignore" | grep -xv "Makefile" | grep -xv "tcp" | xargs rm | true
