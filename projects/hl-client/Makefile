.PHONY: default install-deps build
default: build 
install-deps:
	go install github.com/fyne-io/fyne-cross@latest
build:
	mkdir -p bin
	for arch in amd64 arm64; \
	do \
		for platform in linux windows; \
		do \
			echo "build $${arch}_$${platform}"; \
			fyne-cross $${platform} -arch=$${arch} --app-id hidden.lake.client --icon images/icons/icon.png; \
			if [[ $$platform == "windows" ]] \
			then \
				cp fyne-cross/bin/$${platform}-$${arch}/hidden-lake-client.exe ./bin/hidden-lake-client_$${arch}_$${platform}.exe; \
			else \
				cp fyne-cross/bin/$${platform}-$${arch}/hidden-lake-client ./bin/hidden-lake-client_$${arch}_$${platform}; \
			fi; \
		done; \
	done;
