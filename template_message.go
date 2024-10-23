package weofficial

import (
	"fmt"

	"github.com/hankin-h-k/weofficial/request"
)

var (
	apiAddTeamplate string = "/cgi-bin/template/api_add_template"
	apiTemplateList string = "/cgi-bin/template/get_all_private_template"
	apiSendTemplate string = "/cgi-bin/message/template/send"
)

type AddTemplateResponse struct {
	request.CommonError
	TemplateId string `json:"template_id"`
}

type TemplateList struct {
	TemplateId      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

type TemplateListResponse struct {
	request.CommonError
	TemplateList []TemplateList `json:"template_list"`
}

type SendTemplateResponse struct {
	request.CommonError
	Msgid uint64 `json:"msgid"`
}

func (cli *Client) AddTemplate(template_id_short string, keyword_name_list []string) (*AddTemplateResponse, error) {
	api := baseURL + apiAddTeamplate

	token, err := cli.AccessToken()
	if err != nil {
		return nil, err
	}
	return cli.addTemplate(api, token, template_id_short, keyword_name_list)
}

func (cli *Client) addTemplate(api, token, template_id_short string, keyword_name_list []string) (*AddTemplateResponse, error) {
	api, err := tokenAPI(api, token)
	if err != nil {
		return nil, err
	}

	params := requestParams{
		"template_id_short": template_id_short,
		"keyword_name_list": keyword_name_list,
	}

	res := new(AddTemplateResponse)
	err = cli.request.Post(api, params, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (cli *Client) GetTemplateList() (*TemplateListResponse, error) {
	api := baseURL + apiTemplateList

	token, err := cli.AccessToken()
	if err != nil {
		return nil, err
	}
	return cli.templateList(api, token)
}

func (cli *Client) templateList(api, token string) (*TemplateListResponse, error) {
	api, err := tokenAPI(api, token)
	if err != nil {
		return nil, err
	}

	res := new(TemplateListResponse)
	err = cli.request.Get(api, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (cli *Client) SendTemplate(param map[string]interface{}) (*SendTemplateResponse, error) {
	api := baseURL + apiSendTemplate

	token, err := cli.AccessToken()
	if err != nil {
		return nil, err
	}
	return cli.sendTemplate(api, token, param)
}

func (cli *Client) sendTemplate(api, token string, param map[string]interface{}) (*SendTemplateResponse, error) {
	api, err := tokenAPI(api, token)
	if err != nil {
		return nil, err
	}
	fmt.Println(param)
	params := requestParams{
		"touser":      param["touser"],
		"template_id": param["template_id"],
		"url":         param["url"],
		"miniprogram": param["miniprogram"],
		"data":        param["data"],
	}

	res := new(SendTemplateResponse)
	err = cli.request.Post(api, params, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
