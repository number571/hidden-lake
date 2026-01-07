.PHONY: default clean request-ping
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./simple clean
request-ping:
	./_request/ping.sh
