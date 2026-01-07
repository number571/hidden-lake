.PHONY: default clean request request-exec
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./simple clean
request-exec:
	./_request/exec.sh
