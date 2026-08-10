.PHONY: default clean request-node1 request-node2
default: clean 
clean:
	make -C ./group clean
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./simple clean
request-node1:
	./_request/node1.sh
request-node2:
	./_request/node2.sh
