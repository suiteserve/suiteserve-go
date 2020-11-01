package client

import (
	"github.com/suiteserve/go-runner/internal/sstesting"
	"net/http"
)

type Client struct {
	baseUrl string
	http.Client

	id    string
	cases map[string]string

	idx int64
}

func (c *Client) incIdx() int64 {
	c.idx++
	return c.idx
}

func Open(url, name string, tags []string) (*Client, error) {
	c := &Client{
		baseUrl: url,
		cases:   make(map[string]string),
	}
	if err := c.startSuite(name, tags); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) OnEvent(e *sstesting.TestEvent) error {
	if e.Test == "" {
		return nil
	}
	caseName := e.Package + "/" + e.Test
	switch e.Action {
	case sstesting.TestEventActionRun:
		return c.startCase(caseName)
	case sstesting.TestEventActionOutput:
		return c.addLog(caseName, e.Output)
	case sstesting.TestEventActionPass:
		return c.finishCase(caseName, CaseResultPassed)
	case sstesting.TestEventActionFail:
		return c.finishCase(caseName, CaseResultFailed)
	default:
		return nil
	}
}

func (c *Client) Close() error {
	return c.finishSuite("passed")
}
