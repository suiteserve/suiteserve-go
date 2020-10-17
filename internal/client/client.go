package client

import (
	"bytes"
	"encoding/json"
	"github.com/suiteserve/go-runner/internal/sstesting"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

const contentType = "application/json"

type Client struct {
	http.Client

	id       string
	casesIdx int64
	url      string
}

func Open(url, name string, tags []string) (*Client, error) {
	c := &Client{
		Client: http.Client{},
		url:    url,
	}
	if err := c.startSuite(name, tags); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) OnEvent(e *sstesting.TestEvent) error {
	if e.Action == sstesting.TestEventActionRun {
		return c.startCase(e.Package + "/" + e.Test)
	}
	return nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) startSuite(name string, tags []string) (err error) {
	url := c.url + "/v1/suites"
	resp, err := c.Post(url, contentType, mustMarshalJSON(map[string]interface{}{
		"name":       name,
		"tags":       tags,
		"status":     "started",
		"started_at": timestamp(),
	}))
	if err != nil {
		return err
	}
	defer drainAndClose(resp.Body, &err)
	return json.NewDecoder(resp.Body).Decode(&c.id)
}

func (c *Client) startCase(name string) (err error) {
	now := timestamp()
	url := c.url + "/v1/cases"
	resp, err := c.Post(url, contentType, mustMarshalJSON(map[string]interface{}{
		"suite_id":   c.id,
		"name":       name,
		"idx":        atomic.AddInt64(&c.casesIdx, 1),
		"status":     "started",
		"created_at": now,
		"started_at": now,
	}))
	if err != nil {
		return err
	}
	defer drainAndClose(resp.Body, &err)
	var id string
	// TODO
	return json.NewDecoder(resp.Body).Decode(&id)
}

func mustMarshalJSON(v interface{}) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func drainAndClose(r io.ReadCloser, err *error) {
	if _, cerr := io.Copy(ioutil.Discard, r); cerr != nil && *err == nil {
		*err = cerr
	}
	if cerr := r.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

func timestamp() int64 {
	now := time.Now()
	return now.Unix()*1e3 + int64(now.Nanosecond())/1e6
}