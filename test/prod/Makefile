.PHONY: default test-tcp test-http test-https clean
default: test-tcp test-http test-https

test-tcp: 
	make -C tcp
	echo & echo
test-http: 
	make -C http
	echo & echo
test-https: 
	make -C https
	echo & echo
clean: 
	make -C tcp clean
	make -C http clean
	make -C https clean
