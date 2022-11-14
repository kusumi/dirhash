bin1:
	go build -tags squash1
bin2:
	go build -tags squash2
fmt:
	go fmt
test1:
	go test -v -tags squash1
test2:
	go test -v -tags squash2
clean:
	go clean

lint1:
	golangci-lint run --build-tags squash1
lint2:
	golangci-lint run --build-tags squash2

xxx1:	fmt lint1 test1
xxx2:	fmt lint2 test2
