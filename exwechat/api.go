package exwechat

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wealthworks/go-tencent-api/client"
)

const (
	urlToken   = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	urlGetUser = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	urlAddUser = "https://qyapi.weixin.qq.com/cgi-bin/user/create"
	urlDelUser = "https://qyapi.weixin.qq.com/cgi-bin/user/delete"
)

var (
	corpId, corpSecret string
	holder             *client.TokenHolder
)

func init() {
	corpId = os.Getenv("EXWECHAT_CORP_ID")
	corpSecret = os.Getenv("EXWECHAT_CORP_SECRET")
	if corpId == "" || corpSecret == "" {
		panic("EXWECHAT_CORP_ID or EXWECHAT_CORP_SECRET are empty or not found")
	}
	holder = client.NewTokenHolder(urlToken)
	holder.SetClient(corpId, corpSecret)
}

type API struct {
	c  *client.Client
	th *client.TokenHolder
}

func NewAPI() *API {
	c := client.NewClient()
	c.SetContentType("application/json")
	return &API{c, holder}
}

func (a *API) GetUser(userId string) (*User, error) {
	token, err := a.th.GetAuthToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s&userid=%s", urlGetUser, token, userId)

	body, err := a.c.Get(uri)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal(body, user)

	return user, err
}

func (a *API) AddUser(user *User) (err error) {
	var token string
	token, err = a.th.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", urlAddUser, token)
	var data []byte
	data, err = json.Marshal(user)
	if err != nil {
		return
	}

	_, err = a.c.Post(uri, data)
	return
}

func (a *API) DeleteUser(userId string) (err error) {
	var token string
	token, err = a.th.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&userid=%s", urlDelUser, token, userId)

	_, err = a.c.Get(uri)
	return
}
