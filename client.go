package client

import (
	http "github.com/valyala/fasthttp"
)

type Client struct {
	Jar *Jar
	config *Config
	client *http.Client	
}
	
type Config struct {
	BaseURL string

}

func New() *Client {
	return &Client{
		client: &http.Client{},		
	}
}

func NewWithConfig(config *Config) *Client {
	return &Client{
		client: &http.Client{},
		config: config,
	}
}

func (c *Client) NewRequest(method string, url string, body []byte) *http.Request{
	req := http.AcquireRequest()

	if c.Jar != nil {
		c.Jar.mu.Lock()

		for _, c := range c.Jar.cookies{
			req.Header.SetCookieBytesKV(c.Key(), c.Value())
		}

		c.Jar.mu.Unlock()
	}

	req.SetBody(body)
	req.SetRequestURI(c.buildURL(url))
	req.Header.SetMethod(method)

	return req
}

func (c *Client) Do(req *http.Request) *http.Response{
	resp := http.AcquireResponse()
	defer http.ReleaseRequest(req)

	c.client.Do(req, resp)

	if c.Jar != nil {
		c.Jar.mu.Lock()

		resp.Header.VisitAllCookie(func(key, value []byte){
			cookie := http.AcquireCookie()
			cookie.ParseBytes(value)
			
			c.Jar.cookies[string(cookie.Key())] = cookie

		})

		c.Jar.mu.Unlock()
	}

	return resp
}

func (c *Client) Get(url string) *http.Response {
	req := c.NewRequest("GET", url, nil)

	return c.Do(req)
}

func (c *Client) buildURL(endpoint string) string {
	if c.config == nil {
		return endpoint
	}

	return c.config.BaseURL + endpoint 

}