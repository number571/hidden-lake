N=1
GC=go build
ROOPATH=..
BINPATH=$(ROOPATH)/bin

.PHONY: default build move_hla move_hls remove-std clean 

default: build move_hla move_hls remove-std
build: 
	for app in \
		hlk hlc \
		hls/hls-messenger hls/hls-filesharer hls/hls-remoter hls/hls-pinger \
		hla/hla-tcp hla/hla-http; \
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
move_hla:
	mv $(BINPATH)/hla/* $(BINPATH)/
	rm -r $(BINPATH)/hla
move_hls:
	mv $(BINPATH)/hls/* $(BINPATH)/
	rm -r $(BINPATH)/hls
remove-std:
	make -C $(BINPATH) remove-std

clean:
	make -C $(BINPATH) clean
