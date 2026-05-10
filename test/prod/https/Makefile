N=1
L=1

.PHONY: default run clean
default: run 
run: 
	for i in {1..$(N)}; \
		do echo "https $$i"; \
		go run . $(L); \
	done;
clean: 
	ls | grep -xv "Makefile" | grep -xv "main.go" | grep -xv "hla-https.cert" | grep -xv "hla-https.key" | grep -xv "hla-https.certs" | grep -xv "users.json" | xargs rm | true
