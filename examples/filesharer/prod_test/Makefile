ROOPATH=./../../..
.PHONY: default copy run clean stop
default: copy run 
copy:
	cp \
		$(ROOPATH)/examples/_keys/priv_node1.key \
		./node1_hlf/hls.key
	cp \
		$(ROOPATH)/examples/_keys/priv_node2.key \
		./node2_hlf/hls.key
	mkdir ./node2_hlf/hlf.stg | true
	cp $(ROOPATH)/examples/filesharer/_storage/* ./node2_hlf/hlf.stg
run:
	docker-compose build
	docker-compose up
clean:
	rm -rf ./node2_hlf/hlf.stg
	rm -f ./node1_hlf/hls.key ./node2_hlf/hls.key
stop:
	docker-compose down 
	docker-compose rm -f
