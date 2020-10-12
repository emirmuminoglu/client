package client

import (
	http "github.com/valyala/fasthttp"
	"sync"
)

type Jar struct {
	mu sync.Mutex

	cookies map[string]*http.Cookie
}

func NewJar() *Jar {
	return &Jar{cookies: make(map[string]*http.Cookie)}
}

func (j *Jar) PeekValue(key string) []byte{
	c, ok := j.cookies[key]
	if ok{
		return c.Value()
	}

	return nil
}

func (j *Jar) Peek(key string) *http.Cookie{
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.cookies[key]
}

func (j *Jar) ReleaseCookie(key string){
	j.mu.Lock()
	defer j.mu.Unlock()

	c, ok := j.cookies[key]
	if ok{
		http.ReleaseCookie(c)
		delete(j.cookies, key)
	}

}