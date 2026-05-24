ROOPATH=./../../../..
.PHONY: default run clean stop
default: run 
run:
	docker-compose build
	docker-compose up
stop:
	docker-compose down 
	docker-compose rm -f
