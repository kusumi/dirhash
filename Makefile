bin1:
	go build -o dirhash -tags squash1
bin2:
	go build -o dirhash -tags squash2
fmt:
	go fmt
test1:
	go test -v -tags squash1
test2:
	go test -v -tags squash2
clean:
	go clean
