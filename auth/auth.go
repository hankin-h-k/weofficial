package auth

import "github.com/hankin-h-k/weofficial/request"

type Auth struct {
	request *request.Request

	// 组成完整的 URL 地址
	// 默认包含 AccessToken
	combineURI func(url string, req interface{}, withToken bool) (string, error)
}

func NewAuth(request *request.Request, combineURI func(url string, req interface{}, withToken bool) (string, error)) *Auth {
	sm := Auth{
		request:    request,
		combineURI: combineURI,
	}

	return &sm
}
