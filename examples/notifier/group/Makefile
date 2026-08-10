ROOPATH=./../../..
.PHONY: default copy run clean stop request-node1 request-node2 request-node3
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node1_hlm/hlk.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node2_hlm/hlk.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node3_hlm/hlk.key
run:
	docker-compose build
	docker-compose up
clean:
	rm -f ./node1_hlm/hlk.key ./node2_hlm/hlk.key ./node3_hlm/hlk.key
stop:
	docker-compose down 
	docker-compose rm -f

request-node1:
	./_request/node1.sh
request-node2:
	./_request/node2.sh
request-node3:
	./_request/node3.sh
