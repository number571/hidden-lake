ROOPATH=./../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node1/hlk.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node2.key \
		./node2/hlk.key
run:
	docker-compose build
	docker-compose up -d
	docker-compose cp ./hla_tcp hla_tcp:/hidden-lake
	docker-compose cp ./node1 node1_hlc:/hidden-lake
	docker-compose cp ./node2 node2_hlc:/hidden-lake
	docker-compose restart
	docker-compose up
clean:
	rm -f ./node1/hlk.db ./node2/hlk.db
	rm -f ./node1/hlk.key ./node2/hlk.key
	rm -f ./node1/hls-messenger.db ./node2/hls-messenger.db
	rm -rf ./node1/hls-filesharer.stg ./node2/hls-filesharer.stg
stop:
	docker-compose down 
	docker-compose rm -f
