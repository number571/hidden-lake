.PHONY: default remove-std clean 
default: clean
clean: 
	ls | grep -xv ".gitignore" | grep -xv "Makefile" | xargs rm -r | true
remove-std:
	ls | grep -xv ".gitignore" | grep -xv "Makefile" | grep -xv -E ".*_.*_.*" | xargs rm -r 
