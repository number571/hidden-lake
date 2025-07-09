ROOPATH=./../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./recv_hlc/hlk.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node2.key \
		./send_hls/hlk.key
run:
	docker-compose build
	docker-compose up
clean:
	rm -f ./recv_hlc/hlk.key ./send_hls/hlk.key
stop:
	docker-compose down 
	docker-compose rm -f
