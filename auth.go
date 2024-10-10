package weofficial

import (
	"fmt"

	"github.com/hankin-h-k/weofficial/request"
)

const (
	apiGetAccessToken = "/cgi-bin/token"
	apiGetPaidUnionID = "/wxa/getpaidunionid"
	apiSnsOauth2      = "/sns/oauth2/access_token"
)

// LoginResponse 返回给用户的数据
type LoginResponse struct {
	request.CommonError
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	// 用户在开放平台的唯一标识符
	// 只在满足一定条件的情况下返回
	UnionID string `json:"unionid"`
}

// TokenResponse 获取 access_token 成功返回数据
type TokenResponse struct {
	request.CommonError
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   uint   `json:"expires_in"`   // 凭证有效时间，单位：秒。目前是7200秒之内的值。
}

type SnsOauth2Response struct {
	request.CommonError
	AccessToken    string `json:"access_token"` // 获取到的凭证
	ExpiresIn      uint   `json:"expires_in"`   // 凭证有效时间，单位：秒。目前是7200秒之内的值。
	RefreshToken   string `json:"refresh_token"`
	OpenID         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapshotuser int    `json:"is_snapshotuser"`
	UnionID        string `json:"unionid"`
}

// access_token 缓存 KEY
func (cli *Client) tokenCacheKey() string {
	return fmt.Sprintf("weofficial.%s.access.token", cli.appid)
}

func (cli *Client) GetAccessToken() (*TokenResponse, error) {

	queries := requestQueries{
		"appid":      cli.appid,
		"secret":     cli.secret,
		"grant_type": "client_credential",
	}
	api := baseURL + apiGetAccessToken
	url, err := request.EncodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	res := new(TokenResponse)
	if err := cli.request.Get(url, res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetPaidUnionIDResponse response data
type GetPaidUnionIDResponse struct {
	request.CommonError
	UnionID string `json:"unionid"`
}

// GetPaidUnionID 用户支付完成后，通过微信支付订单号（transaction_id）获取该用户的 UnionId
func (cli *Client) GetPaidUnionID(openID, transactionID string) (*GetPaidUnionIDResponse, error) {
	api := baseURL + apiGetPaidUnionID
	accessToken, err := cli.AccessToken()
	if err != nil {
		return nil, err
	}
	return cli.getPaidUnionID(accessToken, openID, transactionID, api)
}

func (cli *Client) getPaidUnionID(accessToken, openID, transactionID, api string) (*GetPaidUnionIDResponse, error) {
	queries := requestQueries{
		"openid":         openID,
		"access_token":   accessToken,
		"transaction_id": transactionID,
	}

	return cli.getPaidUnionIDRequest(api, queries)
}

// GetPaidUnionIDWithMCH 用户支付完成后，通过微信支付商户订单号和微信支付商户号（out_trade_no 及 mch_id）获取该用户的 UnionId
func (cli *Client) GetPaidUnionIDWithMCH(openID, outTradeNo, mchID string) (*GetPaidUnionIDResponse, error) {
	api := baseURL + apiGetPaidUnionID

	accessToken, err := cli.AccessToken()
	if err != nil {
		return nil, err
	}

	return cli.getPaidUnionIDWithMCH(accessToken, openID, outTradeNo, mchID, api)
}

func (cli *Client) getPaidUnionIDWithMCH(accessToken, openID, outTradeNo, mchID, api string) (*GetPaidUnionIDResponse, error) {
	queries := requestQueries{
		"openid":       openID,
		"mch_id":       mchID,
		"out_trade_no": outTradeNo,
		"access_token": accessToken,
	}

	return cli.getPaidUnionIDRequest(api, queries)
}

func (cli *Client) getPaidUnionIDRequest(api string, queries requestQueries) (*GetPaidUnionIDResponse, error) {
	url, err := request.EncodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	res := new(GetPaidUnionIDResponse)
	if err := cli.request.Get(url, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (cli *Client) Authorize(redirect_uri, scope, state string) (string, error) {
	baseURL := "https://open.weixin.qq.com"
	apiAuthUrl := "/connect/oauth2/authorize"
	api := baseURL + apiAuthUrl
	return cli.authorize(api, redirect_uri, scope, state)
}
func (cli *Client) authorize(api, redirect_uri, scope, state string) (string, error) {
	queries := requestQueries{
		"appid":         cli.appid,
		"redirect_uri":  redirect_uri,
		"response_type": "code",
		"scope":         scope,
		"state":         state,
	}

	url, err := request.EncodeURL(api, queries)
	if err != nil {
		return url, err
	}
	return url + "#wechat_redirect", nil
}

func (cli *Client) Login(code string) (*SnsOauth2Response, error) {

	api := baseURL + apiSnsOauth2
	return cli.login(api, code)
}

func (cli *Client) login(api, code string) (*SnsOauth2Response, error) {
	queries := requestQueries{
		"appid":      cli.appid,
		"secret":     cli.secret,
		"code":       code,
		"grant_type": "authorization_code",
	}

	url, err := request.EncodeURL(api, queries)
	if err != nil {
		return nil, err
	}
	res := new(SnsOauth2Response)
	if err := cli.request.Get(url, res); err != nil {
		return nil, err
	}
	return res, nil
}
