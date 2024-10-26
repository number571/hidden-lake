ROOPATH=./../../../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/anonymity/_keys/priv_node1.key \
		./node1_hlm/hls.key
	cp \
		$(ROOPATH)/examples/anonymity/_keys/priv_node2.key \
		./node2_hlm/hls.key
run:
	docker-compose build
	docker-compose up
clean:
	rm -f ./node1_hlm/hls.key ./node2_hlm/hls.key
stop:
	docker-compose down 
	docker-compose rm -f
