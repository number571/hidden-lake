ROOPATH=./../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node1_hlm/hlk.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node2.key \
		./node2_hlm/hlk.key
run:
	docker-compose build
	docker-compose up
clean:
	rm -f ./node1_hlm/hlk.key ./node2_hlm/hlk.key
stop:
	docker-compose down 
	docker-compose rm -f
