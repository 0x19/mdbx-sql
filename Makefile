
test:
	go clean -testcache && go test ./... -v -cover

bench:
	go clean -testcache && go test -bench=./... -benchmem