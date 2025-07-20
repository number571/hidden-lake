ROOPATH=./../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node1_hlf/hlk.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node2.key \
		./node2_hlf/hlk.key
	mkdir ./node2_hlf/hls-filesharer.stg | true
	cp $(ROOPATH)/examples/filesharer/_storage/* ./node2_hlf/hls-filesharer.stg
run:
	docker-compose build
	docker-compose up
clean:
	rm -rf ./node2_hlf/hls-filesharer.stg
	rm -f ./node1_hlf/hlk.key ./node2_hlf/hlk.key
stop:
	docker-compose down 
	docker-compose rm -f
