suiteservego:
	go build -o suiteservego github.com/suiteserve/go-runner/cmd/suiteservego

.PHONY: test
test:
	go test -count 1 ./internal/client/clienttest

.PHONY: clean
clean:
	rm -f suiteservego
