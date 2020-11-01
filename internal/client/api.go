package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type jsonObj map[string]interface{}

type SuiteResult string

const (
	SuiteResultPassed SuiteResult = "passed"
	SuiteResultFailed SuiteResult = "failed"
)

type CaseResult string

const (
	CaseResultPassed  CaseResult = "passed"
	CaseResultFailed  CaseResult = "failed"
	CaseResultSkipped CaseResult = "skipped"
	CaseResultAborted CaseResult = "aborted"
	CaseResultErrored CaseResult = "errored"
)

type apiUrl string

const (
	suitesUrl apiUrl = "/v1/suites"
	casesUrl  apiUrl = "/v1/cases"
	logsUrl   apiUrl = "/v1/logs"
)

func (c *Client) buildUrl(base apiUrl, segment string,
	queries ...string) apiUrl {
	if len(queries)%2 != 0 {
		panic("expected query pairs")
	}
	var query string
	if len(queries) >= 2 {
		vals := url.Values{}
		for i := 1; i < len(queries); i += 2 {
			vals.Set(queries[i-1], queries[i])
		}
		query = "?" + vals.Encode()
	}
	return apiUrl(fmt.Sprintf("%s/%s%s", base, url.PathEscape(segment), query))
}

func (c *Client) suiteUrl(id string) apiUrl {
	return c.buildUrl(suitesUrl, id)
}

func (c *Client) finishSuiteUrl(id string) apiUrl {
	return c.buildUrl(suitesUrl, id, "finish", "true")
}

func (c *Client) caseUrl(id string) apiUrl {
	return c.buildUrl(casesUrl, id)
}

func (c *Client) logUrl(id string) apiUrl {
	return c.buildUrl(logsUrl, id)
}

func (c *Client) startSuite(name string, tags []string) error {
	return c.sendReq(http.MethodPost, suitesUrl, jsonObj{
		"name":      name,
		"tags":      tags,
		"status":    "started",
		"startedAt": timestamp(),
	}, &c.id)
}

func (c *Client) finishSuite(res SuiteResult) error {
	return c.sendReq(http.MethodPatch, c.finishSuiteUrl(c.id), jsonObj{
		"result": res,
		"at":     timestamp(),
	})
}

func (c *Client) startCase(name string) error {
	var id string
	now := timestamp()
	if err := c.sendReq(http.MethodPost, casesUrl, jsonObj{
		"suiteId":   c.id,
		"name":      name,
		"idx":       c.incIdx(),
		"status":    "started",
		"createdAt": now,
		"startedAt": now,
	}, &id); err != nil {
		return err
	}
	c.cases[name] = id
	return nil
}

func (c *Client) finishCase(name string, res CaseResult) error {
	u := c.buildUrl(casesUrl, c.cases[name], "finish", "true")
	if err := c.sendReq(http.MethodPatch, u, jsonObj{
		"result": res,
		"at":     timestamp(),
	}); err != nil {
		return err
	}
	delete(c.cases, name)
	return nil
}

func (c *Client) addLog(caseName, line string) (err error) {
	return c.sendReq(http.MethodPost, logsUrl, jsonObj{
		"caseId": c.cases[caseName],
		"idx":    c.incIdx(),
		"line":   line,
	})
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

func (c *Client) sendReq(method string, u apiUrl,
	body interface{}, v ...interface{}) (err error) {
	req, err := http.NewRequest(method, c.baseUrl+string(u),
		mustMarshalJSON(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer drainAndClose(resp.Body, &err)
	if len(v) == 0 {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(v[0])
}

func mustMarshalJSON(v interface{}) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}
