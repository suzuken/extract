deps:
	godep restore

test:
	godep go test -v
