N=1
GC=go build
ROOPATH=..
BINPATH=$(ROOPATH)/bin

.PHONY: default build move_apps remove-std clean 

default: build move_apps remove-std
build: 
	for app in \
		hlk 				hlc \
		hla/hla-tcp 		hla/hla-http \
		hls/hls-messenger 	hls/hls-messenger/client/hls-messenger-cli \
		hls/hls-filesharer 	hls/hls-filesharer/client/hls-filesharer-cli \
		hls/hls-remoter 	hls/hls-remoter/client/hls-remoter-cli \
		hls/hls-pinger 		hls/hls-pinger/client/hls-pinger-cli; \
	do \
		$(GC) -o $(BINPATH)/$${app} ./$${app}; \
		for arch in amd64 arm64; \
		do \
			for platform in linux windows darwin; \
			do \
				echo "build $${app}_$${arch}_$${platform}"; \
				if [[ $$platform == "windows" ]] \
				then \
					CGO_ENABLED=0 GOOS=$${platform} GOARCH=$${arch} go build -o $(BINPATH)/$${app}_$${arch}_$${platform}.exe ./$${app}; \
				else \
					CGO_ENABLED=0 GOOS=$${platform} GOARCH=$${arch} go build -o $(BINPATH)/$${app}_$${arch}_$${platform} ./$${app}; \
				fi; \
			done; \
		done; \
	done;

move_apps:
	mv $(BINPATH)/hla/* $(BINPATH)/
	mv $(BINPATH)/hls/* $(BINPATH)/
	rm -r $(BINPATH)/hla $(BINPATH)/hls 
	
remove-std:
	make -C $(BINPATH) remove-std

clean:
	make -C $(BINPATH) clean
