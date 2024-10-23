package weofficial

import (
	"fmt"

	"github.com/hankin-h-k/weofficial/request"
)

var (
	apiSendSubscribe string = "/cgi-bin/message/subscribe/bizsend"
)

type SendSubscribeResponse struct {
	request.CommonError
}

func (cli *Client) SendSubscribe(param map[string]interface{}) (*SendSubscribeResponse, error) {
	api := baseURL + apiSendSubscribe

	token, err := cli.AccessToken()
	if err != nil {
		return nil, err
	}
	return cli.sendSubscribe(api, token, param)
}

func (cli *Client) sendSubscribe(api, token string, param map[string]interface{}) (*SendSubscribeResponse, error) {
	api, err := tokenAPI(api, token)
	if err != nil {
		return nil, err
	}
	fmt.Println(param)
	params := requestParams{
		"touser":      param["touser"],
		"template_id": param["template_id"],
		"page":        param["page"],
		"miniprogram": param["miniprogram"],
		"data":        param["data"],
	}

	res := new(SendSubscribeResponse)
	err = cli.request.Post(api, params, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
