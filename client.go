/*
Mailgun client in Go.
*/
package mailgun

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient  *http.Client
	key         string
  apiVersion  int
  apiEndpoint string
}

func New(key string) *Client {
	return &Client{httpClient: &http.Client{}, key: key, apiVersion: 2, apiEndpoint: "api.mailgun.net"}
}

func (c *Client) SetEndpoint(str string) {
  c.apiEndpoint = str
}

// make an api request
func (c *Client) api(method string, path string, fields url.Values) (body []byte, err error) {
	var req *http.Request
	url := fmt.Sprintf("https://%s/v%d%s", c.apiEndpoint, c.apiVersion, path)

	if method == "POST" && fields != nil {
		req, err = http.NewRequest(method, url, strings.NewReader(fields.Encode()))

    if err != nil {
      panic(err)
    }

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	} else {
		if fields != nil {
			url += "?" + fields.Encode()
		}
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return
	}
	req.SetBasicAuth("api", c.key)
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	if rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		err = fmt.Errorf("mailgun error: %d %s", rsp.StatusCode, body)
	}
	return
}
