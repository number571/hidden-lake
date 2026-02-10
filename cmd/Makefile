N=1
GC=go build
ROOPATH=..
BINPATH=$(ROOPATH)/bin

.PHONY: default build build-cli move_apps remove-std clean 
default: build build-cli move_apps remove-std

build: 
	for app in \
		hlc \
		hlk	\
		hla/hla-tcp	\
		hla/hla-http \
		hls/hls-messenger \
		hls/hls-filesharer \
		hls/hls-pinger; \
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

build-cli: 
	for app in \
		hlk/client/hlk-cli \
		hls/hls-messenger/client/hls-messenger-cli \
		hls/hls-filesharer/client/hls-filesharer-cli \
		hls/hls-pinger/client/hls-pinger-cli; \
	do \
		$(GC) -o $(BINPATH)/cli/$$(basename "$${app}") ./$${app}; \
		for arch in amd64 arm64; \
		do \
			for platform in linux windows darwin; \
			do \
				echo "build $${app}_$${arch}_$${platform}"; \
				if [[ $$platform == "windows" ]] \
				then \
					CGO_ENABLED=0 GOOS=$${platform} GOARCH=$${arch} go build -o $(BINPATH)/cli/$$(basename "$${app}")_$${arch}_$${platform}.exe ./$${app}; \
				else \
					CGO_ENABLED=0 GOOS=$${platform} GOARCH=$${arch} go build -o $(BINPATH)/cli/$$(basename "$${app}")_$${arch}_$${platform} ./$${app}; \
				fi; \
			done; \
		done; \
	done;

move_apps:
	mv $(BINPATH)/hla/* $(BINPATH)/
	mv $(BINPATH)/hls/* $(BINPATH)/
	mv $(BINPATH)/cli/* $(BINPATH)/
	rm -r $(BINPATH)/hla $(BINPATH)/hls 

remove-std:
	make -C $(BINPATH) remove-std

clean:
	make -C $(BINPATH) clean
