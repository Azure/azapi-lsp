fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w