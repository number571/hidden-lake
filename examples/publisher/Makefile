PUBLIC ?= false
ARG ?= ""

.PHONY: default clean 
	request-test-listen
	request-test-size request-test-size-public request-test-size-remote request-test-size-remote-public
	request-test-push request-test-push-public request-test-push-remote
	request-test-load-by-uid request-test-load-by-uid-public request-test-load-by-uid-remote request-test-load-by-uid-remote-public
	request-test-load-by-nid request-test-load-by-nid-public request-test-load-by-nid-remote request-test-load-by-nid-remote-public
default: clean 
clean:
	make -C ./prod_test clean
	make -C ./routing clean
	make -C ./simple clean

request-test-listen:
	./_request/listen.sh

request-test-size:
	./_request/size.sh
request-test-size-public:
	./_request/size_public.sh
request-test-size-remote:
	./_request/size_remote.sh
request-test-size-remote-public:
	./_request/size_remote_public.sh

request-test-push: 
	./_request/push.sh $(ARG)
request-test-push-public: 
	./_request/push_public.sh $(ARG)
request-test-push-remote: 
	./_request/push_remote.sh $(ARG)

request-test-load-by-uid:
	./_request/load_by_uid.sh $(ARG)
request-test-load-by-uid-public:
	./_request/load_by_uid_public.sh $(ARG)
request-test-load-by-uid-remote:
	./_request/load_by_uid_remote.sh $(ARG)
request-test-load-by-uid-remote-public:
	./_request/load_by_uid_remote_public.sh $(ARG)

request-test-load-by-nid:
	./_request/load_by_nid.sh $(ARG)
request-test-load-by-nid-public:
	./_request/load_by_nid_public.sh $(ARG)
request-test-load-by-nid-remote:
	./_request/load_by_nid_remote.sh $(ARG)
request-test-load-by-nid-remote-public:
	./_request/load_by_nid_remote_public.sh $(ARG)
