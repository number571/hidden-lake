N=1
GC=go build
ROOPATH=..
BINPATH=$(ROOPATH)/bin

.PHONY: default build move_apps remove-std clean 

default: build move_apps remove-std
build: 
	for app in \
		hlk \
		hlc \
		hla/hla-tcp hla/hla-http \
		hls/hls-messenger hls/hls-filesharer hls/hls-remoter hls/hls-pinger \
		hlp/hlp-chat; \
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
	mv $(BINPATH)/hlp/* $(BINPATH)/
	rm -r $(BINPATH)/hla $(BINPATH)/hls $(BINPATH)/hlp
	
remove-std:
	make -C $(BINPATH) remove-std

clean:
	make -C $(BINPATH) clean
