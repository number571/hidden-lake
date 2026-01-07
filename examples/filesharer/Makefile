.PHONY: default clean request-list request-info request-load
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./simple clean
request-list:
	./_request/list.sh
request-info:
	./_request/info.sh
request-load:
	./_request/load.sh
