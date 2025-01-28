ROOPATH=./../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node1_hln/hls.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node2.key \
		./node2_hln/hls.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node3.key \
		./node3_hln/hls.key
run:
	docker-compose build
	docker-compose up
clean:
	rm -f ./node1_hln/hls.key ./node2_hln/hls.key ./node3_hln/hls.key
stop:
	docker-compose down 
	docker-compose rm -f
