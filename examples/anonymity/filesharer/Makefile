.PHONY: default clean request
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./simple clean
	make -C ./secret_channel clean
request:
	go run ./_request/main.go 
