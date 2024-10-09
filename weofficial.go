package weofficial

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hankin-h-k/weofficial/auth"
	"github.com/hankin-h-k/weofficial/cache"
	"github.com/hankin-h-k/weofficial/logger"
	"github.com/hankin-h-k/weofficial/request"

	"github.com/mitchellh/mapstructure"
)

const (
	baseURL = "https://api.weixin.qq.com"
)

type Client struct {
	// HTTP请求客户端
	request *request.Request
	// 数据缓存器
	cache cache.Cache
	// 日志记录器
	logger logger.Logger
	// 公众号appid
	appid string
	// 公众号秘钥
	secret string
	// 用户自定义获取access_token的方法
	accessTokenGetter AccessTokenGetter
}

// 用户自定义获取access_token的方法
type AccessTokenGetter func(appid, secret string) (token string, expireIn uint)

// 初始化客户端并用自定义配置替换默认配置
func NewClient(appid, secret string, opts ...func(*Client)) *Client {
	cli := &Client{
		appid:  appid,
		secret: secret,
	}

	// 执行额外的配置函数
	for _, fn := range opts {
		fn(cli)
	}

	if cli.cache == nil {
		cli.cache = cache.NewMemoryCache()
	}

	if cli.request == nil {
		cli.request = request.NewRequest(http.DefaultClient, request.ContentTypeJSON, cli.Logger)
	}

	if cli.logger == nil {
		cli.logger = logger.NewLogger(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Info, true)
	}

	return cli
}

// 自定义 HTTP Client
func WithHttpClient(hc *http.Client) func(*Client) {
	return func(cli *Client) {
		cli.request = request.NewRequest(hc, request.ContentTypeJSON, cli.Logger)
	}
}

// 自定义缓存
func WithCache(cc cache.Cache) func(*Client) {
	return func(cli *Client) {
		cli.cache = cc
	}
}

// 自定义获取access_token的方法
func WithAccessTokenSetter(getter AccessTokenGetter) func(*Client) {
	return func(cli *Client) {
		cli.accessTokenGetter = getter
	}
}

// 自定义日志
func WithLogger(logger logger.Logger) func(*Client) {
	return func(cli *Client) {
		cli.logger = logger
	}
}

// POST 参数
type requestParams map[string]interface{}

// URL 参数
type requestQueries map[string]interface{}

// tokenAPI 获取带 token 的 API 地址
func tokenAPI(api, token string) (string, error) {
	queries := requestQueries{
		"access_token": token,
	}

	return request.EncodeURL(api, queries)
}

// convert bool to int
func bool2int(ok bool) uint8 {

	if ok {
		return 1
	}

	return 0
}

// 获取日志记录器
func (cli *Client) Logger() logger.Logger {
	return cli.logger
}

// 设置日志等级
func (cli *Client) SetLogLevel(lv logger.Level) {
	if cli.logger != nil {
		cli.logger.SetLevel(lv)
	}
}

// 获取公众号全局唯一后台接口调用凭据（access_token）。
// 调调用绝大多数后台接口时都需使用 access_token，开发者需要进行妥善保存，注意缓存。
func (cli *Client) AccessToken() (string, error) {
	key := cli.tokenCacheKey()
	data, ok := cli.cache.Get(key)
	if ok {
		return data.(string), nil
	}

	if cli.accessTokenGetter != nil {
		token, expireIn := cli.accessTokenGetter(cli.appid, cli.secret)
		cli.cache.Set(key, token, time.Duration(expireIn)*time.Second)
		return token, nil
	} else {

		req := auth.GetStableAccessTokenRequest{
			Appid:     cli.appid,
			Secret:    cli.secret,
			GrantType: "client_credential",
		}
		rsp, err := cli.NewAuth().GetStableAccessToken(&req)
		if err != nil {
			return "", err
		}

		if err := rsp.GetResponseError(); err != nil {
			return "", err
		}

		cli.cache.Set(key, rsp.AccessToken, time.Duration(rsp.ExpiresIn)*time.Second)
		return rsp.AccessToken, nil
	}
}

// 获取稳定版接口调用凭据
func (cli *Client) StableAccessToken(forceRefresh bool) (string, error) {

	key := cli.tokenCacheKey()
	data, ok := cli.cache.Get(key)
	if ok {
		return data.(string), nil
	}

	if !forceRefresh && cli.accessTokenGetter != nil {
		token, expireIn := cli.accessTokenGetter(cli.appid, cli.secret)
		cli.cache.Set(key, token, time.Duration(expireIn)*time.Second)
		return token, nil
	} else {

		req := auth.GetStableAccessTokenRequest{
			Appid:        cli.appid,
			Secret:       cli.secret,
			GrantType:    "client_credential",
			ForceRefresh: forceRefresh,
		}
		rsp, err := cli.NewAuth().GetStableAccessToken(&req)
		if err != nil {
			return "", err
		}

		if err := rsp.GetResponseError(); err != nil {
			return "", err
		}
		cli.cache.Set(key, rsp.AccessToken, time.Duration(rsp.ExpiresIn)*time.Second)
		return rsp.AccessToken, nil
	}
}

// 拼凑完整的 URI
func (cli *Client) combineURI(url string, req interface{}, withToken bool) (string, error) {
	output := make(map[string]interface{})

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &output,
		TagName:  "query",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return "", err
	}

	err = decoder.Decode(req)
	if err != nil {
		return "", err
	}

	if withToken {
		token, err := cli.AccessToken()
		if err != nil {
			return "", err
		}

		output["access_token"] = token
	}

	return request.EncodeURL(baseURL+url, output)
}

// 用户信息
func (cli *Client) NewAuth() *auth.Auth {
	return auth.NewAuth(cli.request, cli.combineURI)
}