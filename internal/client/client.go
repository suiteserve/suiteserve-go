package client

import (
	"github.com/suiteserve/suiteserve-go/internal/sstesting"
	"net/http"
)

type Client struct {
	http.Client

	init    bool
	closed  bool
	baseUrl string
	name    string
	tags    []string
	id      string
	cases   map[string]string
	idx     int64
}

func (c *Client) incIdx() int64 {
	c.idx++
	return c.idx
}

func Open(url, name string, tags []string) *Client {
	return &Client{
		baseUrl: url,
		name:    name,
		tags:    tags,
		cases:   make(map[string]string),
	}
}

func (c *Client) OnEvent(e *sstesting.TestEvent) error {
	if c.closed {
		panic("client closed")
	}
	if !c.init {
		if err := c.startSuite(c.name, c.tags); err != nil {
			return err
		}
		c.init = true
	}
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
	case sstesting.TestEventActionSkip:
		return c.finishCase(caseName, CaseResultSkipped)
	default:
		return nil
	}
}

func (c *Client) Close() error {
	if !c.init {
		return nil
	}
	c.closed = true
	return c.finishSuite("passed")
}
