.PHONY: test 

stub:
	@echo It's a stub

mock:
	go generate -v ./...

test:
	go test tests/*