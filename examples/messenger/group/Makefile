ROOPATH=./../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node1_hlm/hls.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node2_hlm/hls.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node3_hlm/hls.key
run:
	docker-compose build
	docker-compose up
clean:
	rm -f ./node1_hlm/hls.yml ./node2_hlm/hls.yml ./node3_hlm/hls.yml
stop:
	docker-compose down 
	docker-compose rm -f
