package client

import (
	http "github.com/valyala/fasthttp"
)

type Client struct {
	Jar    *Jar
	config *Config
	client *http.Client
}

type Config struct {
	BaseURL string
	Pre     func(*http.Request)
	Post    func(*http.Response)
}

func New() *Client {
	return &Client{
		client: &http.Client{},
		config: &Config{},
	}
}

func NewWithConfig(config *Config) *Client {
	if config == nil {
		config = &Config{}
	}
	return &Client{
		client: &http.Client{},
		config: config,
	}
}

func (c *Client) NewRequest(method string, url string, body []byte, options ...*Option) *http.Request {
	req := http.AcquireRequest()

	for _, opt := range options {
		opt.Transform(req)
		ReleaseOption(opt)
	}

	if c.Jar != nil {
		c.Jar.mu.Lock()

		for _, c := range c.Jar.cookies {
			req.Header.SetCookieBytesKV(c.Key(), c.Value())
		}

		c.Jar.mu.Unlock()
	}

	if body != nil {
		req.SetBody(body)
	}
	req.SetRequestURI(c.buildURL(url))
	req.Header.SetMethod(method)

	return req
}

func (c *Client) Do(req *http.Request) *http.Response {
	if c.config.Pre != nil {
		c.config.Pre(req)
	}

	resp := http.AcquireResponse()
	defer http.ReleaseRequest(req)

	c.client.Do(req, resp)

	if c.Jar != nil {
		c.Jar.mu.Lock()

		resp.Header.VisitAllCookie(func(key, value []byte) {
			cookie := http.AcquireCookie()
			cookie.ParseBytes(value)

			c.Jar.cookies[string(cookie.Key())] = cookie

		})

		c.Jar.mu.Unlock()
	}

	if c.config.Post != nil {
		c.config.Post(resp)
	}

	return resp
}

func (c *Client) Get(url string) *http.Response {
	req := c.NewRequest(http.MethodGet, url, nil)

	return c.Do(req)
}

func (c *Client) buildURL(endpoint string) string {
	if c.config == nil {
		return endpoint
	}

	return c.config.BaseURL + endpoint
}
