N=1
GC=go build
ROOPATH=..
BINPATH=$(ROOPATH)/bin

.PHONY: default build move_hla remove-std clean 

default: build move_hla remove-std
build: 
	for app in hls hlc hlm hlf hlr hlp hla/hla_tcp hla/hla_http; \
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
remove-std:
	make -C $(BINPATH) remove-std

clean:
	make -C $(BINPATH) clean
