.PHONY: default clean request
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./modes clean
	make -C ./anonalg clean
	make -C ./schemes clean
request:
	./_request/raw/request.sh
