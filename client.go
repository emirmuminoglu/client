package client

import (
	http "github.com/valyala/fasthttp"
)

type Client struct {
	Jar *Jar	
}
	
func NewRequest(method string, url string, body []byte) *http.Request{
	req := http.AcquireRequest()

	req.SetBody(body)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	return req
}

func (c *Client) Do(req *http.Request) *http.Response{
	resp := http.AcquireResponse()
	defer http.ReleaseRequest(req)

	if c.Jar != nil {
		c.Jar.mu.Lock()

		for _, c := range c.Jar.cookies{
			req.Header.SetCookieBytesKV(c.Key(), c.Value())
		}

		c.Jar.mu.Unlock()
	}

	http.Do(req, resp)

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
	req := NewRequest("GET", url, nil)

	return c.Do(req)
}