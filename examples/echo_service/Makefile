.PHONY: default clean request
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./modes clean
	make -C ./networks clean
request:
	go run ./_request/main.go
