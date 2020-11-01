suiteservego:
	go build -o suiteservego github.com/suiteserve/go-runner/cmd/suiteservego

.PHONY: test
test:
	go test -json . | go run github.com/suiteserve/go-runner/cmd/suiteservego \
		-reprint https://localhost:8080

.PHONY: clean
clean:
	rm -f suiteservego
