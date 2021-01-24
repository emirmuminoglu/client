package client

import (
	"encoding/base64"
	http "github.com/valyala/fasthttp"
	"sync"
)

var (
	optionPool sync.Pool
	zeroOption = &Option{}
)

func AcquireOption() *Option {
	v := optionPool.Get()
	if v == nil {
		return new(Option)
	}

	return v.(*Option)
}

func ReleaseOption(opt *Option) {
	*opt = *zeroOption

	optionPool.Put(opt)
}

type Option struct {
	Transform func(*http.Request)
}

func BasicAuth(username, password string) *Option {
	return &Option{
		Transform: func(req *http.Request) {
			toEncode := []byte(username + ":" + password)

			req.Header.Add(http.HeaderAuthorization, basic+base64.StdEncoding.EncodeToString(toEncode))
		},
	}
}

func JSON() *Option {
	return &Option{
		Transform: func(req *http.Request) {
			req.Header.Add(http.HeaderContentType, jsonContentType)
		},
	}
}
