N=1

.PHONY: default test-tcp test-http test-https clean
default: test-tcp test-http test-https

test-tcp: 
	for i in {1..$(N)}; \
		do echo $$i; \
		make -C tcp; \
	done;
test-http: 
	for i in {1..$(N)}; \
		do echo $$i; \
		make -C http; \
	done;
test-https: 
	for i in {1..$(N)}; \
		do echo $$i; \
		make -C https; \
	done;
clean: 
	make -C tcp clean
	make -C http clean
	make -C https clean
