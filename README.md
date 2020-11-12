# SuiteServe for Go

SuiteServe client for the Go programming language.

## Quickstart

Ensure that `$GOPATH/bin` is in your path and then install with:
```bash
$ go get github.com/suiteserve/suiteserve-go/cmd/suiteservego
```

Start the SuiteServe server and note its URL. If you've checked out [suiteserve/suiteserve](https://github.com/suiteserve/suiteserve) and are running it locally, this will be `https://localhost:8080` by default.

From within a Go project of your choice, run its tests with the `-json` flag and pipe the output to `suiteservego <url>`. For example:
```bash
$ go test -count=1 -json ./... | suiteservego https://localhost:8080
```

The `-count=1` option is the idiomatic way to disable test caching explicitly, often useful for end-to-end testing.
